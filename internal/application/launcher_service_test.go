package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

func TestLauncherService_GetSessionStatus(t *testing.T) {
	projectRepo := newMockRepo()
	sessionRepo := newMockSessionRepo()
	generatorSvc := NewGeneratorService(sessionRepo, projectRepo)

	project := createTestProject(projectRepo, "TestProject")

	// Create a session
	session := domain.Session{
		ID:                "session-1",
		Name:              "dev",
		SessionName:       "testproject_dev",
		TargetMultiplexer: domain.MultiplexerTmux,
		Windows: []domain.Window{
			{
				ID:   "w1",
				Name: "main",
				Panes: []domain.Pane{
					{ID: "p1", Name: "main", Size: 100},
				},
			},
		},
	}
	sessionRepo.ensureProject(project.ID)
	sessionRepo.sessions[project.ID][session.ID] = session

	launcherSvc, err := NewLauncherService(sessionRepo, projectRepo, generatorSvc)
	require.NoError(t, err)

	t.Run("returns status for existing session", func(t *testing.T) {
		status, err := launcherSvc.GetSessionStatus(project.ID, session.ID)
		require.NoError(t, err)

		assert.Equal(t, "testproject_dev", status.SessionName)
		assert.Equal(t, domain.MultiplexerTmux, status.Multiplexer)
		// Running state depends on actual tmux, will be false in test
		assert.False(t, status.Running)
		assert.False(t, status.Attached)
	})

	t.Run("returns status by session name", func(t *testing.T) {
		status, err := launcherSvc.GetSessionStatus(project.ID, "dev")
		require.NoError(t, err)

		assert.Equal(t, "testproject_dev", status.SessionName)
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		_, err := launcherSvc.GetSessionStatus(project.ID, "non-existent")
		assert.Error(t, err)
	})

	t.Run("returns error for non-existent project", func(t *testing.T) {
		_, err := launcherSvc.GetSessionStatus("non-existent", session.ID)
		assert.Error(t, err)
	})
}

func TestLauncherService_ListRunningSessions(t *testing.T) {
	projectRepo := newMockRepo()
	sessionRepo := newMockSessionRepo()
	generatorSvc := NewGeneratorService(sessionRepo, projectRepo)

	project1 := createTestProject(projectRepo, "Project1")
	project2 := createTestProject(projectRepo, "Project2")

	// Create sessions
	session1 := domain.Session{
		ID:                "session-1",
		Name:              "dev",
		SessionName:       "project1_dev",
		TargetMultiplexer: domain.MultiplexerTmux,
		Windows: []domain.Window{
			{ID: "w1", Name: "main", Panes: []domain.Pane{{ID: "p1", Size: 100}}},
		},
	}
	session2 := domain.Session{
		ID:                "session-2",
		Name:              "test",
		SessionName:       "project2_test",
		TargetMultiplexer: domain.MultiplexerZellij,
		Windows: []domain.Window{
			{ID: "w1", Name: "main", Panes: []domain.Pane{{ID: "p1", Size: 100}}},
		},
	}

	sessionRepo.ensureProject(project1.ID)
	sessionRepo.sessions[project1.ID][session1.ID] = session1
	sessionRepo.ensureProject(project2.ID)
	sessionRepo.sessions[project2.ID][session2.ID] = session2

	launcherSvc, err := NewLauncherService(sessionRepo, projectRepo, generatorSvc)
	require.NoError(t, err)

	t.Run("returns empty list when no sessions running", func(t *testing.T) {
		running, err := launcherSvc.ListRunningSessions()
		require.NoError(t, err)
		// No sessions are actually running in tests
		assert.Empty(t, running)
	})
}

func TestNewLauncherService(t *testing.T) {
	projectRepo := newMockRepo()
	sessionRepo := newMockSessionRepo()
	generatorSvc := NewGeneratorService(sessionRepo, projectRepo)

	svc, err := NewLauncherService(sessionRepo, projectRepo, generatorSvc)
	require.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestSessionStatus_Fields(t *testing.T) {
	status := SessionStatus{
		Running:     true,
		Attached:    true,
		SessionName: "test-session",
		Multiplexer: domain.MultiplexerTmux,
	}

	assert.True(t, status.Running)
	assert.True(t, status.Attached)
	assert.Equal(t, "test-session", status.SessionName)
	assert.Equal(t, domain.MultiplexerTmux, status.Multiplexer)
}

func TestRunningSession_Fields(t *testing.T) {
	running := RunningSession{
		ProjectID:   "proj-1",
		ProjectName: "TestProject",
		Session: domain.SessionMetadata{
			ID:                "session-1",
			Name:              "dev",
			TargetMultiplexer: domain.MultiplexerTmux,
		},
		Status: SessionStatus{
			Running:     true,
			Attached:    false,
			SessionName: "testproject_dev",
			Multiplexer: domain.MultiplexerTmux,
		},
	}

	assert.Equal(t, "proj-1", running.ProjectID)
	assert.Equal(t, "TestProject", running.ProjectName)
	assert.Equal(t, "session-1", running.Session.ID)
	assert.True(t, running.Status.Running)
}
