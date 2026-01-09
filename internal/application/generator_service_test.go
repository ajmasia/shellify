package application

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

func TestGeneratorService_GenerateSession(t *testing.T) {
	projectRepo := newMockRepo()
	sessionRepo := newMockSessionRepo()
	service := NewGeneratorService(sessionRepo, projectRepo)

	// Create a project
	project := createTestProject(projectRepo, "TestProject")

	// Create a tmux session
	session := domain.Session{
		ID:                "session-1",
		Name:              "dev",
		SessionName:       "testproject_dev",
		WorkingDirectory:  "~/projects",
		TargetMultiplexer: domain.MultiplexerTmux,
		Windows: []domain.Window{
			{
				ID:   "w1",
				Name: "editor",
				Panes: []domain.Pane{
					{ID: "p1", Name: "editor", Command: "nvim .", Size: 100},
				},
			},
		},
		DefaultWindowID: "w1",
	}
	sessionRepo.ensureProject(project.ID)
	sessionRepo.sessions[project.ID][session.ID] = session

	t.Run("generates tmux script", func(t *testing.T) {
		result, err := service.GenerateSession(project.ID, session.ID)
		require.NoError(t, err)

		assert.Equal(t, ".sh", result.FileExtension)
		assert.Equal(t, "tmux", result.Multiplexer)
		assert.Contains(t, result.Content, "#!/bin/bash")
		assert.Contains(t, result.Content, session.SessionName)
	})

	t.Run("generates for session by name", func(t *testing.T) {
		result, err := service.GenerateSession(project.Name, session.Name)
		require.NoError(t, err)

		assert.Contains(t, result.Content, "#!/bin/bash")
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		_, err := service.GenerateSession(project.ID, "non-existent")
		assert.Error(t, err)
	})

	t.Run("returns error for non-existent project", func(t *testing.T) {
		_, err := service.GenerateSession("non-existent", session.ID)
		assert.Error(t, err)
	})
}

func TestGeneratorService_GenerateSession_Zellij(t *testing.T) {
	projectRepo := newMockRepo()
	sessionRepo := newMockSessionRepo()
	service := NewGeneratorService(sessionRepo, projectRepo)

	// Create a project
	project := createTestProject(projectRepo, "TestProject")

	// Create a zellij session
	session := domain.Session{
		ID:                "session-1",
		Name:              "dev",
		SessionName:       "testproject_dev",
		WorkingDirectory:  "~/projects",
		TargetMultiplexer: domain.MultiplexerZellij,
		Windows: []domain.Window{
			{
				ID:   "w1",
				Name: "editor",
				Panes: []domain.Pane{
					{ID: "p1", Name: "editor", Command: "nvim .", Size: 100},
				},
			},
		},
		DefaultWindowID: "w1",
	}
	sessionRepo.ensureProject(project.ID)
	sessionRepo.sessions[project.ID][session.ID] = session

	result, err := service.GenerateSession(project.ID, session.ID)
	require.NoError(t, err)

	assert.Equal(t, ".sh", result.FileExtension)
	assert.Equal(t, "zellij", result.Multiplexer)
	assert.Contains(t, result.Content, "layout {")
}

func TestGeneratorService_GenerateSessionToFile(t *testing.T) {
	projectRepo := newMockRepo()
	sessionRepo := newMockSessionRepo()
	service := NewGeneratorService(sessionRepo, projectRepo)

	// Create a project
	project := createTestProject(projectRepo, "TestProject")

	// Create a tmux session
	session := domain.Session{
		ID:                "session-1",
		Name:              "dev",
		SessionName:       "testproject_dev",
		WorkingDirectory:  "~/projects",
		TargetMultiplexer: domain.MultiplexerTmux,
		Windows: []domain.Window{
			{
				ID:   "w1",
				Name: "editor",
				Panes: []domain.Pane{
					{ID: "p1", Name: "editor", Command: "nvim .", Size: 100},
				},
			},
		},
		DefaultWindowID: "w1",
	}
	sessionRepo.ensureProject(project.ID)
	sessionRepo.sessions[project.ID][session.ID] = session

	tmpDir := t.TempDir()

	t.Run("writes tmux script to file", func(t *testing.T) {
		outputPath := filepath.Join(tmpDir, "output", "test.sh")
		resultPath, err := service.GenerateSessionToFile(project.ID, session.ID, outputPath)
		require.NoError(t, err)

		assert.Equal(t, outputPath, resultPath)

		// Check file exists and content
		content, err := os.ReadFile(outputPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "#!/bin/bash")

		// Check executable permission
		info, err := os.Stat(outputPath)
		require.NoError(t, err)
		assert.True(t, info.Mode()&0100 != 0, "file should be executable")
	})

	t.Run("uses default filename when output path is empty", func(t *testing.T) {
		// Change to temp dir to avoid polluting cwd
		oldDir, _ := os.Getwd()
		require.NoError(t, os.Chdir(tmpDir))
		defer func() { _ = os.Chdir(oldDir) }()

		resultPath, err := service.GenerateSessionToFile(project.ID, session.ID, "")
		require.NoError(t, err)

		assert.Equal(t, session.SessionName+".sh", resultPath)

		// Cleanup
		_ = os.Remove(resultPath)
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		outputPath := filepath.Join(tmpDir, "output2", "test.sh")
		_, err := service.GenerateSessionToFile(project.ID, "non-existent", outputPath)
		assert.Error(t, err)
	})
}

func TestGeneratorService_GenerateSessionToFile_Zellij(t *testing.T) {
	projectRepo := newMockRepo()
	sessionRepo := newMockSessionRepo()
	service := NewGeneratorService(sessionRepo, projectRepo)

	// Create a project
	project := createTestProject(projectRepo, "TestProject")

	// Create a zellij session
	session := domain.Session{
		ID:                "session-1",
		Name:              "dev",
		SessionName:       "testproject_dev",
		WorkingDirectory:  "~/projects",
		TargetMultiplexer: domain.MultiplexerZellij,
		Windows: []domain.Window{
			{
				ID:   "w1",
				Name: "editor",
				Panes: []domain.Pane{
					{ID: "p1", Name: "editor", Command: "nvim .", Size: 100},
				},
			},
		},
		DefaultWindowID: "w1",
	}
	sessionRepo.ensureProject(project.ID)
	sessionRepo.sessions[project.ID][session.ID] = session

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output", "test.sh")
	resultPath, err := service.GenerateSessionToFile(project.ID, session.ID, outputPath)
	require.NoError(t, err)

	assert.Equal(t, outputPath, resultPath)

	// Check file exists and content
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "layout {")

	// Shell files should be executable (zellij now uses .sh wrapper)
	info, err := os.Stat(outputPath)
	require.NoError(t, err)
	assert.True(t, info.Mode()&0100 != 0, "sh file should be executable")
}
