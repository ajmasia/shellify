import type { ButtonHTMLAttributes, ReactNode } from 'react'
import styles from './IconButton.module.css'

export type IconButtonVariant = 'default' | 'danger' | 'ghost' | 'primary'
type IconButtonSize = 'small' | 'medium' | 'large'

interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: IconButtonVariant
  size?: IconButtonSize
  children: ReactNode
  'aria-label'?: string
}

export function IconButton({
  variant = 'default',
  size = 'medium',
  children,
  className,
  ...props
}: IconButtonProps) {
  const classes = [styles.iconButton, styles[variant], styles[size], className]
    .filter(Boolean)
    .join(' ')

  return (
    <button className={classes} {...props}>
      {children}
    </button>
  )
}
