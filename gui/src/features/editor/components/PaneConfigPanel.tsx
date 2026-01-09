import { observer } from 'mobx-react-lite'
import { useSessionEditor } from '../hooks/useSessionEditor'
import styles from './PaneConfigPanel.module.css'

export const PaneConfigPanel = observer(function PaneConfigPanel() {
  const store = useSessionEditor()
  const pane = store.currentPane

  if (!pane) {
    return (
      <div className={styles.panel}>
        <div className={styles.empty}>Select a pane to edit</div>
      </div>
    )
  }

  const handleChange = (field: 'name' | 'command' | 'workingDirectory', value: string) => {
    store.updatePane(pane.id, { [field]: value || undefined })
  }

  return (
    <div className={styles.panel}>
      <h3 className={styles.title}>Pane Settings</h3>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="pane-name">
          Name
        </label>
        <input
          id="pane-name"
          type="text"
          className={styles.input}
          value={pane.name || ''}
          onChange={(e) => handleChange('name', e.target.value)}
          placeholder="pane"
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="pane-command">
          Command
        </label>
        <input
          id="pane-command"
          type="text"
          className={styles.input}
          value={pane.command || ''}
          onChange={(e) => handleChange('command', e.target.value)}
          placeholder="nvim ."
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="pane-workdir">
          Working Directory
        </label>
        <input
          id="pane-workdir"
          type="text"
          className={styles.input}
          value={pane.workingDirectory || ''}
          onChange={(e) => handleChange('workingDirectory', e.target.value)}
          placeholder="~/projects"
        />
      </div>
    </div>
  )
})
