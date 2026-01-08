import type { ReactNode } from 'react'
import { Header } from './Header'
import styles from './MainLayout.module.css'

interface MainLayoutProps {
  children: ReactNode
}

export function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className={styles.layout}>
      <Header />
      <main className={styles.main}>
        <span className={styles.watermark}>Shellify</span>
        {children}
      </main>
    </div>
  )
}
