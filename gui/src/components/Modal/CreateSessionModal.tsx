import { useState } from 'react'
import { BaseModal } from './BaseModal'
import { Button } from '@/components/Button'
import type { CreateSessionModalProps } from './types'
import styles from './Modal.module.css'

export function CreateSessionModal({
  isOpen,
  onClose,
  onSubmit,
  onSubmitAndEdit,
  loading = false,
}: CreateSessionModalProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [targetMultiplexer, setTargetMultiplexer] = useState<'tmux' | 'zellij'>('tmux')

  const getFormData = () => ({
    name: name.trim(),
    description: description.trim() || undefined,
    targetMultiplexer,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) return
    onSubmit(getFormData())
  }

  const handleSubmitAndEdit = () => {
    if (!name.trim() || !onSubmitAndEdit) return
    onSubmitAndEdit(getFormData())
  }

  const handleClose = () => {
    setName('')
    setDescription('')
    setTargetMultiplexer('tmux')
    onClose()
  }

  return (
    <BaseModal isOpen={isOpen} onClose={handleClose} title="Create Session">
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="session-name">
            Name
          </label>
          <input
            id="session-name"
            type="text"
            className={styles.input}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="My Session"
            autoFocus
            disabled={loading}
          />
        </div>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="session-description">
            Description (optional)
          </label>
          <textarea
            id="session-description"
            className={`${styles.input} ${styles.textarea}`}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Session description..."
            disabled={loading}
          />
        </div>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="session-multiplexer">
            Multiplexer
          </label>
          <select
            id="session-multiplexer"
            className={styles.select}
            value={targetMultiplexer}
            onChange={(e) => setTargetMultiplexer(e.target.value as 'tmux' | 'zellij')}
            disabled={loading}
          >
            <option value="tmux">tmux</option>
            <option value="zellij">zellij</option>
          </select>
        </div>
        <div className={styles.footer}>
          <Button type="button" variant="ghost" onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <div className={styles.footerActions}>
            {onSubmitAndEdit && (
              <Button
                type="button"
                variant="secondary"
                onClick={handleSubmitAndEdit}
                disabled={!name.trim() || loading}
              >
                Create + Edit
              </Button>
            )}
            <Button type="submit" disabled={!name.trim() || loading}>
              Create
            </Button>
          </div>
        </div>
      </form>
    </BaseModal>
  )
}
