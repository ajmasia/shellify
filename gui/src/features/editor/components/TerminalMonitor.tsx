import type { ReactNode } from 'react'
import { observer } from 'mobx-react-lite'
import { useSessionEditor } from '../hooks/useSessionEditor'
import { WindowTabs } from './WindowTabs'
import { SettingsIcon } from '@/components/Icons'
import { Tooltip } from '@/components/Tooltip'
import { Chip } from '@/components/Chip'
import styles from './TerminalMonitor.module.css'

interface TerminalMonitorProps {
  children: ReactNode
  onSettingsClick?: () => void
}

export const TerminalMonitor = observer(function TerminalMonitor({
  children,
  onSettingsClick,
}: TerminalMonitorProps) {
  const store = useSessionEditor()
  // Show the multiplexer session name directly (e.g., "project_session-name")
  const displayName = store.session.sessionName || store.session.name || 'session'

  return (
    <div className={styles.monitor}>
      <div className={styles.titleBar}>
        <div className={styles.windowControls}>
          <span className={styles.controlRed} />
          <span className={styles.controlYellow} />
          <span className={styles.controlGreen} />
        </div>
        <span className={styles.titleText}>Session Preview</span>
        <div className={styles.spacer} />
      </div>
      <div className={styles.windowBar}>
        <WindowTabs />
        <div className={styles.windowBarRight}>
          <Chip>{displayName}</Chip>
          <Tooltip content="Session settings">
            <button className={styles.settingsButton} onClick={onSettingsClick}>
              <SettingsIcon size={16} />
            </button>
          </Tooltip>
        </div>
      </div>
      <div className={styles.screen}>{children}</div>
    </div>
  )
})
