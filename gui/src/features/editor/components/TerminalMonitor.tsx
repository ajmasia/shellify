import type { ReactNode } from 'react'
import { WindowTabs } from './WindowTabs'
import styles from './TerminalMonitor.module.css'

interface TerminalMonitorProps {
  children: ReactNode
}

export function TerminalMonitor({ children }: TerminalMonitorProps) {
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
      <WindowTabs />
      <div className={styles.screen}>{children}</div>
    </div>
  )
}
