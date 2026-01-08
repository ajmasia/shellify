package domain

// ProjectRepository defines the interface for project persistence operations.
type ProjectRepository interface {
	ListProjects() ([]Project, error)
	GetProject(id string) (Project, error)
	GetProjectByName(name string) (Project, error)
	CreateProject(project Project) (Project, error)
	UpdateProject(project Project) error
	DeleteProject(id string) error
	ProjectExists(id string) bool
	ProjectExistsByName(name string) bool
	CountSessions(projectID string) (int, error)
}

// SessionRepository defines the interface for session persistence operations.
type SessionRepository interface {
	ListSessions(projectID string) ([]SessionMetadata, error)
	GetSession(projectID, sessionID string) (Session, error)
	GetSessionByName(projectID, name string) (Session, error)
	CreateSession(projectID string, session Session) (Session, error)
	UpdateSession(projectID string, session Session) error
	DeleteSession(projectID, sessionID string) error
	SessionExists(projectID, sessionID string) bool
	SessionExistsByName(projectID, name string) bool
}
