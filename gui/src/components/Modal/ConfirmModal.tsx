import { BaseModal } from './BaseModal'
import { Button } from '@/components/Button'
import type { ConfirmModalProps } from './types'
import styles from './Modal.module.css'

export function ConfirmModal({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  variant = 'default',
  loading = false,
}: ConfirmModalProps) {
  const handleConfirm = () => {
    onConfirm()
  }

  return (
    <BaseModal isOpen={isOpen} onClose={onClose} title={title}>
      <p className={styles.message}>{message}</p>
      <div className={styles.footer}>
        <Button variant="ghost" onClick={onClose} disabled={loading}>
          {cancelText}
        </Button>
        <Button
          variant={variant === 'danger' ? 'danger' : 'primary'}
          onClick={handleConfirm}
          disabled={loading}
        >
          {confirmText}
        </Button>
      </div>
    </BaseModal>
  )
}
