import { createContext } from 'react'
import type { SessionEditorStore } from './SessionEditorStore'

export const SessionEditorContext = createContext<SessionEditorStore | null>(null)
