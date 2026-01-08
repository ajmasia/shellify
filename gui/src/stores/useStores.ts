import { useContext } from 'react'
import { StoreContext } from './context'
import type { UIStore } from './UIStore'

export function useStores() {
  const context = useContext(StoreContext)
  if (!context) {
    throw new Error('useStores must be used within a StoreProvider')
  }
  return context
}

export function useUIStore(): UIStore {
  const { uiStore } = useStores()
  return uiStore
}
