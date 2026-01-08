package application

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ajmasia/shellify/internal/domain"
)

// SessionService provides business logic for session operations.
type SessionService struct {
	sessionRepo domain.SessionRepository
	projectRepo domain.ProjectRepository
}

// NewSessionService creates a new SessionService instance.
func NewSessionService(sessionRepo domain.SessionRepository, projectRepo domain.ProjectRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		projectRepo: projectRepo,
	}
}

// SessionWithProject represents a session with its project info for listing.
type SessionWithProject struct {
	ProjectID   string
	ProjectName string
	Session     domain.SessionMetadata
}

// CreateSessionInput contains the input for creating a new session.
type CreateSessionInput struct {
	Name             string
	Description      string
	WorkingDirectory string
	Multiplexer      domain.MultiplexerType
	Environment      map[string]string
	PreCommands      []string
	PostCommands     []string
	Windows          []WindowInput
	DefaultWindowID  string
}

// WindowInput contains the input for creating a window with panes.
type WindowInput struct {
	Name             string
	WorkingDirectory string
	Command          string
	RootDirection    domain.Direction
	Panes            []PaneInput
}

// PaneInput contains the input for creating a pane.
type PaneInput struct {
	Name             string
	Command          string
	WorkingDirectory string
	Size             float64
	Direction        string
}

// UpdateSessionInput contains optional fields for updating a session.
type UpdateSessionInput struct {
	Name             *string
	Description      *string
	WorkingDirectory *string
	Multiplexer      *domain.MultiplexerType
}

// ListSessions returns all sessions for a specific project.
func (s *SessionService) ListSessions(projectID string) ([]domain.SessionMetadata, error) {
	// Resolve project by ID or name
	project, err := s.getProject(projectID)
	if err != nil {
		return nil, err
	}
	return s.sessionRepo.ListSessions(project.ID)
}

// ListAllSessions returns sessions from all projects with project info.
func (s *SessionService) ListAllSessions() ([]SessionWithProject, error) {
	projects, err := s.projectRepo.ListProjects()
	if err != nil {
		return nil, err
	}

	var result []SessionWithProject
	for _, project := range projects {
		sessions, err := s.sessionRepo.ListSessions(project.ID)
		if err != nil {
			continue
		}
		for _, session := range sessions {
			result = append(result, SessionWithProject{
				ProjectID:   project.ID,
				ProjectName: project.Name,
				Session:     session,
			})
		}
	}

	return result, nil
}

// GetSession retrieves a session by ID or name within a project.
func (s *SessionService) GetSession(projectID, idOrName string) (domain.Session, error) {
	project, err := s.getProject(projectID)
	if err != nil {
		return domain.Session{}, err
	}

	// Try by ID first
	session, err := s.sessionRepo.GetSession(project.ID, idOrName)
	if err == nil {
		return session, nil
	}

	// If not found by ID, try by name
	if err == domain.ErrSessionNotFound {
		return s.sessionRepo.GetSessionByName(project.ID, idOrName)
	}

	return domain.Session{}, err
}

// CreateSession creates a new session with validation.
func (s *SessionService) CreateSession(projectID string, input CreateSessionInput) (domain.Session, error) {
	project, err := s.getProject(projectID)
	if err != nil {
		return domain.Session{}, err
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.Session{}, domain.ErrSessionNameEmpty
	}

	if s.sessionRepo.SessionExistsByName(project.ID, name) {
		return domain.Session{}, domain.ErrSessionAlreadyExists
	}

	if !input.Multiplexer.IsValid() {
		return domain.Session{}, domain.ErrInvalidMultiplexer
	}

	if len(input.Windows) == 0 {
		return domain.Session{}, domain.ErrNoWindows
	}

	// Build session
	session := domain.Session{
		ID:                uuid.New().String(),
		Name:              name,
		SessionName:       generateSessionName(project.SessionPrefix, name),
		Description:       strings.TrimSpace(input.Description),
		WorkingDirectory:  strings.TrimSpace(input.WorkingDirectory),
		Environment:       input.Environment,
		PreCommands:       input.PreCommands,
		PostCommands:      input.PostCommands,
		TargetMultiplexer: input.Multiplexer,
		Windows:           make([]domain.Window, len(input.Windows)),
	}

	// Build windows with pane trees
	for i, w := range input.Windows {
		windowName := strings.TrimSpace(w.Name)
		if windowName == "" {
			windowName = fmt.Sprintf("window-%d", i+1)
		}

		direction := w.RootDirection
		if direction == "" {
			direction = domain.DirectionHorizontal
		}

		window := domain.Window{
			ID:            fmt.Sprintf("w%d", i+1),
			Name:          windowName,
			RootDirection: direction,
		}

		// Build panes
		panes := make([]domain.Pane, len(w.Panes))
		sizePerPane := 100.0
		if len(w.Panes) > 0 {
			sizePerPane = 100.0 / float64(len(w.Panes))
		}

		for j, p := range w.Panes {
			paneName := strings.TrimSpace(p.Name)
			if paneName == "" {
				paneName = fmt.Sprintf("pane-%d", j+1)
			}

			panes[j] = domain.Pane{
				ID:               fmt.Sprintf("p%d", j+1),
				Name:             paneName,
				Command:          strings.TrimSpace(p.Command),
				WorkingDirectory: strings.TrimSpace(p.WorkingDirectory),
				Size:             sizePerPane,
			}
		}

		// Single pane: store directly
		// Multiple panes: wrap in root container
		if len(panes) == 1 {
			window.Panes = panes
		} else if len(panes) > 1 {
			rootPane := domain.Pane{
				ID:        "root",
				Size:      100,
				Direction: direction,
				Children:  panes,
			}
			window.Panes = []domain.Pane{rootPane}
		} else {
			// No panes provided, create default
			window.Panes = []domain.Pane{
				{
					ID:   "p1",
					Name: windowName,
					Size: 100,
				},
			}
		}

		session.Windows[i] = window
	}

	// Set default window
	if input.DefaultWindowID != "" {
		session.DefaultWindowID = input.DefaultWindowID
	} else if len(session.Windows) > 0 {
		session.DefaultWindowID = session.Windows[0].ID
	}

	return s.sessionRepo.CreateSession(project.ID, session)
}

// UpdateSession updates an existing session.
func (s *SessionService) UpdateSession(projectID, idOrName string, input UpdateSessionInput) (domain.Session, error) {
	project, err := s.getProject(projectID)
	if err != nil {
		return domain.Session{}, err
	}

	session, err := s.GetSession(project.ID, idOrName)
	if err != nil {
		return domain.Session{}, err
	}

	if input.Name != nil {
		newName := strings.TrimSpace(*input.Name)
		if newName == "" {
			return domain.Session{}, domain.ErrSessionNameEmpty
		}

		// Check if name is taken by another session
		if !strings.EqualFold(newName, session.Name) && s.sessionRepo.SessionExistsByName(project.ID, newName) {
			return domain.Session{}, domain.ErrSessionAlreadyExists
		}

		session.Name = newName
		session.SessionName = generateSessionName(project.SessionPrefix, newName)
	}

	if input.Description != nil {
		session.Description = strings.TrimSpace(*input.Description)
	}

	if input.WorkingDirectory != nil {
		session.WorkingDirectory = strings.TrimSpace(*input.WorkingDirectory)
	}

	if input.Multiplexer != nil {
		if !input.Multiplexer.IsValid() {
			return domain.Session{}, domain.ErrInvalidMultiplexer
		}
		session.TargetMultiplexer = *input.Multiplexer
	}

	if err := s.sessionRepo.UpdateSession(project.ID, session); err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

// DeleteSession deletes a session.
func (s *SessionService) DeleteSession(projectID, idOrName string) error {
	project, err := s.getProject(projectID)
	if err != nil {
		return err
	}

	session, err := s.GetSession(project.ID, idOrName)
	if err != nil {
		return err
	}

	return s.sessionRepo.DeleteSession(project.ID, session.ID)
}

// getProject resolves a project by ID or name.
func (s *SessionService) getProject(idOrName string) (domain.Project, error) {
	project, err := s.projectRepo.GetProject(idOrName)
	if err == nil {
		return project, nil
	}

	if err == domain.ErrProjectNotFound {
		return s.projectRepo.GetProjectByName(idOrName)
	}

	return domain.Project{}, err
}

// CloneSession creates a copy of an existing session with a new name.
func (s *SessionService) CloneSession(projectID, idOrName, newName string) (domain.Session, error) {
	project, err := s.getProject(projectID)
	if err != nil {
		return domain.Session{}, err
	}

	// Get the original session
	original, err := s.GetSession(project.ID, idOrName)
	if err != nil {
		return domain.Session{}, err
	}

	// Validate new name
	name := strings.TrimSpace(newName)
	if name == "" {
		name = original.Name + " (copy)"
	}

	if s.sessionRepo.SessionExistsByName(project.ID, name) {
		return domain.Session{}, domain.ErrSessionAlreadyExists
	}

	// Create cloned session
	cloned := original
	cloned.ID = uuid.New().String()
	cloned.Name = name
	cloned.SessionName = generateSessionName(project.SessionPrefix, name)

	return s.sessionRepo.CreateSession(project.ID, cloned)
}

// FindSessionByID finds a session by ID across all projects.
// Returns the session and the project ID it belongs to.
func (s *SessionService) FindSessionByID(sessionID string) (domain.Session, string, error) {
	projects, err := s.projectRepo.ListProjects()
	if err != nil {
		return domain.Session{}, "", err
	}

	for _, project := range projects {
		session, err := s.sessionRepo.GetSession(project.ID, sessionID)
		if err == nil {
			return session, project.ID, nil
		}
	}

	return domain.Session{}, "", domain.ErrSessionNotFound
}

// generateSessionName creates a session name from project prefix and session name.
func generateSessionName(prefix, name string) string {
	sessionName := strings.ToLower(strings.TrimSpace(name))
	sessionName = strings.ReplaceAll(sessionName, " ", "-")

	if prefix != "" {
		return prefix + "_" + sessionName
	}
	return sessionName
}
