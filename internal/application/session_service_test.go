package application

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

// mockSessionRepository is a test double for domain.SessionRepository.
type mockSessionRepository struct {
	sessions map[string]map[string]domain.Session // projectID -> sessionID -> Session
}

func newMockSessionRepo() *mockSessionRepository {
	return &mockSessionRepository{
		sessions: make(map[string]map[string]domain.Session),
	}
}

func (m *mockSessionRepository) ensureProject(projectID string) {
	if m.sessions[projectID] == nil {
		m.sessions[projectID] = make(map[string]domain.Session)
	}
}

func (m *mockSessionRepository) ListSessions(projectID string) ([]domain.SessionMetadata, error) {
	m.ensureProject(projectID)
	result := make([]domain.SessionMetadata, 0)
	for _, s := range m.sessions[projectID] {
		result = append(result, s.ToMetadata())
	}
	return result, nil
}

func (m *mockSessionRepository) GetSession(projectID, sessionID string) (domain.Session, error) {
	m.ensureProject(projectID)
	if s, ok := m.sessions[projectID][sessionID]; ok {
		return s, nil
	}
	return domain.Session{}, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) GetSessionByName(projectID, name string) (domain.Session, error) {
	m.ensureProject(projectID)
	for _, s := range m.sessions[projectID] {
		if strings.EqualFold(s.Name, name) {
			return s, nil
		}
	}
	return domain.Session{}, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) CreateSession(projectID string, session domain.Session) (domain.Session, error) {
	m.ensureProject(projectID)
	m.sessions[projectID][session.ID] = session
	return session, nil
}

func (m *mockSessionRepository) UpdateSession(projectID string, session domain.Session) error {
	m.ensureProject(projectID)
	if _, ok := m.sessions[projectID][session.ID]; !ok {
		return domain.ErrSessionNotFound
	}
	m.sessions[projectID][session.ID] = session
	return nil
}

func (m *mockSessionRepository) DeleteSession(projectID, sessionID string) error {
	m.ensureProject(projectID)
	if _, ok := m.sessions[projectID][sessionID]; !ok {
		return domain.ErrSessionNotFound
	}
	delete(m.sessions[projectID], sessionID)
	return nil
}

func (m *mockSessionRepository) SessionExists(projectID, sessionID string) bool {
	m.ensureProject(projectID)
	_, ok := m.sessions[projectID][sessionID]
	return ok
}

func (m *mockSessionRepository) SessionExistsByName(projectID, name string) bool {
	m.ensureProject(projectID)
	for _, s := range m.sessions[projectID] {
		if strings.EqualFold(s.Name, name) {
			return true
		}
	}
	return false
}

// Helper to create a project in the mock
func createTestProject(repo *mockProjectRepository, name string) domain.Project {
	project := domain.Project{
		ID:            "proj-" + name,
		Name:          name,
		SessionPrefix: strings.ToLower(name),
	}
	repo.projects[project.ID] = project
	return project
}

func TestSessionService_CreateSession(t *testing.T) {
	t.Run("creates session successfully", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Description: "Development session",
			Multiplexer: domain.MultiplexerTmux,
			Windows: []WindowInput{
				{Name: "editor", Command: "nvim ."},
			},
		}

		session, err := svc.CreateSession(project.ID, input)
		require.NoError(t, err)
		assert.Equal(t, "dev", session.Name)
		assert.Equal(t, "testproject_dev", session.SessionName)
		assert.Equal(t, "Development session", session.Description)
		assert.Equal(t, domain.MultiplexerTmux, session.TargetMultiplexer)
		assert.Len(t, session.Windows, 1)
		assert.Equal(t, "editor", session.Windows[0].Name)
		assert.Len(t, session.Windows[0].Panes, 1)
		assert.Equal(t, "nvim .", session.Windows[0].Panes[0].Command)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		_, err := svc.CreateSession(project.ID, input)
		assert.ErrorIs(t, err, domain.ErrSessionNameEmpty)
	})

	t.Run("fails with invalid multiplexer", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerType("invalid"),
			Windows:     []WindowInput{{Name: "main"}},
		}

		_, err := svc.CreateSession(project.ID, input)
		assert.ErrorIs(t, err, domain.ErrInvalidMultiplexer)
	})

	t.Run("fails with no windows", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{},
		}

		_, err := svc.CreateSession(project.ID, input)
		assert.ErrorIs(t, err, domain.ErrNoWindows)
	})

	t.Run("fails with duplicate name", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		_, err := svc.CreateSession(project.ID, input)
		require.NoError(t, err)

		_, err = svc.CreateSession(project.ID, input)
		assert.ErrorIs(t, err, domain.ErrSessionAlreadyExists)
	})

	t.Run("fails with non-existent project", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		_, err := svc.CreateSession("non-existent", input)
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})

	t.Run("auto-generates window name if empty", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "", Command: "bash"}},
		}

		session, err := svc.CreateSession(project.ID, input)
		require.NoError(t, err)
		assert.Equal(t, "window-1", session.Windows[0].Name)
	})
}

func TestSessionService_GetSession(t *testing.T) {
	t.Run("gets session by ID", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)
		found, err := svc.GetSession(project.ID, created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("gets session by name", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)
		found, err := svc.GetSession(project.ID, "dev")

		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		_, err := svc.GetSession(project.ID, "non-existent")
		assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	})
}

func TestSessionService_UpdateSession(t *testing.T) {
	t.Run("updates session name", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)

		newName := "development"
		updated, err := svc.UpdateSession(project.ID, created.ID, UpdateSessionInput{Name: &newName})

		require.NoError(t, err)
		assert.Equal(t, "development", updated.Name)
		assert.Equal(t, "testproject_development", updated.SessionName)
	})

	t.Run("updates session description", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)

		newDesc := "Updated description"
		updated, err := svc.UpdateSession(project.ID, created.ID, UpdateSessionInput{Description: &newDesc})

		require.NoError(t, err)
		assert.Equal(t, "Updated description", updated.Description)
	})

	t.Run("updates multiplexer", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)

		newMux := domain.MultiplexerZellij
		updated, err := svc.UpdateSession(project.ID, created.ID, UpdateSessionInput{Multiplexer: &newMux})

		require.NoError(t, err)
		assert.Equal(t, domain.MultiplexerZellij, updated.TargetMultiplexer)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)

		emptyName := ""
		_, err := svc.UpdateSession(project.ID, created.ID, UpdateSessionInput{Name: &emptyName})

		assert.ErrorIs(t, err, domain.ErrSessionNameEmpty)
	})

	t.Run("fails with duplicate name", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		_, _ = svc.CreateSession(project.ID, CreateSessionInput{
			Name:        "session-a",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		})

		created, _ := svc.CreateSession(project.ID, CreateSessionInput{
			Name:        "session-b",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		})

		dupName := "session-a"
		_, err := svc.UpdateSession(project.ID, created.ID, UpdateSessionInput{Name: &dupName})

		assert.ErrorIs(t, err, domain.ErrSessionAlreadyExists)
	})

	t.Run("fails with invalid multiplexer", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)

		invalidMux := domain.MultiplexerType("invalid")
		_, err := svc.UpdateSession(project.ID, created.ID, UpdateSessionInput{Multiplexer: &invalidMux})

		assert.ErrorIs(t, err, domain.ErrInvalidMultiplexer)
	})
}

func TestSessionService_DeleteSession(t *testing.T) {
	t.Run("deletes session successfully", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		created, _ := svc.CreateSession(project.ID, input)
		err := svc.DeleteSession(project.ID, created.ID)

		require.NoError(t, err)
		_, err = svc.GetSession(project.ID, created.ID)
		assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	})

	t.Run("deletes session by name", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		input := CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		}

		_, _ = svc.CreateSession(project.ID, input)
		err := svc.DeleteSession(project.ID, "dev")

		require.NoError(t, err)
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project := createTestProject(projectRepo, "TestProject")

		svc := NewSessionService(sessionRepo, projectRepo)

		err := svc.DeleteSession(project.ID, "non-existent")
		assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	})
}

func TestSessionService_ListAllSessions(t *testing.T) {
	t.Run("lists sessions from all projects", func(t *testing.T) {
		projectRepo := newMockRepo()
		sessionRepo := newMockSessionRepo()
		project1 := createTestProject(projectRepo, "Project1")
		project2 := createTestProject(projectRepo, "Project2")

		svc := NewSessionService(sessionRepo, projectRepo)

		_, _ = svc.CreateSession(project1.ID, CreateSessionInput{
			Name:        "dev",
			Multiplexer: domain.MultiplexerTmux,
			Windows:     []WindowInput{{Name: "main"}},
		})

		_, _ = svc.CreateSession(project2.ID, CreateSessionInput{
			Name:        "test",
			Multiplexer: domain.MultiplexerZellij,
			Windows:     []WindowInput{{Name: "main"}},
		})

		all, err := svc.ListAllSessions()
		require.NoError(t, err)
		assert.Len(t, all, 2)
	})
}

func TestGenerateSessionName(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		session  string
		expected string
	}{
		{"simple", "proj", "dev", "proj_dev"},
		{"with spaces", "proj", "my session", "proj_my-session"},
		{"no prefix", "", "dev", "dev"},
		{"uppercase", "proj", "DEV", "proj_dev"},
		{"trim whitespace", "proj", "  dev  ", "proj_dev"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSessionName(tt.prefix, tt.session)
			assert.Equal(t, tt.expected, result)
		})
	}
}
