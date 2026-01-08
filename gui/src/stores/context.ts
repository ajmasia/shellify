import { createContext } from 'react'
import type { UIStore } from './UIStore'

export interface StoreContextValue {
  uiStore: UIStore
}

export const StoreContext = createContext<StoreContextValue | null>(null)
