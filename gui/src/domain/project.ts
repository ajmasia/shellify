export interface Project {
  id: string
  name: string
  description?: string
  sessionPrefix: string
  createdAt: string
  updatedAt: string
}

export interface ProjectWithCount extends Project {
  sessionCount: number
}

export interface CreateProjectRequest {
  name: string
  description?: string
}

export interface UpdateProjectRequest {
  name?: string
  description?: string
}
