import { useState, useRef, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { useSessionEditor } from '../hooks/useSessionEditor'
import { CloseIcon, AddIcon } from '@/components/Icons'
import { Tooltip } from '@/components/Tooltip'
import styles from './WindowTabs.module.css'

export const WindowTabs = observer(function WindowTabs() {
  const store = useSessionEditor()
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editValue, setEditValue] = useState('')
  const [draggedId, setDraggedId] = useState<string | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (editingId && inputRef.current) {
      inputRef.current.focus()
      inputRef.current.select()
    }
  }, [editingId])

  const handleDoubleClick = (windowId: string, currentName: string) => {
    setEditingId(windowId)
    setEditValue(currentName)
  }

  const handleEditSubmit = () => {
    if (editingId && editValue.trim()) {
      store.updateWindow(editingId, { name: editValue.trim() })
    }
    setEditingId(null)
    setEditValue('')
  }

  const handleEditKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleEditSubmit()
    } else if (e.key === 'Escape') {
      setEditingId(null)
      setEditValue('')
    }
  }

  const handleRemoveWindow = (e: React.MouseEvent, windowId: string) => {
    e.stopPropagation()
    store.removeWindow(windowId)
  }

  const handleDragStart = (e: React.DragEvent, windowId: string) => {
    setDraggedId(windowId)
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', windowId)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
  }

  const handleDrop = (e: React.DragEvent, targetId: string) => {
    e.preventDefault()
    if (draggedId && draggedId !== targetId) {
      store.reorderWindows(draggedId, targetId)
    }
    setDraggedId(null)
  }

  const handleDragEnd = () => {
    setDraggedId(null)
  }

  return (
    <div className={styles.tabs}>
      {store.session.windows.map((window) => (
        <div
          key={window.id}
          className={`${styles.tab} ${store.selectedWindowId === window.id ? styles.active : ''} ${draggedId === window.id ? styles.dragging : ''}`}
          onClick={() => store.selectWindow(window.id)}
          draggable={editingId !== window.id}
          onDragStart={(e) => handleDragStart(e, window.id)}
          onDragOver={handleDragOver}
          onDrop={(e) => handleDrop(e, window.id)}
          onDragEnd={handleDragEnd}
        >
          {editingId === window.id ? (
            <input
              ref={inputRef}
              type="text"
              className={styles.tabInput}
              value={editValue}
              onChange={(e) => setEditValue(e.target.value)}
              onBlur={handleEditSubmit}
              onKeyDown={handleEditKeyDown}
              onClick={(e) => e.stopPropagation()}
            />
          ) : (
            <span
              className={styles.tabName}
              onDoubleClick={() => handleDoubleClick(window.id, window.name)}
            >
              {window.name}
            </span>
          )}
          {store.session.windows.length > 1 && (
            <Tooltip content="Close window">
              <button
                className={styles.closeButton}
                onClick={(e) => handleRemoveWindow(e, window.id)}
              >
                <CloseIcon />
              </button>
            </Tooltip>
          )}
        </div>
      ))}
      <Tooltip content="Add window">
        <button className={styles.addButton} onClick={() => store.addWindow()}>
          <AddIcon />
        </button>
      </Tooltip>
    </div>
  )
})
