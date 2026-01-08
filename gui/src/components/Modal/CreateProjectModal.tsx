import { useState } from 'react'
import { BaseModal } from './BaseModal'
import { Button } from '@/components/Button'
import type { CreateProjectModalProps } from './types'
import styles from './Modal.module.css'

export function CreateProjectModal({
  isOpen,
  onClose,
  onSubmit,
  loading = false,
}: CreateProjectModalProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) return
    onSubmit({ name: name.trim(), description: description.trim() || undefined })
  }

  const handleClose = () => {
    setName('')
    setDescription('')
    onClose()
  }

  return (
    <BaseModal isOpen={isOpen} onClose={handleClose} title="Create Project">
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="project-name">
            Name
          </label>
          <input
            id="project-name"
            type="text"
            className={styles.input}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="My Project"
            autoFocus
            disabled={loading}
          />
        </div>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="project-description">
            Description (optional)
          </label>
          <textarea
            id="project-description"
            className={`${styles.input} ${styles.textarea}`}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Project description..."
            disabled={loading}
          />
        </div>
        <div className={styles.footer}>
          <Button type="button" variant="ghost" onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <Button type="submit" disabled={!name.trim() || loading}>
            Create
          </Button>
        </div>
      </form>
    </BaseModal>
  )
}
