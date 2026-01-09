import { useNavigate, useLocation } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { IconButton } from '@/components/IconButton'
import { BackIcon, SettingsIcon, FolderIcon, TerminalIcon } from '@/components/Icons'
import { Tooltip } from '@/components/Tooltip'
import { useUIStore } from '@/stores'
import styles from './Header.module.css'

export const Header = observer(function Header() {
  const navigate = useNavigate()
  const location = useLocation()
  const uiStore = useUIStore()

  const isHome = location.pathname === '/'
  const projectName = uiStore.currentProjectName

  const handleBack = () => {
    uiStore.clearProject()
    navigate('/')
  }

  return (
    <header className={styles.header}>
      <div className={styles.left}>
        {!isHome && (
          <Tooltip content="Back to projects">
            <IconButton aria-label="Back" onClick={handleBack} variant="ghost">
              <BackIcon />
            </IconButton>
          </Tooltip>
        )}
        {isHome ? (
          <FolderIcon size={24} className={styles.titleIcon} />
        ) : (
          <TerminalIcon size={24} className={styles.titleIcon} />
        )}
        <h1 className={styles.title}>{isHome ? 'Projects' : (projectName ?? 'Sessions')}</h1>
      </div>
      <div className={styles.right}>
        <Tooltip content="Coming soon">
          <IconButton aria-label="Settings" disabled>
            <SettingsIcon size={20} />
          </IconButton>
        </Tooltip>
      </div>
    </header>
  )
})
