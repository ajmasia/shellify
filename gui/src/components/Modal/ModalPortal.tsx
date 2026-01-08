import { useEffect, useState, type ReactNode } from 'react'
import { createPortal } from 'react-dom'

interface ModalPortalProps {
  children: ReactNode
}

export function ModalPortal({ children }: ModalPortalProps) {
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    return () => setMounted(false)
  }, [])

  if (!mounted) return null

  const container = document.getElementById('modal-root')
  if (!container) return null

  return createPortal(children, container)
}
