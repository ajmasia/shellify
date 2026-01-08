import { useState, useCallback, type ReactNode } from 'react'
import { ModalContext } from './context'
import { ConfirmModal } from './ConfirmModal'
import { CreateProjectModal } from './CreateProjectModal'
import { CreateSessionModal } from './CreateSessionModal'
import type {
  ModalType,
  ConfirmModalProps,
  CreateProjectModalProps,
  CreateSessionModalProps,
} from './types'

interface ModalProviderProps {
  children: ReactNode
}

export function ModalProvider({ children }: ModalProviderProps) {
  const [modal, setModal] = useState<ModalType>(null)

  const closeModal = useCallback(() => {
    setModal(null)
  }, [])

  const openConfirmModal = useCallback((props: Omit<ConfirmModalProps, 'isOpen' | 'onClose'>) => {
    setModal({ type: 'confirm', props })
  }, [])

  const openCreateProjectModal = useCallback(
    (props: Omit<CreateProjectModalProps, 'isOpen' | 'onClose'>) => {
      setModal({ type: 'createProject', props })
    },
    []
  )

  const openCreateSessionModal = useCallback(
    (props: Omit<CreateSessionModalProps, 'isOpen' | 'onClose'>) => {
      setModal({ type: 'createSession', props })
    },
    []
  )

  return (
    <ModalContext.Provider
      value={{ openConfirmModal, openCreateProjectModal, openCreateSessionModal, closeModal }}
    >
      {children}
      {modal?.type === 'confirm' && (
        <ConfirmModal isOpen={true} onClose={closeModal} {...modal.props} />
      )}
      {modal?.type === 'createProject' && (
        <CreateProjectModal isOpen={true} onClose={closeModal} {...modal.props} />
      )}
      {modal?.type === 'createSession' && (
        <CreateSessionModal isOpen={true} onClose={closeModal} {...modal.props} />
      )}
    </ModalContext.Provider>
  )
}
