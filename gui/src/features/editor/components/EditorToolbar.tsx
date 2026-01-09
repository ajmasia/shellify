import { observer } from 'mobx-react-lite'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { useSessionEditor } from '../hooks/useSessionEditor'
import { IconButton } from '@/components/IconButton'
import { Tooltip } from '@/components/Tooltip'
import {
  BackIcon,
  NewWindowIcon,
  SplitVerticalIcon,
  SplitHorizontalIcon,
  DeleteIcon,
  UndoIcon,
  RedoIcon,
  SaveIcon,
} from '@/components/Icons'
import { api } from '@/api'
import styles from './EditorToolbar.module.css'

interface EditorToolbarProps {
  sessionId?: string
  projectId?: string
  isNew?: boolean
}

export const EditorToolbar = observer(function EditorToolbar({
  sessionId,
  projectId,
  isNew,
}: EditorToolbarProps) {
  const store = useSessionEditor()
  const navigate = useNavigate()

  const handleBack = () => {
    if (projectId) {
      navigate(`/projects/${projectId}`)
    } else {
      navigate('/')
    }
  }

  const handleAddWindow = () => {
    store.addWindow()
  }

  const handleSplitVertical = () => {
    if (store.selectedPaneId) {
      store.splitPane(store.selectedPaneId, 'vertical')
    }
  }

  const handleSplitHorizontal = () => {
    if (store.selectedPaneId) {
      store.splitPane(store.selectedPaneId, 'horizontal')
    }
  }

  const handleDelete = () => {
    if (store.selectedPaneId) {
      store.removePane(store.selectedPaneId)
    }
  }

  const handleSave = async () => {
    const session = store.toJSON()

    // Flatten recursive pane tree to flat array for API
    type PaneApiFormat = {
      name?: string
      command?: string
      workingDirectory?: string
      size?: number
    }
    function flattenPanes(panes: (typeof session.windows)[0]['panes']): PaneApiFormat[] {
      const result: PaneApiFormat[] = []

      function traverse(pane: (typeof panes)[0]) {
        if (pane.children && pane.children.length > 0) {
          pane.children.forEach(traverse)
        } else {
          result.push({
            name: pane.name || undefined,
            command: pane.command || undefined,
            workingDirectory: pane.workingDirectory,
            size: pane.size,
          })
        }
      }

      panes.forEach(traverse)
      return result
    }

    // Convert editor windows to API format
    const windowsForApi = session.windows.map((w) => ({
      name: w.name,
      workingDirectory: w.workingDirectory,
      panes: flattenPanes(w.panes),
    }))

    try {
      if (isNew && projectId) {
        const result = await api.sessions.create({
          projectId,
          session: {
            name: session.name,
            description: session.description,
            workingDirectory: session.workingDirectory,
            targetMultiplexer: session.targetMultiplexer,
            environment: session.environment,
            preCommands: session.preCommands,
            postCommands: session.postCommands,
            windows: windowsForApi,
          },
        })
        store.markAsSaved()
        toast.success(`Session "${session.name}" created`)
        navigate(`/sessions/${result.session.id}/edit`, { replace: true })
      } else if (sessionId) {
        await api.sessions.update(sessionId, {
          name: session.name,
          description: session.description,
          workingDirectory: session.workingDirectory,
          targetMultiplexer: session.targetMultiplexer,
          environment: session.environment,
          preCommands: session.preCommands,
          postCommands: session.postCommands,
          windows: windowsForApi,
        })
        store.markAsSaved()
        toast.success(`Session "${session.name}" saved`)
      }
    } catch {
      toast.error('Failed to save session')
    }
  }

  return (
    <div className={styles.toolbar}>
      <div className={styles.group}>
        <Tooltip content="Back">
          <IconButton onClick={handleBack} variant="ghost">
            <BackIcon />
          </IconButton>
        </Tooltip>

        <div className={styles.separator} />

        <Tooltip content="Add window">
          <IconButton onClick={handleAddWindow} variant="ghost">
            <NewWindowIcon />
          </IconButton>
        </Tooltip>

        <Tooltip content="Split vertical">
          <IconButton
            onClick={handleSplitVertical}
            variant="ghost"
            disabled={!store.selectedPaneId}
          >
            <SplitVerticalIcon />
          </IconButton>
        </Tooltip>

        <Tooltip content="Split horizontal">
          <IconButton
            onClick={handleSplitHorizontal}
            variant="ghost"
            disabled={!store.selectedPaneId}
          >
            <SplitHorizontalIcon />
          </IconButton>
        </Tooltip>

        <Tooltip content="Delete pane">
          <IconButton onClick={handleDelete} variant="ghost" disabled={!store.selectedPaneId}>
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </div>

      <div className={styles.group}>
        <Tooltip content="Undo">
          <IconButton onClick={() => store.undo()} variant="ghost" disabled={!store.canUndo}>
            <UndoIcon />
          </IconButton>
        </Tooltip>

        <Tooltip content="Redo">
          <IconButton onClick={() => store.redo()} variant="ghost" disabled={!store.canRedo}>
            <RedoIcon />
          </IconButton>
        </Tooltip>

        <div className={styles.separator} />

        <Tooltip content="Save">
          <IconButton onClick={handleSave} variant="primary">
            <SaveIcon />
          </IconButton>
        </Tooltip>
      </div>
    </div>
  )
})
