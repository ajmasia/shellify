package storage

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ajmasia/shellify/internal/domain"
)

// Storage provides filesystem-based persistence for projects and sessions.
type Storage struct {
	paths *Paths
}

// NewStorage creates a new Storage instance and ensures required directories exist.
func NewStorage() (*Storage, error) {
	paths, err := NewPaths()
	if err != nil {
		return nil, fmt.Errorf("resolving paths: %w", err)
	}

	// Ensure base directories exist
	if err := EnsureDir(paths.ConfigDir); err != nil {
		return nil, fmt.Errorf("creating config directory: %w", err)
	}

	if err := EnsureDir(paths.ProjectsDir); err != nil {
		return nil, fmt.Errorf("creating projects directory: %w", err)
	}

	return &Storage{paths: paths}, nil
}

// Config operations

// LoadConfig loads the application configuration from disk.
// Returns default config if file doesn't exist.
func (s *Storage) LoadConfig() (domain.Config, error) {
	config, err := readJSON[domain.Config](s.paths.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.DefaultConfig(), nil
		}
		return domain.Config{}, fmt.Errorf("reading config: %w", err)
	}
	return config, nil
}

// SaveConfig saves the application configuration to disk.
func (s *Storage) SaveConfig(config domain.Config) error {
	if err := writeJSON(s.paths.ConfigFile, config); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// Project operations

// ListProjects returns all projects sorted by UpdatedAt (newest first).
func (s *Storage) ListProjects() ([]domain.Project, error) {
	entries, err := os.ReadDir(s.paths.ProjectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Project{}, nil
		}
		return nil, fmt.Errorf("reading projects directory: %w", err)
	}

	var projects []domain.Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		project, err := s.GetProject(entry.Name())
		if err != nil {
			// Skip invalid projects
			continue
		}
		projects = append(projects, project)
	}

	// Sort by UpdatedAt descending
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].UpdatedAt.After(projects[j].UpdatedAt)
	})

	return projects, nil
}

// GetProject retrieves a project by ID.
func (s *Storage) GetProject(projectID string) (domain.Project, error) {
	project, err := readJSON[domain.Project](s.paths.ProjectFile(projectID))
	if err != nil {
		if os.IsNotExist(err) {
			return domain.Project{}, domain.ErrProjectNotFound
		}
		return domain.Project{}, fmt.Errorf("reading project: %w", err)
	}
	return project, nil
}

// CreateProject creates a new project.
func (s *Storage) CreateProject(project domain.Project) (domain.Project, error) {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}

	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	// Create project directory
	if err := EnsureDir(s.paths.ProjectDir(project.ID)); err != nil {
		return domain.Project{}, fmt.Errorf("creating project directory: %w", err)
	}

	// Create sessions directory
	if err := EnsureDir(s.paths.SessionsDir(project.ID)); err != nil {
		return domain.Project{}, fmt.Errorf("creating sessions directory: %w", err)
	}

	// Write project file
	if err := writeJSON(s.paths.ProjectFile(project.ID), project); err != nil {
		return domain.Project{}, fmt.Errorf("writing project: %w", err)
	}

	return project, nil
}

// UpdateProject updates an existing project.
func (s *Storage) UpdateProject(project domain.Project) error {
	existing, err := s.GetProject(project.ID)
	if err != nil {
		return err
	}

	// Preserve CreatedAt
	project.CreatedAt = existing.CreatedAt
	project.UpdatedAt = time.Now()

	if err := writeJSON(s.paths.ProjectFile(project.ID), project); err != nil {
		return fmt.Errorf("writing project: %w", err)
	}

	return nil
}

// DeleteProject deletes a project and all its sessions.
func (s *Storage) DeleteProject(projectID string) error {
	projectDir := s.paths.ProjectDir(projectID)

	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return domain.ErrProjectNotFound
	}

	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("deleting project directory: %w", err)
	}

	return nil
}

// Session operations

// storedSession wraps a session with metadata for storage.
type storedSession struct {
	domain.SessionMetadata
	Session domain.Session `json:"session"`
}

// ListSessions returns all sessions for a project sorted by UpdatedAt (newest first).
func (s *Storage) ListSessions(projectID string) ([]domain.SessionMetadata, error) {
	sessionsDir := s.paths.SessionsDir(projectID)

	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.SessionMetadata{}, nil
		}
		return nil, fmt.Errorf("reading sessions directory: %w", err)
	}

	var sessions []domain.SessionMetadata
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		stored, err := readJSON[storedSession](s.paths.SessionFile(projectID, entry.Name()[:len(entry.Name())-5])) // Remove .json
		if err != nil {
			// Skip invalid sessions
			continue
		}
		sessions = append(sessions, stored.SessionMetadata)
	}

	// Sort by UpdatedAt descending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})

	return sessions, nil
}

// GetSession retrieves a session by project ID and session ID.
func (s *Storage) GetSession(projectID, sessionID string) (domain.Session, error) {
	stored, err := readJSON[storedSession](s.paths.SessionFile(projectID, sessionID))
	if err != nil {
		if os.IsNotExist(err) {
			return domain.Session{}, domain.ErrSessionNotFound
		}
		return domain.Session{}, fmt.Errorf("reading session: %w", err)
	}
	return stored.Session, nil
}

// CreateSession creates a new session in a project.
func (s *Storage) CreateSession(projectID string, session domain.Session) (domain.Session, error) {
	// Verify project exists
	if _, err := s.GetProject(projectID); err != nil {
		return domain.Session{}, err
	}

	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now

	stored := storedSession{
		SessionMetadata: session.ToMetadata(),
		Session:         session,
	}

	if err := writeJSON(s.paths.SessionFile(projectID, session.ID), stored); err != nil {
		return domain.Session{}, fmt.Errorf("writing session: %w", err)
	}

	return session, nil
}

// UpdateSession updates an existing session.
func (s *Storage) UpdateSession(projectID string, session domain.Session) error {
	existing, err := s.GetSession(projectID, session.ID)
	if err != nil {
		return err
	}

	// Preserve CreatedAt
	session.CreatedAt = existing.CreatedAt
	session.UpdatedAt = time.Now()

	stored := storedSession{
		SessionMetadata: session.ToMetadata(),
		Session:         session,
	}

	if err := writeJSON(s.paths.SessionFile(projectID, session.ID), stored); err != nil {
		return fmt.Errorf("writing session: %w", err)
	}

	return nil
}

// DeleteSession deletes a session from a project.
func (s *Storage) DeleteSession(projectID, sessionID string) error {
	sessionFile := s.paths.SessionFile(projectID, sessionID)

	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		return domain.ErrSessionNotFound
	}

	if err := os.Remove(sessionFile); err != nil {
		return fmt.Errorf("deleting session file: %w", err)
	}

	return nil
}

// Additional project operations for ProjectRepository interface

// GetProjectByName retrieves a project by name (case-insensitive).
func (s *Storage) GetProjectByName(name string) (domain.Project, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return domain.Project{}, err
	}

	for _, p := range projects {
		if strings.EqualFold(p.Name, name) {
			return p, nil
		}
	}

	return domain.Project{}, domain.ErrProjectNotFound
}

// ProjectExists checks if a project exists by ID.
func (s *Storage) ProjectExists(id string) bool {
	_, err := os.Stat(s.paths.ProjectFile(id))
	return err == nil
}

// ProjectExistsByName checks if a project with the given name exists (case-insensitive).
func (s *Storage) ProjectExistsByName(name string) bool {
	_, err := s.GetProjectByName(name)
	return err == nil
}

// CountSessions returns the number of sessions in a project.
func (s *Storage) CountSessions(projectID string) (int, error) {
	sessions, err := s.ListSessions(projectID)
	if err != nil {
		return 0, err
	}
	return len(sessions), nil
}
