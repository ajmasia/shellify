package application

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ajmasia/shellify/internal/domain"
	"github.com/ajmasia/shellify/internal/infrastructure/generator"
)

// GeneratorService provides business logic for generating session scripts/layouts.
type GeneratorService struct {
	sessionRepo domain.SessionRepository
	projectRepo domain.ProjectRepository
}

// NewGeneratorService creates a new GeneratorService instance.
func NewGeneratorService(sessionRepo domain.SessionRepository, projectRepo domain.ProjectRepository) *GeneratorService {
	return &GeneratorService{
		sessionRepo: sessionRepo,
		projectRepo: projectRepo,
	}
}

// GenerateResult contains the result of generating a session script/layout.
type GenerateResult struct {
	Content       string
	FileExtension string
	Multiplexer   string
}

// GenerateSession generates a script or layout for a session.
func (s *GeneratorService) GenerateSession(projectID, sessionIDOrName string) (GenerateResult, error) {
	project, err := s.getProject(projectID)
	if err != nil {
		return GenerateResult{}, err
	}

	session, err := s.getSession(project.ID, sessionIDOrName)
	if err != nil {
		return GenerateResult{}, err
	}

	gen, err := generator.New(session.TargetMultiplexer)
	if err != nil {
		return GenerateResult{}, fmt.Errorf("creating generator: %w", err)
	}

	content, err := gen.Generate(session)
	if err != nil {
		return GenerateResult{}, fmt.Errorf("generating script: %w", err)
	}

	return GenerateResult{
		Content:       content,
		FileExtension: gen.FileExtension(),
		Multiplexer:   gen.Name(),
	}, nil
}

// GenerateSessionToFile generates a session script/layout and writes it to a file.
func (s *GeneratorService) GenerateSessionToFile(projectID, sessionIDOrName, outputPath string) (string, error) {
	result, err := s.GenerateSession(projectID, sessionIDOrName)
	if err != nil {
		return "", err
	}

	// Determine final output path
	finalPath := outputPath
	if finalPath == "" {
		// Get session for name
		project, _ := s.getProject(projectID)
		session, _ := s.getSession(project.ID, sessionIDOrName)
		finalPath = session.SessionName + result.FileExtension
	}

	// Ensure directory exists
	dir := filepath.Dir(finalPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("creating output directory: %w", err)
		}
	}

	// Write file
	if err := os.WriteFile(finalPath, []byte(result.Content), 0644); err != nil {
		return "", fmt.Errorf("writing file: %w", err)
	}

	// Make shell scripts executable
	if result.FileExtension == ".sh" {
		if err := os.Chmod(finalPath, 0755); err != nil {
			return "", fmt.Errorf("setting executable permission: %w", err)
		}
	}

	return finalPath, nil
}

// getProject resolves a project by ID or name.
func (s *GeneratorService) getProject(idOrName string) (domain.Project, error) {
	project, err := s.projectRepo.GetProject(idOrName)
	if err == nil {
		return project, nil
	}

	if err == domain.ErrProjectNotFound {
		return s.projectRepo.GetProjectByName(idOrName)
	}

	return domain.Project{}, err
}

// getSession resolves a session by ID or name within a project.
func (s *GeneratorService) getSession(projectID, idOrName string) (domain.Session, error) {
	session, err := s.sessionRepo.GetSession(projectID, idOrName)
	if err == nil {
		return session, nil
	}

	if err == domain.ErrSessionNotFound {
		return s.sessionRepo.GetSessionByName(projectID, idOrName)
	}

	return domain.Session{}, err
}
