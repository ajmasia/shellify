import type { ReactNode } from 'react'

export interface BaseModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  children: ReactNode
}

export interface ConfirmModalProps {
  isOpen: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  message: ReactNode
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'default'
  loading?: boolean
}

export interface CreateProjectModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: { name: string; description?: string }) => void
  loading?: boolean
}

export interface CreateSessionModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: {
    name: string
    description?: string
    targetMultiplexer: 'tmux' | 'zellij'
  }) => void
  loading?: boolean
}

export type ModalType =
  | { type: 'confirm'; props: Omit<ConfirmModalProps, 'isOpen' | 'onClose'> }
  | { type: 'createProject'; props: Omit<CreateProjectModalProps, 'isOpen' | 'onClose'> }
  | { type: 'createSession'; props: Omit<CreateSessionModalProps, 'isOpen' | 'onClose'> }
  | null
