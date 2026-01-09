import type { IconProps } from './types'

export function ShellifyIcon({ size, className, ...props }: IconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size ?? '100%'}
      height={size ?? '100%'}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.5"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      {...props}
    >
      {/* Terminal window frame */}
      <rect x="3" y="3" width="18" height="18" rx="2" />
      {/* Title bar */}
      <line x1="3" y1="7" x2="21" y2="7" />
      {/* Window controls */}
      <circle cx="6" cy="5" r="0.5" fill="currentColor" />
      <circle cx="8.5" cy="5" r="0.5" fill="currentColor" />
      <circle cx="11" cy="5" r="0.5" fill="currentColor" />
      {/* Vertical split */}
      <line x1="12" y1="7" x2="12" y2="21" />
      {/* Horizontal split on right pane */}
      <line x1="12" y1="14" x2="21" y2="14" />
      {/* Command prompt in left pane */}
      <path d="M5.5 10.5l1.5 1.5-1.5 1.5" />
    </svg>
  )
}
