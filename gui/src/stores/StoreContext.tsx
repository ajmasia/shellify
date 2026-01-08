import { createContext, useContext, ReactNode } from 'react'
import { UIStore } from './UIStore'

interface StoreContextValue {
  uiStore: UIStore
}

const StoreContext = createContext<StoreContextValue | null>(null)

const uiStore = new UIStore()

interface StoreProviderProps {
  children: ReactNode
}

export function StoreProvider({ children }: StoreProviderProps) {
  return <StoreContext.Provider value={{ uiStore }}>{children}</StoreContext.Provider>
}

export function useStores(): StoreContextValue {
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
