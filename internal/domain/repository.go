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
