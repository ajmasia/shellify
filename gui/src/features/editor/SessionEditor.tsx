import { useEffect, useState } from 'react'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { api } from '@/api'
import { SessionEditorProvider } from './stores/SessionEditorProvider'
import { useSessionEditor } from './hooks/useSessionEditor'
import { EditorToolbar } from './components/EditorToolbar'
import { WindowTabs } from './components/WindowTabs'
import { PaneGrid } from './components/PaneGrid'
import { PaneConfigPanel } from './components/PaneConfigPanel'
import { TerminalMonitor } from './components/TerminalMonitor'
import type { EditorSession } from './types'
import type { MultiplexerType } from '@/domain'
import styles from './SessionEditor.module.css'

function SessionEditorContent() {
  const { sessionId } = useParams<{ sessionId: string }>()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const store = useSessionEditor()
  const [isLoading, setIsLoading] = useState(true)
  const [loadedProjectId, setLoadedProjectId] = useState<string | undefined>()

  const isNew = sessionId === 'new'
  const projectIdFromParams = searchParams.get('projectId') ?? undefined
  const projectId = projectIdFromParams ?? loadedProjectId
  const multiplexer = (searchParams.get('multiplexer') as MultiplexerType) ?? 'tmux'

  useEffect(() => {
    const loadSession = async () => {
      setIsLoading(true)
      try {
        if (isNew) {
          store.resetSession(multiplexer)
        } else if (sessionId) {
          const result = await api.sessions.get(sessionId)
          setLoadedProjectId(result.projectId)
          store.loadSession(result.session as unknown as EditorSession)
        }
      } catch {
        toast.error('Failed to load session')
        navigate('/')
      } finally {
        setIsLoading(false)
      }
    }

    loadSession()
  }, [sessionId, isNew, multiplexer, store, navigate])

  // Keyboard shortcuts for undo/redo
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey || e.metaKey) {
        if (e.key === 'z' && !e.shiftKey) {
          e.preventDefault()
          store.undo()
        } else if (e.key === 'y' || (e.key === 'z' && e.shiftKey)) {
          e.preventDefault()
          store.redo()
        }
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [store])

  if (isLoading) {
    return (
      <div className={styles.editor}>
        <div className={styles.loading}>Loading session...</div>
      </div>
    )
  }

  const currentWindow = store.currentWindow

  return (
    <div className={styles.editor}>
      <EditorToolbar
        sessionId={isNew ? undefined : sessionId}
        projectId={projectId}
        isNew={isNew}
      />
      <WindowTabs />
      <div className={styles.content}>
        <div className={styles.preview}>
          <TerminalMonitor>
            {currentWindow && (
              <PaneGrid panes={currentWindow.panes} direction={currentWindow.rootDirection} />
            )}
          </TerminalMonitor>
        </div>
        <PaneConfigPanel />
      </div>
    </div>
  )
}

const SessionEditorContentObserver = observer(SessionEditorContent)

export function SessionEditor() {
  return (
    <SessionEditorProvider>
      <SessionEditorContentObserver />
    </SessionEditorProvider>
  )
}
