import { useState, type ReactNode } from 'react'
import { SessionEditorContext } from './context'
import { SessionEditorStore } from './SessionEditorStore'

interface SessionEditorProviderProps {
  children: ReactNode
}

export function SessionEditorProvider({ children }: SessionEditorProviderProps) {
  const [store] = useState(() => new SessionEditorStore())

  return <SessionEditorContext.Provider value={store}>{children}</SessionEditorContext.Provider>
}
