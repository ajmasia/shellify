import { makeAutoObservable, toJS } from 'mobx'
import type { MultiplexerType } from '@/domain'
import type { EditorSession, EditorWindow, EditorPane, Direction } from '../types'
import { createDefaultPane, createDefaultWindow, createDefaultSession, generateId } from '../types'

const MAX_WINDOWS = 10

function normalizePane(pane: EditorPane): EditorPane {
  if (!pane.children || pane.children.length === 0) {
    return pane
  }

  pane.children = pane.children.map(normalizePane)

  // Merge children that share the same direction as this container into this level
  const merged: EditorPane[] = []
  for (const child of pane.children) {
    if (child.children && child.children.length > 0 && child.direction === pane.direction) {
      for (const grandchild of child.children) {
        merged.push({ ...grandchild, size: (child.size * grandchild.size) / 100 })
      }
    } else {
      merged.push(child)
    }
  }
  pane.children = merged

  // Collapse single-child containers by hoisting the child's content into this pane
  if (pane.children.length === 1) {
    const only = pane.children[0]
    pane.name = only.name
    pane.command = only.command
    pane.workingDirectory = only.workingDirectory
    pane.environment = only.environment
    if (only.children && only.children.length > 0) {
      pane.children = only.children
      pane.direction = only.direction
    } else {
      delete pane.children
      delete pane.direction
    }
    return normalizePane(pane)
  }

  return pane
}

export class SessionEditorStore {
  session: EditorSession = createDefaultSession()
  selectedPaneId: string | null = null
  selectedWindowId: string | null = null
  isDirty = false
  isResizing = false

  private history: string[] = []
  private historyIndex = -1
  private readonly maxHistorySize = 50
  private isUndoRedo = false
  private lastSavedSnapshot = ''

  constructor() {
    makeAutoObservable(this)
    this.saveSnapshot()
    this.markAsSaved()
  }

  // Undo/Redo
  get canUndo(): boolean {
    return this.historyIndex > 0
  }

  get canRedo(): boolean {
    return this.historyIndex < this.history.length - 1
  }

  undo(): void {
    if (!this.canUndo) return
    this.isUndoRedo = true
    this.historyIndex--
    this.session = JSON.parse(this.history[this.historyIndex])
    this.isUndoRedo = false
    this.updateDirtyState()
  }

  redo(): void {
    if (!this.canRedo) return
    this.isUndoRedo = true
    this.historyIndex++
    this.session = JSON.parse(this.history[this.historyIndex])
    this.isUndoRedo = false
    this.updateDirtyState()
  }

  private saveSnapshot(): void {
    if (this.isUndoRedo) return

    const snapshot = JSON.stringify(toJS(this.session))

    if (this.historyIndex < this.history.length - 1) {
      this.history = this.history.slice(0, this.historyIndex + 1)
    }

    this.history.push(snapshot)

    if (this.history.length > this.maxHistorySize) {
      this.history.shift()
    } else {
      this.historyIndex++
    }

    this.updateDirtyState()
  }

  private updateDirtyState(): void {
    const currentSnapshot = JSON.stringify(toJS(this.session))
    this.isDirty = currentSnapshot !== this.lastSavedSnapshot
  }

  markAsSaved(): void {
    this.lastSavedSnapshot = JSON.stringify(toJS(this.session))
    this.isDirty = false
  }

  // Session management
  loadSession(session: EditorSession): void {
    this.session = session
    this.normalizeAllWindows()
    this.history = []
    this.historyIndex = -1
    this.saveSnapshot()
    this.markAsSaved()

    // Select first pane of first window
    if (session.windows.length > 0) {
      const firstWindow = session.windows[0]
      this.selectedWindowId = firstWindow.id
      const firstPane = this.findFirstLeafPane(firstWindow.panes[0])
      this.selectedPaneId = firstPane
    }
  }

  resetSession(multiplexer: MultiplexerType = 'tmux'): void {
    this.session = createDefaultSession(multiplexer)
    this.history = []
    this.historyIndex = -1
    this.saveSnapshot()
    this.markAsSaved()

    // Select first pane
    if (this.session.windows.length > 0) {
      const firstWindow = this.session.windows[0]
      this.selectedWindowId = firstWindow.id
      this.selectedPaneId = firstWindow.panes[0]?.id ?? null
    }
  }

  updateSession(
    updates: Partial<
      Pick<
        EditorSession,
        | 'name'
        | 'sessionName'
        | 'description'
        | 'workingDirectory'
        | 'targetMultiplexer'
        | 'environment'
        | 'preCommands'
        | 'postCommands'
        | 'defaultWindowId'
      >
    >
  ): void {
    Object.assign(this.session, updates)
    this.saveSnapshot()
  }

  // Window management
  get currentWindow(): EditorWindow | null {
    return this.session.windows.find((w) => w.id === this.selectedWindowId) ?? null
  }

  get canAddWindow(): boolean {
    return this.session.windows.length < MAX_WINDOWS
  }

  addWindow(): string | null {
    if (!this.canAddWindow) return null

    const windowNumber = this.session.windows.length + 1
    const newWindow = createDefaultWindow(`window-${windowNumber}`)
    this.session.windows.push(newWindow)
    this.selectedWindowId = newWindow.id
    this.selectedPaneId = newWindow.panes[0]?.id ?? null
    this.saveSnapshot()
    return newWindow.id
  }

  removeWindow(windowId: string): void {
    if (this.session.windows.length <= 1) return

    const index = this.session.windows.findIndex((w) => w.id === windowId)
    if (index !== -1) {
      this.session.windows.splice(index, 1)

      // Select another window
      const nextIndex = Math.min(index, this.session.windows.length - 1)
      const nextWindow = this.session.windows[nextIndex]
      this.selectedWindowId = nextWindow.id
      this.selectedPaneId = this.findFirstLeafPane(nextWindow.panes[0])

      this.saveSnapshot()
    }
  }

  updateWindow(
    windowId: string,
    updates: Partial<Pick<EditorWindow, 'name' | 'workingDirectory' | 'environment'>>
  ): void {
    const window = this.session.windows.find((w) => w.id === windowId)
    if (window) {
      Object.assign(window, updates)
      this.saveSnapshot()
    }
  }

  reorderWindows(draggedId: string, targetId: string): void {
    const windows = this.session.windows
    const draggedIndex = windows.findIndex((w) => w.id === draggedId)
    const targetIndex = windows.findIndex((w) => w.id === targetId)

    if (draggedIndex !== -1 && targetIndex !== -1) {
      const [dragged] = windows.splice(draggedIndex, 1)
      windows.splice(targetIndex, 0, dragged)
      this.saveSnapshot()
    }
  }

  selectWindow(windowId: string): void {
    const window = this.session.windows.find((w) => w.id === windowId)
    if (window) {
      this.selectedWindowId = windowId
      this.selectedPaneId = this.findFirstLeafPane(window.panes[0])
    }
  }

  // Pane management
  get currentPane(): EditorPane | null {
    if (!this.selectedPaneId) return null
    return this.findPane(this.selectedPaneId)
  }

  selectPane(paneId: string): void {
    const pane = this.findPane(paneId)
    if (pane) {
      this.selectedPaneId = paneId
      // Also update selected window
      const window = this.findWindowByPaneId(paneId)
      if (window) {
        this.selectedWindowId = window.id
      }
    }
  }

  clearSelection(): void {
    this.selectedPaneId = null
  }

  findPane(paneId: string, panes?: EditorPane[]): EditorPane | null {
    const searchPanes = panes ?? this.session.windows.flatMap((w) => w.panes)

    for (const pane of searchPanes) {
      if (pane.id === paneId) return pane
      if (pane.children) {
        const found = this.findPane(paneId, pane.children)
        if (found) return found
      }
    }
    return null
  }

  findWindowByPaneId(paneId: string): EditorWindow | null {
    for (const window of this.session.windows) {
      if (this.findPane(paneId, window.panes)) {
        return window
      }
    }
    return null
  }

  private findFirstLeafPane(pane: EditorPane): string {
    if (!pane.children || pane.children.length === 0) {
      return pane.id
    }
    return this.findFirstLeafPane(pane.children[0])
  }

  splitPane(paneId: string, direction: Direction): string | null {
    const pane = this.findPane(paneId)
    if (!pane) return null

    let newPaneId: string

    if (pane.children && pane.children.length > 0) {
      if (pane.direction === direction) {
        // Same direction: add a new sibling inside this container
        const newPane = createDefaultPane()
        newPane.size = 100 / (pane.children.length + 1)
        pane.children.forEach((child) => {
          child.size = 100 / (pane.children!.length + 1)
        })
        pane.children.push(newPane)
        newPaneId = newPane.id
      } else {
        // Different direction: wrap existing children and new pane in a new container
        const newPane = createDefaultPane()
        newPane.size = 50

        const existingWrapper: EditorPane = {
          id: generateId(),
          name: '',
          command: '',
          size: 50,
          children: pane.children.map((c) => ({ ...c, size: c.size })),
          direction: pane.direction,
        }

        pane.children = [existingWrapper, newPane]
        pane.direction = direction
        newPaneId = newPane.id
      }
    } else {
      // Leaf pane: convert to container with original content + new pane
      const originalContent: EditorPane = {
        id: generateId(),
        name: pane.name,
        command: pane.command,
        workingDirectory: pane.workingDirectory,
        environment: pane.environment,
        size: 50,
      }

      const newPane = createDefaultPane()
      newPane.size = 50

      pane.name = ''
      pane.command = ''
      pane.workingDirectory = undefined
      pane.environment = undefined
      pane.children = [originalContent, newPane]
      pane.direction = direction

      newPaneId = newPane.id
    }

    this.normalizeAllWindows()
    this.selectedPaneId = newPaneId
    this.saveSnapshot()
    return newPaneId
  }

  removePane(paneId: string): string | null {
    for (let i = 0; i < this.session.windows.length; i++) {
      const window = this.session.windows[i]

      // Check if this is the last pane in the window
      if (window.panes.length === 1 && window.panes[0].id === paneId && !window.panes[0].children) {
        // If there are other windows, delete this window
        if (this.session.windows.length > 1) {
          this.session.windows.splice(i, 1)
          const nextWindow = this.session.windows[0]
          this.selectedWindowId = nextWindow.id
          this.normalizeAllWindows()
          const nextPaneId = this.findFirstLeafPane(nextWindow.panes[0])
          this.selectedPaneId = nextPaneId
          this.saveSnapshot()
          return nextPaneId
        }
        // Can't delete the last pane of the last window
        return null
      }

      const result = this.removePaneFromArray(paneId, window.panes, window)
      if (result.removed) {
        this.normalizeAllWindows()
        this.saveSnapshot()
        // Always select a leaf pane to avoid subsequent operations targeting a container
        const nextPaneId = this.findFirstLeafPane(window.panes[0])
        this.selectedPaneId = nextPaneId
        return nextPaneId
      }
    }
    return null
  }

  private normalizeAllWindows(): void {
    for (const window of this.session.windows) {
      window.panes = window.panes.map(normalizePane)
    }
  }

  private removePaneFromArray(
    paneId: string,
    panes: EditorPane[],
    parent: EditorWindow | EditorPane
  ): { removed: boolean; nextPaneId: string | null } {
    const index = panes.findIndex((p) => p.id === paneId)

    if (index !== -1) {
      if (panes.length === 1) {
        return { removed: false, nextPaneId: null }
      }

      const removedSize = panes[index].size
      panes.splice(index, 1)

      const sizePerPane = removedSize / panes.length
      panes.forEach((p) => (p.size += sizePerPane))

      // If only one child left in a container, flatten
      if (panes.length === 1 && 'children' in parent && parent.children) {
        const remainingPane = panes[0]
        parent.name = remainingPane.name
        parent.command = remainingPane.command
        parent.workingDirectory = remainingPane.workingDirectory
        parent.environment = remainingPane.environment
        parent.children = remainingPane.children
        parent.direction = remainingPane.direction
        return { removed: true, nextPaneId: parent.id }
      }

      const nextIndex = Math.min(index, panes.length - 1)
      return { removed: true, nextPaneId: panes[nextIndex].id }
    }

    for (const pane of panes) {
      if (pane.children) {
        const result = this.removePaneFromArray(paneId, pane.children, pane)
        if (result.removed) {
          return result
        }
      }
    }

    return { removed: false, nextPaneId: null }
  }

  updatePane(
    paneId: string,
    updates: Partial<Pick<EditorPane, 'name' | 'command' | 'workingDirectory' | 'environment'>>
  ): void {
    const pane = this.findPane(paneId)
    if (pane) {
      Object.assign(pane, updates)
      this.saveSnapshot()
    }
  }

  resizePane(paneId: string, newSize: number): void {
    const pane = this.findPane(paneId)
    if (!pane) return

    // Find siblings
    let siblings: EditorPane[] | undefined

    for (const window of this.session.windows) {
      if (window.panes.some((p) => p.id === paneId)) {
        siblings = window.panes
        break
      }
    }

    if (!siblings) {
      // Look in nested panes
      for (const window of this.session.windows) {
        siblings = this.findSiblings(paneId, window.panes)
        if (siblings) break
      }
    }

    if (!siblings || siblings.length < 2) return

    const paneIndex = siblings.findIndex((p) => p.id === paneId)
    if (paneIndex === -1 || paneIndex === siblings.length - 1) return

    const nextPane = siblings[paneIndex + 1]
    const totalSize = pane.size + nextPane.size

    const clampedSize = Math.round(Math.max(10, Math.min(newSize, totalSize - 10)))
    const delta = clampedSize - pane.size

    pane.size = clampedSize
    nextPane.size = Math.round(nextPane.size - delta)
  }

  private findSiblings(paneId: string, panes: EditorPane[]): EditorPane[] | undefined {
    for (const pane of panes) {
      if (pane.children) {
        if (pane.children.some((p) => p.id === paneId)) {
          return pane.children
        }
        const found = this.findSiblings(paneId, pane.children)
        if (found) return found
      }
    }
    return undefined
  }

  startResize(): void {
    this.isResizing = true
  }

  commitResize(): void {
    this.isResizing = false
    this.saveSnapshot()
  }

  // Export
  toJSON(): EditorSession {
    return toJS(this.session)
  }
}
