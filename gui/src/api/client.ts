import type {
  Project,
  ProjectWithCount,
  CreateProjectRequest,
  UpdateProjectRequest,
  Session,
  CreateSessionRequest,
  UpdateSessionRequest,
} from '@/domain'

const API_BASE = import.meta.env.VITE_API_URL || '/api'

interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_BASE}${endpoint}`

  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })

  const json: ApiResponse<T> = await response.json()

  if (!response.ok || !json.success) {
    throw new Error(json.error || `Request failed: ${response.status}`)
  }

  return json.data as T
}

export const api = {
  // Health
  health: () => request<{ status: string }>('/health'),

  // Projects
  projects: {
    list: () => request<ProjectWithCount[]>('/projects'),

    get: (id: string) => request<ProjectWithCount>(`/projects/${id}`),

    create: (data: CreateProjectRequest) =>
      request<Project>('/projects', {
        method: 'POST',
        body: JSON.stringify(data),
      }),

    update: (id: string, data: UpdateProjectRequest) =>
      request<Project>(`/projects/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),

    delete: (id: string, force = false) =>
      request<{ deleted: string }>(`/projects/${id}${force ? '?force=true' : ''}`, {
        method: 'DELETE',
      }),

    sessions: (projectId: string) => request<Session[]>(`/projects/${projectId}/sessions`),
  },

  // Sessions
  sessions: {
    get: (id: string) => request<{ projectId: string; session: Session }>(`/sessions/${id}`),

    create: (data: CreateSessionRequest) =>
      request<{ projectId: string; session: Session }>('/sessions', {
        method: 'POST',
        body: JSON.stringify(data),
      }),

    update: (id: string, data: UpdateSessionRequest) =>
      request<Session>(`/sessions/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),

    delete: (id: string) =>
      request<{ deleted: string }>(`/sessions/${id}`, {
        method: 'DELETE',
      }),
  },
}
