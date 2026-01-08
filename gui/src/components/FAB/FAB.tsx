import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { PlusIcon } from '@/components/Icons'
import { Tooltip } from '@/components/Tooltip'
import styles from './FAB.module.css'

interface FABProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon?: ReactNode
  tooltip?: string
}

export function FAB({ icon, tooltip, className, ...props }: FABProps) {
  const classes = [styles.fab, className].filter(Boolean).join(' ')

  const button = (
    <button className={classes} {...props}>
      {icon ?? <PlusIcon size={28} />}
    </button>
  )

  if (tooltip) {
    return (
      <Tooltip content={tooltip} side="left">
        {button}
      </Tooltip>
    )
  }

  return button
}
