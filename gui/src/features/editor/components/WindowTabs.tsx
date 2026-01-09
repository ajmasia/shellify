import { observer } from 'mobx-react-lite'
import { useSessionEditor } from '../hooks/useSessionEditor'
import { CloseIcon } from '@/components/Icons'
import { Tooltip } from '@/components/Tooltip'
import styles from './WindowTabs.module.css'

export const WindowTabs = observer(function WindowTabs() {
  const store = useSessionEditor()

  const handleRemoveWindow = (e: React.MouseEvent, windowId: string) => {
    e.stopPropagation()
    store.removeWindow(windowId)
  }

  return (
    <div className={styles.tabs}>
      {store.session.windows.map((window) => (
        <div
          key={window.id}
          className={`${styles.tab} ${store.selectedWindowId === window.id ? styles.active : ''}`}
          onClick={() => store.selectWindow(window.id)}
        >
          <span className={styles.tabName}>{window.name}</span>
          {store.session.windows.length > 1 && (
            <Tooltip content="Close window">
              <button
                className={styles.closeButton}
                onClick={(e) => handleRemoveWindow(e, window.id)}
              >
                <CloseIcon />
              </button>
            </Tooltip>
          )}
        </div>
      ))}
    </div>
  )
})
