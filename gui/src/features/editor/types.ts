import type { Session, Window, Pane, MultiplexerType } from '@/domain'

export type Direction = 'horizontal' | 'vertical'

export interface EditorPane extends Omit<Pane, 'children'> {
  children?: EditorPane[]
  direction?: Direction
}

export interface EditorWindow extends Omit<Window, 'panes'> {
  panes: EditorPane[]
  rootDirection: Direction
}

export interface EditorSession extends Omit<Session, 'windows'> {
  windows: EditorWindow[]
}

export interface EditorState {
  session: EditorSession
  selectedPaneId: string | null
  selectedWindowId: string | null
  isDirty: boolean
}

export function createDefaultPane(name = 'pane'): EditorPane {
  return {
    id: generateId(),
    name,
    command: '',
    size: 100,
  }
}

export function createDefaultWindow(name = 'window'): EditorWindow {
  return {
    id: generateId(),
    name,
    panes: [createDefaultPane()],
    rootDirection: 'horizontal',
  }
}

export function createDefaultSession(multiplexer: MultiplexerType = 'tmux'): EditorSession {
  return {
    id: generateId(),
    name: 'new-session',
    sessionName: 'new-session',
    targetMultiplexer: multiplexer,
    windows: [createDefaultWindow('main')],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }
}

export function generateId(): string {
  return Math.random().toString(36).substring(2, 9)
}
