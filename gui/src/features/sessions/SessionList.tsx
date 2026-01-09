import { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { FAB } from '@/components/FAB'
import { useModal } from '@/components/Modal'
import { api } from '@/api'
import { useUIStore } from '@/stores'
import type { Session } from '@/domain'
import { useSessions } from './useSessions'
import { SessionCard } from './SessionCard'
import styles from './SessionList.module.css'

export function SessionList() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const uiStore = useUIStore()
  const { sessions, loading, error, refetch } = useSessions(projectId ?? '')
  const { openCreateSessionModal, openConfirmModal, closeModal } = useModal()
  const [actionLoading, setActionLoading] = useState(false)

  useEffect(() => {
    if (!projectId) {
      navigate('/', { replace: true })
      return
    }

    if (!uiStore.currentProjectName) {
      api.projects
        .get(projectId)
        .then((project) => {
          uiStore.setCurrentProject(project.id, project.name)
        })
        .catch(() => {
          toast.error('Failed to load project')
          navigate('/', { replace: true })
        })
    }
  }, [projectId, navigate, uiStore])

  const handleCreate = useCallback(() => {
    if (!projectId) return

    const createSession = async (data: {
      name: string
      description?: string
      targetMultiplexer: 'tmux' | 'zellij'
    }) => {
      setActionLoading(true)
      try {
        const result = await api.sessions.create({
          projectId,
          session: {
            name: data.name,
            description: data.description,
            targetMultiplexer: data.targetMultiplexer,
            windows: [{ name: 'main' }],
          },
        })
        closeModal()
        toast.success('Session created')
        await refetch()
        return result.session.id
      } catch (err) {
        toast.error(err instanceof Error ? err.message : 'Failed to create session')
        return null
      } finally {
        setActionLoading(false)
      }
    }

    openCreateSessionModal({
      onSubmit: createSession,
      onSubmitAndEdit: async (data) => {
        const sessionId = await createSession(data)
        if (sessionId) {
          navigate(`/sessions/${sessionId}/edit`)
        }
      },
      loading: actionLoading,
    })
  }, [projectId, openCreateSessionModal, closeModal, refetch, actionLoading, navigate])

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.altKey && e.key === 'n') {
        e.preventDefault()
        handleCreate()
      }
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [handleCreate])

  const handleDelete = (session: Session) => {
    openConfirmModal({
      title: 'Delete Session',
      message: (
        <>
          Are you sure you want to delete <strong>{session.name}</strong>?
        </>
      ),
      confirmText: 'Delete',
      variant: 'danger',
      onConfirm: async () => {
        setActionLoading(true)
        try {
          await api.sessions.delete(session.id)
          closeModal()
          toast.success('Session deleted')
          await refetch()
        } catch (err) {
          toast.error(err instanceof Error ? err.message : 'Failed to delete session')
        } finally {
          setActionLoading(false)
        }
      },
      loading: actionLoading,
    })
  }

  if (!projectId) {
    return <div className={styles.loading}>Redirecting...</div>
  }

  if (loading) {
    return <div className={styles.loading}>Loading sessions...</div>
  }

  if (error) {
    return <div className={styles.error}>{error}</div>
  }

  return (
    <div className={styles.container}>
      {sessions.length === 0 ? (
        <div className={styles.empty}>
          <p>No sessions yet</p>
          <p className={styles.hint}>Click the + button to create your first session</p>
        </div>
      ) : (
        <div className={styles.grid}>
          {sessions.map((session) => (
            <SessionCard key={session.id} session={session} onDelete={handleDelete} />
          ))}
        </div>
      )}
      <FAB onClick={handleCreate} tooltip="Create session (Alt+N)" />
    </div>
  )
}
