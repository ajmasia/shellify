import type { ReactNode } from 'react'
import styles from './Chip.module.css'

type ChipVariant = 'default' | 'tmux' | 'zellij'

interface ChipProps {
  variant?: ChipVariant
  children: ReactNode
  selected?: boolean
  onClick?: () => void
}

export function Chip({ variant = 'default', children, selected, onClick }: ChipProps) {
  const isClickable = !!onClick
  const classes = [
    styles.chip,
    styles[variant],
    selected && styles.selected,
    isClickable && styles.clickable,
  ]
    .filter(Boolean)
    .join(' ')

  if (isClickable) {
    return (
      <button type="button" className={classes} onClick={onClick}>
        {children}
      </button>
    )
  }

  return <span className={classes}>{children}</span>
}
