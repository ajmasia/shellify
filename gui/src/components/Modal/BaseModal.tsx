import { useEffect } from 'react'
import { ModalPortal } from './ModalPortal'
import { IconButton } from '@/components/IconButton'
import { CloseIcon } from '@/components/Icons'
import type { BaseModalProps } from './types'
import styles from './Modal.module.css'

export function BaseModal({ isOpen, onClose, title, children }: BaseModalProps) {
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose()
      }
    }

    if (isOpen) {
      document.addEventListener('keydown', handleEscape)
      document.body.style.overflow = 'hidden'
    }

    return () => {
      document.removeEventListener('keydown', handleEscape)
      document.body.style.overflow = ''
    }
  }, [isOpen, onClose])

  if (!isOpen) return null

  return (
    <ModalPortal>
      <div className={styles.backdrop} onClick={onClose}>
        <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
          <div className={styles.header}>
            <h2 className={styles.title}>{title}</h2>
            <IconButton aria-label="Close modal" onClick={onClose} className={styles.closeButton}>
              <CloseIcon size={20} />
            </IconButton>
          </div>
          <div className={styles.content}>{children}</div>
        </div>
      </div>
    </ModalPortal>
  )
}
