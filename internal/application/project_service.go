package application

import (
	"strings"

	"github.com/ajmasia/shellify/internal/domain"
)

// ProjectService provides business logic for project operations.
type ProjectService struct {
	repo domain.ProjectRepository
}

// NewProjectService creates a new ProjectService instance.
func NewProjectService(repo domain.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// ListProjects returns all projects sorted by update time.
func (s *ProjectService) ListProjects() ([]domain.Project, error) {
	return s.repo.ListProjects()
}

// GetProject retrieves a project by ID or name.
func (s *ProjectService) GetProject(idOrName string) (domain.Project, error) {
	// Try by ID first
	project, err := s.repo.GetProject(idOrName)
	if err == nil {
		return project, nil
	}

	// If not found by ID, try by name
	if err == domain.ErrProjectNotFound {
		return s.repo.GetProjectByName(idOrName)
	}

	return domain.Project{}, err
}

// CreateProject creates a new project with validation.
func (s *ProjectService) CreateProject(name, description string) (domain.Project, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return domain.Project{}, domain.ErrProjectNameEmpty
	}

	if s.repo.ProjectExistsByName(name) {
		return domain.Project{}, domain.ErrProjectAlreadyExists
	}

	project := domain.Project{
		Name:          name,
		Description:   strings.TrimSpace(description),
		SessionPrefix: generatePrefix(name),
	}

	return s.repo.CreateProject(project)
}

// UpdateProject updates an existing project.
func (s *ProjectService) UpdateProject(idOrName string, name, description *string) (domain.Project, error) {
	project, err := s.GetProject(idOrName)
	if err != nil {
		return domain.Project{}, err
	}

	if name != nil {
		newName := strings.TrimSpace(*name)
		if newName == "" {
			return domain.Project{}, domain.ErrProjectNameEmpty
		}

		// Check if name is taken by another project
		if !strings.EqualFold(newName, project.Name) && s.repo.ProjectExistsByName(newName) {
			return domain.Project{}, domain.ErrProjectAlreadyExists
		}

		project.Name = newName
		project.SessionPrefix = generatePrefix(newName)
	}

	if description != nil {
		project.Description = strings.TrimSpace(*description)
	}

	if err := s.repo.UpdateProject(project); err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

// DeleteProject deletes a project. Returns error if project has sessions unless force is true.
func (s *ProjectService) DeleteProject(idOrName string, force bool) error {
	project, err := s.GetProject(idOrName)
	if err != nil {
		return err
	}

	if !force {
		count, err := s.repo.CountSessions(project.ID)
		if err != nil {
			return err
		}
		if count > 0 {
			return &ProjectHasSessionsError{Count: count}
		}
	}

	return s.repo.DeleteProject(project.ID)
}

// CountSessions returns the number of sessions in a project.
func (s *ProjectService) CountSessions(idOrName string) (int, error) {
	project, err := s.GetProject(idOrName)
	if err != nil {
		return 0, err
	}
	return s.repo.CountSessions(project.ID)
}

// ProjectHasSessionsError is returned when trying to delete a project with sessions.
type ProjectHasSessionsError struct {
	Count int
}

func (e *ProjectHasSessionsError) Error() string {
	return "project has sessions"
}

// generatePrefix creates a session prefix from a project name.
func generatePrefix(name string) string {
	var result []rune
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result = append(result, r)
		}
	}
	return string(result)
}
