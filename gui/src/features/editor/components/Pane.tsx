import { useState, useRef, useCallback } from 'react'
import { observer } from 'mobx-react-lite'
import type { EditorPane, Direction } from '../types'
import { useSessionEditor } from '../hooks/useSessionEditor'
import { PaneGrid } from './PaneGrid'
import { PaneEditDropdown } from './PaneEditDropdown'
import { FolderIcon, DeleteIcon, EditIcon } from '@/components/Icons'
import { Tooltip } from '@/components/Tooltip'
import styles from './Pane.module.css'

interface PaneProps {
  pane: EditorPane
  direction: Direction
  isLast: boolean
}

export const Pane = observer(function Pane({ pane, direction, isLast }: PaneProps) {
  const store = useSessionEditor()
  const paneRef = useRef<HTMLDivElement>(null)
  const resizeHandleRef = useRef<HTMLDivElement>(null)
  const isResizingRef = useRef(false)
  const [isEditing, setIsEditing] = useState(false)

  const isSelected = store.selectedPaneId === pane.id
  const hasChildren = pane.children && pane.children.length > 0

  const handleClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    if (!hasChildren) {
      store.selectPane(pane.id)
    }
  }

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation()
    store.removePane(pane.id)
  }

  const handleEdit = (e: React.MouseEvent) => {
    e.stopPropagation()
    store.selectPane(pane.id)
    setIsEditing(true)
  }

  const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      if (isLast) return
      e.preventDefault()
      isResizingRef.current = true
      store.startResize()

      const parentElement = paneRef.current?.parentElement
      if (!parentElement) return

      const startPos = direction === 'horizontal' ? e.clientX : e.clientY
      const startSize = pane.size

      const handleMouseMove = (moveEvent: MouseEvent) => {
        if (!isResizingRef.current || !parentElement) return

        const parentSize =
          direction === 'horizontal' ? parentElement.clientWidth : parentElement.clientHeight

        const currentPos = direction === 'horizontal' ? moveEvent.clientX : moveEvent.clientY
        const delta = currentPos - startPos
        const deltaPercent = (delta / parentSize) * 100

        store.resizePane(pane.id, startSize + deltaPercent)
      }

      const handleMouseUp = () => {
        if (isResizingRef.current) {
          isResizingRef.current = false
          store.commitResize()
        }
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
        document.body.style.cursor = ''
        document.body.style.userSelect = ''
      }

      document.addEventListener('mousemove', handleMouseMove)
      document.addEventListener('mouseup', handleMouseUp)
      document.body.style.cursor = direction === 'horizontal' ? 'col-resize' : 'row-resize'
      document.body.style.userSelect = 'none'
    },
    [direction, isLast, pane.id, pane.size, store]
  )

  const sizeStyle =
    direction === 'horizontal' ? { width: `${pane.size}%` } : { height: `${pane.size}%` }

  const paneClasses = [
    styles.pane,
    !hasChildren && styles.leaf,
    !hasChildren && isSelected && styles.selected,
  ]
    .filter(Boolean)
    .join(' ')

  return (
    <div ref={paneRef} className={paneClasses} style={sizeStyle} onClick={handleClick}>
      {hasChildren ? (
        <PaneGrid panes={pane.children!} direction={pane.direction!} />
      ) : (
        <>
          {/* Terminal simulation */}
          <div className={styles.terminal}>
            {/* Prompt line with command */}
            <div className={styles.promptLine}>
              <span className={styles.prompt}>❯</span>
              {pane.command && <span className={styles.command}>{pane.command}</span>}
            </div>

            {/* Path info */}
            {pane.workingDirectory && (
              <div className={styles.pathInfo}>
                <FolderIcon />
                <span className={styles.pathText}>{pane.workingDirectory}</span>
              </div>
            )}
          </div>

          {/* Pane name badge */}
          <div className={styles.paneBadge}>{pane.name || 'pane'}</div>

          {/* Resize percentage indicator */}
          {store.isResizing && <div className={styles.sizeIndicator}>{Math.round(pane.size)}%</div>}

          {/* Edit dropdown */}
          {isEditing && <PaneEditDropdown pane={pane} onClose={() => setIsEditing(false)} />}

          {/* Actions */}
          <div className={`${styles.paneActions} ${isEditing ? styles.actionsVisible : ''}`}>
            <Tooltip content="Edit pane">
              <button
                className={`${styles.actionButton} ${isEditing ? styles.editActive : ''}`}
                onClick={handleEdit}
              >
                <EditIcon />
              </button>
            </Tooltip>
            <Tooltip content="Delete pane">
              <button
                className={`${styles.actionButton} ${styles.deleteButton}`}
                onClick={handleDelete}
              >
                <DeleteIcon />
              </button>
            </Tooltip>
          </div>
        </>
      )}

      {!isLast && (
        <div
          ref={resizeHandleRef}
          className={`${styles.resizeHandle} ${direction === 'horizontal' ? styles.vertical : styles.horizontal}`}
          onMouseDown={handleMouseDown}
        />
      )}
    </div>
  )
})
