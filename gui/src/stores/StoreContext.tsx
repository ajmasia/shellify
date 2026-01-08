import type { ReactNode } from 'react'
import { StoreContext } from './context'
import { UIStore } from './UIStore'

const uiStore = new UIStore()

interface StoreProviderProps {
  children: ReactNode
}

export function StoreProvider({ children }: StoreProviderProps) {
  return <StoreContext.Provider value={{ uiStore }}>{children}</StoreContext.Provider>
}
