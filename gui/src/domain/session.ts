export type MultiplexerType = 'tmux' | 'zellij'
export type Direction = 'horizontal' | 'vertical'

export interface Pane {
  id: string
  name: string
  command?: string
  workingDirectory?: string
  environment?: Record<string, string>
  size: number
  children?: Pane[]
  direction?: Direction
}

export interface Window {
  id: string
  name: string
  workingDirectory?: string
  environment?: Record<string, string>
  panes: Pane[]
  rootDirection: Direction
}

export interface Session {
  id: string
  name: string
  sessionName: string
  description?: string
  workingDirectory?: string
  environment?: Record<string, string>
  preCommands?: string[]
  postCommands?: string[]
  windows: Window[]
  defaultWindowId?: string
  targetMultiplexer: MultiplexerType
  createdAt: string
  updatedAt: string
}

export interface SessionMetadata {
  id: string
  name: string
  description?: string
  targetMultiplexer: MultiplexerType
  createdAt: string
  updatedAt: string
}

export interface CreateSessionRequest {
  projectId: string
  session: {
    name: string
    description?: string
    workingDirectory?: string
    targetMultiplexer: MultiplexerType
    environment?: Record<string, string>
    preCommands?: string[]
    postCommands?: string[]
    windows: WindowInput[]
    defaultWindowId?: string
  }
}

export interface WindowInput {
  name: string
  workingDirectory?: string
  command?: string
  panes?: PaneInput[]
}

export interface PaneInput {
  name?: string
  command?: string
  workingDirectory?: string
  size?: number
}

export interface UpdateSessionRequest {
  name?: string
  description?: string
  workingDirectory?: string
  targetMultiplexer?: MultiplexerType
  environment?: Record<string, string>
  preCommands?: string[]
  postCommands?: string[]
  windows?: WindowInput[]
  defaultWindowId?: string
}
