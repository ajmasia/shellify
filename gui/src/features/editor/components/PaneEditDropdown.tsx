import { useRef, useEffect } from 'react'
import type { EditorPane } from '../types'
import { useSessionEditor } from '../hooks/useSessionEditor'
import styles from './PaneEditDropdown.module.css'

interface PaneEditDropdownProps {
  pane: EditorPane
  onClose: () => void
}

export function PaneEditDropdown({ pane, onClose }: PaneEditDropdownProps) {
  const store = useSessionEditor()
  const dropdownRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        onClose()
      }
    }

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose()
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    document.addEventListener('keydown', handleEscape)

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
      document.removeEventListener('keydown', handleEscape)
    }
  }, [onClose])

  const handleChange = (field: 'name' | 'command' | 'workingDirectory', value: string) => {
    store.updatePane(pane.id, { [field]: value || undefined })
  }

  return (
    <div ref={dropdownRef} className={styles.dropdown} onClick={(e) => e.stopPropagation()}>
      <div className={styles.field}>
        <label className={styles.label} htmlFor={`pane-name-${pane.id}`}>
          Name
        </label>
        <input
          id={`pane-name-${pane.id}`}
          type="text"
          className={styles.input}
          value={pane.name || ''}
          onChange={(e) => handleChange('name', e.target.value)}
          placeholder="pane"
          autoFocus
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor={`pane-command-${pane.id}`}>
          Command
        </label>
        <input
          id={`pane-command-${pane.id}`}
          type="text"
          className={styles.input}
          value={pane.command || ''}
          onChange={(e) => handleChange('command', e.target.value)}
          placeholder="nvim ."
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor={`pane-workdir-${pane.id}`}>
          Working Directory
        </label>
        <input
          id={`pane-workdir-${pane.id}`}
          type="text"
          className={styles.input}
          value={pane.workingDirectory || ''}
          onChange={(e) => handleChange('workingDirectory', e.target.value)}
          placeholder="~/projects"
        />
      </div>
    </div>
  )
}
