package application

import (
	"fmt"

	"github.com/ajmasia/shellify/internal/domain"
	"github.com/ajmasia/shellify/internal/infrastructure/multiplexer"
)

// LauncherService provides business logic for launching, attaching, and stopping sessions.
type LauncherService struct {
	sessionRepo  domain.SessionRepository
	projectRepo  domain.ProjectRepository
	generatorSvc *GeneratorService
	launcher     *multiplexer.Launcher
}

// NewLauncherService creates a new LauncherService instance.
func NewLauncherService(
	sessionRepo domain.SessionRepository,
	projectRepo domain.ProjectRepository,
	generatorSvc *GeneratorService,
) (*LauncherService, error) {
	launcher, err := multiplexer.NewLauncher()
	if err != nil {
		return nil, fmt.Errorf("creating launcher: %w", err)
	}

	return &LauncherService{
		sessionRepo:  sessionRepo,
		projectRepo:  projectRepo,
		generatorSvc: generatorSvc,
		launcher:     launcher,
	}, nil
}

// SessionStatus holds the status of a session.
type SessionStatus struct {
	Running     bool
	Attached    bool
	SessionName string
	Multiplexer domain.MultiplexerType
}

// RunningSession represents a session with its running status.
type RunningSession struct {
	ProjectID   string
	ProjectName string
	Session     domain.SessionMetadata
	Status      SessionStatus
}

// LaunchSession launches a session by ID or name.
// This replaces the current process (use for CLI).
func (s *LauncherService) LaunchSession(projectID, sessionIDOrName string) error {
	session, err := s.resolveSession(projectID, sessionIDOrName)
	if err != nil {
		return err
	}

	avail := multiplexer.CheckMultiplexers()
	if !avail.HasMultiplexer(session.TargetMultiplexer) {
		return fmt.Errorf("%s is not installed.\n\n%s",
			session.TargetMultiplexer,
			multiplexer.GetInstallInstructions(session.TargetMultiplexer))
	}

	result, err := s.generatorSvc.GenerateSession(projectID, session.ID)
	if err != nil {
		return fmt.Errorf("generating session: %w", err)
	}

	switch session.TargetMultiplexer {
	case domain.MultiplexerTmux:
		return s.launcher.LaunchTmux(session.SessionName, result.Content)
	case domain.MultiplexerZellij:
		return s.launcher.LaunchZellij(session.SessionName, result.Content)
	default:
		return fmt.Errorf("unsupported multiplexer: %s", session.TargetMultiplexer)
	}
}

// LaunchSessionDetached launches a session in a new terminal window.
// Use for API/server where we don't want to replace the process.
func (s *LauncherService) LaunchSessionDetached(projectID, sessionIDOrName string) error {
	session, err := s.resolveSession(projectID, sessionIDOrName)
	if err != nil {
		return err
	}

	avail := multiplexer.CheckMultiplexers()
	if !avail.HasMultiplexer(session.TargetMultiplexer) {
		return fmt.Errorf("%s is not installed.\n\n%s",
			session.TargetMultiplexer,
			multiplexer.GetInstallInstructions(session.TargetMultiplexer))
	}

	result, err := s.generatorSvc.GenerateSession(projectID, session.ID)
	if err != nil {
		return fmt.Errorf("generating session: %w", err)
	}

	switch session.TargetMultiplexer {
	case domain.MultiplexerTmux:
		return s.launcher.LaunchTmuxDetached(session.SessionName, result.Content)
	case domain.MultiplexerZellij:
		return s.launcher.LaunchZellijDetached(session.SessionName, result.Content)
	default:
		return fmt.Errorf("unsupported multiplexer: %s", session.TargetMultiplexer)
	}
}

// AttachSession attaches to a running session.
// This replaces the current process.
func (s *LauncherService) AttachSession(projectID, sessionIDOrName string) error {
	session, err := s.resolveSession(projectID, sessionIDOrName)
	if err != nil {
		return err
	}

	if !s.launcher.IsRunning(session.SessionName, session.TargetMultiplexer) {
		return fmt.Errorf("session '%s' is not running", session.SessionName)
	}

	switch session.TargetMultiplexer {
	case domain.MultiplexerTmux:
		return s.launcher.AttachTmux(session.SessionName)
	case domain.MultiplexerZellij:
		return s.launcher.AttachZellij(session.SessionName)
	default:
		return fmt.Errorf("unsupported multiplexer: %s", session.TargetMultiplexer)
	}
}

// StopSession stops a running session.
func (s *LauncherService) StopSession(projectID, sessionIDOrName string) error {
	session, err := s.resolveSession(projectID, sessionIDOrName)
	if err != nil {
		return err
	}

	if !s.launcher.IsRunning(session.SessionName, session.TargetMultiplexer) {
		return fmt.Errorf("session '%s' is not running", session.SessionName)
	}

	switch session.TargetMultiplexer {
	case domain.MultiplexerTmux:
		return s.launcher.StopTmux(session.SessionName)
	case domain.MultiplexerZellij:
		return s.launcher.StopZellij(session.SessionName)
	default:
		return fmt.Errorf("unsupported multiplexer: %s", session.TargetMultiplexer)
	}
}

// GetSessionStatus returns the status of a session.
func (s *LauncherService) GetSessionStatus(projectID, sessionIDOrName string) (SessionStatus, error) {
	session, err := s.resolveSession(projectID, sessionIDOrName)
	if err != nil {
		return SessionStatus{}, err
	}

	return SessionStatus{
		Running:     s.launcher.IsRunning(session.SessionName, session.TargetMultiplexer),
		Attached:    s.launcher.IsAttached(session.SessionName, session.TargetMultiplexer),
		SessionName: session.SessionName,
		Multiplexer: session.TargetMultiplexer,
	}, nil
}

// ListRunningSessions returns all sessions with their running status.
func (s *LauncherService) ListRunningSessions() ([]RunningSession, error) {
	projects, err := s.projectRepo.ListProjects()
	if err != nil {
		return nil, err
	}

	var result []RunningSession
	for _, project := range projects {
		sessions, err := s.sessionRepo.ListSessions(project.ID)
		if err != nil {
			continue
		}

		for _, session := range sessions {
			fullSession, err := s.sessionRepo.GetSession(project.ID, session.ID)
			if err != nil {
				continue
			}

			running := s.launcher.IsRunning(fullSession.SessionName, fullSession.TargetMultiplexer)
			if running {
				result = append(result, RunningSession{
					ProjectID:   project.ID,
					ProjectName: project.Name,
					Session:     session,
					Status: SessionStatus{
						Running:     true,
						Attached:    s.launcher.IsAttached(fullSession.SessionName, fullSession.TargetMultiplexer),
						SessionName: fullSession.SessionName,
						Multiplexer: fullSession.TargetMultiplexer,
					},
				})
			}
		}
	}

	return result, nil
}

// resolveSession resolves a session by project ID and session ID or name.
func (s *LauncherService) resolveSession(projectID, sessionIDOrName string) (domain.Session, error) {
	project, err := s.getProject(projectID)
	if err != nil {
		return domain.Session{}, err
	}

	session, err := s.sessionRepo.GetSession(project.ID, sessionIDOrName)
	if err == nil {
		return session, nil
	}

	if err == domain.ErrSessionNotFound {
		return s.sessionRepo.GetSessionByName(project.ID, sessionIDOrName)
	}

	return domain.Session{}, err
}

// getProject resolves a project by ID or name.
func (s *LauncherService) getProject(idOrName string) (domain.Project, error) {
	project, err := s.projectRepo.GetProject(idOrName)
	if err == nil {
		return project, nil
	}

	if err == domain.ErrProjectNotFound {
		return s.projectRepo.GetProjectByName(idOrName)
	}

	return domain.Project{}, err
}
