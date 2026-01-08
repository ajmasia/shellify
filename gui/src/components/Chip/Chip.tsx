import type { ReactNode } from 'react'
import styles from './Chip.module.css'

type ChipVariant = 'default' | 'tmux' | 'zellij'

interface ChipProps {
  variant?: ChipVariant
  children: ReactNode
}

export function Chip({ variant = 'default', children }: ChipProps) {
  const classes = [styles.chip, styles[variant]].join(' ')

  return <span className={classes}>{children}</span>
}
