import { useContext } from 'react'
import { SessionEditorContext } from '../stores/context'
import type { SessionEditorStore } from '../stores/SessionEditorStore'

export function useSessionEditor(): SessionEditorStore {
  const store = useContext(SessionEditorContext)
  if (!store) {
    throw new Error('useSessionEditor must be used within a SessionEditorProvider')
  }
  return store
}
