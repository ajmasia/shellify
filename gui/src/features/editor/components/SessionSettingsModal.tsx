import { useState, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { useSessionEditor } from '../hooks/useSessionEditor'
import { BaseModal } from '@/components/Modal/BaseModal'
import { Button } from '@/components/Button'
import { Chip } from '@/components/Chip'
import type { MultiplexerType } from '@/domain'
import styles from './SessionSettingsModal.module.css'

interface EnvVar {
  key: string
  value: string
}

interface SessionSettingsModalProps {
  isOpen: boolean
  onClose: () => void
}

export const SessionSettingsModal = observer(function SessionSettingsModal({
  isOpen,
  onClose,
}: SessionSettingsModalProps) {
  const store = useSessionEditor()

  const [sessionName, setSessionName] = useState('')
  const [description, setDescription] = useState('')
  const [workingDirectory, setWorkingDirectory] = useState('')
  const [targetMultiplexer, setTargetMultiplexer] = useState<MultiplexerType>('tmux')
  const [defaultWindowId, setDefaultWindowId] = useState('')
  const [envVars, setEnvVars] = useState<EnvVar[]>([])
  const [preCommands, setPreCommands] = useState('')
  const [postCommands, setPostCommands] = useState('')

  useEffect(() => {
    if (isOpen) {
      setSessionName(store.session.sessionName || '')
      setDescription(store.session.description || '')
      setWorkingDirectory(store.session.workingDirectory || '')
      setTargetMultiplexer(store.session.targetMultiplexer || 'tmux')
      setDefaultWindowId(store.session.defaultWindowId || store.session.windows[0]?.id || '')
      setEnvVars(
        store.session.environment
          ? Object.entries(store.session.environment).map(([key, value]) => ({ key, value }))
          : []
      )
      setPreCommands(store.session.preCommands?.join('\n') || '')
      setPostCommands(store.session.postCommands?.join('\n') || '')
    }
  }, [isOpen, store.session])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    const environment: Record<string, string> = {}
    envVars.forEach((env) => {
      if (env.key.trim()) {
        environment[env.key.trim()] = env.value
      }
    })

    store.updateSession({
      sessionName: sessionName.trim() || undefined,
      description: description.trim() || undefined,
      workingDirectory: workingDirectory.trim() || undefined,
      targetMultiplexer,
      defaultWindowId: defaultWindowId || undefined,
      environment: Object.keys(environment).length > 0 ? environment : undefined,
      preCommands: preCommands.trim()
        ? preCommands
            .split('\n')
            .map((c) => c.trim())
            .filter(Boolean)
        : undefined,
      postCommands: postCommands.trim()
        ? postCommands
            .split('\n')
            .map((c) => c.trim())
            .filter(Boolean)
        : undefined,
    })
    onClose()
  }

  const handleAddEnvVar = () => {
    setEnvVars([...envVars, { key: '', value: '' }])
  }

  const handleRemoveEnvVar = (index: number) => {
    setEnvVars(envVars.filter((_, i) => i !== index))
  }

  const handleEnvVarChange = (index: number, field: 'key' | 'value', value: string) => {
    const newEnvVars = [...envVars]
    newEnvVars[index][field] = value
    setEnvVars(newEnvVars)
  }

  const handleInputKeyDown = (e: React.KeyboardEvent) => {
    e.stopPropagation()
  }

  return (
    <BaseModal
      isOpen={isOpen}
      onClose={onClose}
      title="Session Settings"
      className={styles.wideModal}
    >
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="settings-session-name">
            Session Name
          </label>
          <input
            id="settings-session-name"
            type="text"
            className={styles.input}
            value={sessionName}
            onChange={(e) => setSessionName(e.target.value)}
            onKeyDown={handleInputKeyDown}
            placeholder="my-session"
            autoFocus
          />
          <span className={styles.hint}>Name used in the multiplexer (no spaces)</span>
        </div>

        <div className={styles.field}>
          <label className={styles.label} htmlFor="settings-description">
            Description
          </label>
          <textarea
            id="settings-description"
            className={`${styles.input} ${styles.textarea}`}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            onKeyDown={handleInputKeyDown}
            placeholder="Brief description of this session"
            rows={2}
          />
        </div>

        <div className={styles.field}>
          <label className={styles.label} htmlFor="settings-workdir">
            Working Directory
          </label>
          <input
            id="settings-workdir"
            type="text"
            className={styles.input}
            value={workingDirectory}
            onChange={(e) => setWorkingDirectory(e.target.value)}
            onKeyDown={handleInputKeyDown}
            placeholder="e.g., ~/projects/my-app"
          />
        </div>

        <div className={styles.field}>
          <label className={styles.label}>Target Multiplexer</label>
          <div className={styles.chipsContainer}>
            <Chip
              selected={targetMultiplexer === 'tmux'}
              onClick={() => setTargetMultiplexer('tmux')}
            >
              tmux
            </Chip>
            <Chip
              selected={targetMultiplexer === 'zellij'}
              onClick={() => setTargetMultiplexer('zellij')}
            >
              zellij
            </Chip>
          </div>
        </div>

        <div className={styles.field}>
          <label className={styles.label}>Default Window</label>
          <div className={styles.chipsContainer}>
            {store.session.windows.map((window) => (
              <Chip
                key={window.id}
                selected={defaultWindowId === window.id}
                onClick={() => setDefaultWindowId(window.id)}
              >
                {window.name}
              </Chip>
            ))}
          </div>
        </div>

        <div className={styles.field}>
          <div className={styles.labelRow}>
            <label className={styles.label}>Environment Variables</label>
            <button type="button" className={styles.addButton} onClick={handleAddEnvVar}>
              +
            </button>
          </div>
          <div className={styles.envVarList}>
            {envVars.map((env, index) => (
              <div key={index} className={styles.envVarRow}>
                <input
                  type="text"
                  value={env.key}
                  onChange={(e) => handleEnvVarChange(index, 'key', e.target.value)}
                  onKeyDown={handleInputKeyDown}
                  placeholder="KEY"
                  className={`${styles.input} ${styles.envKey}`}
                />
                <span className={styles.envEquals}>=</span>
                <input
                  type="text"
                  value={env.value}
                  onChange={(e) => handleEnvVarChange(index, 'value', e.target.value)}
                  onKeyDown={handleInputKeyDown}
                  placeholder="value"
                  className={`${styles.input} ${styles.envValue}`}
                />
                <button
                  type="button"
                  className={styles.removeButton}
                  onClick={() => handleRemoveEnvVar(index)}
                >
                  ×
                </button>
              </div>
            ))}
            {envVars.length === 0 && (
              <div className={styles.emptyEnv}>No environment variables defined</div>
            )}
          </div>
        </div>

        <div className={styles.field}>
          <label className={styles.label} htmlFor="settings-pre-commands">
            Pre-session Commands
          </label>
          <textarea
            id="settings-pre-commands"
            className={`${styles.input} ${styles.textarea} ${styles.monospace}`}
            value={preCommands}
            onChange={(e) => setPreCommands(e.target.value)}
            onKeyDown={handleInputKeyDown}
            placeholder="Commands to run before creating windows (one per line)"
            rows={3}
          />
        </div>

        <div className={styles.field}>
          <label className={styles.label} htmlFor="settings-post-commands">
            Post-session Commands
          </label>
          <textarea
            id="settings-post-commands"
            className={`${styles.input} ${styles.textarea} ${styles.monospace}`}
            value={postCommands}
            onChange={(e) => setPostCommands(e.target.value)}
            onKeyDown={handleInputKeyDown}
            placeholder="Commands to run when closing session (one per line)"
            rows={3}
          />
        </div>

        <div className={styles.footer}>
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit">Save</Button>
        </div>
      </form>
    </BaseModal>
  )
})
