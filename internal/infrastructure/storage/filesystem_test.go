package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

// setupTestStorage creates a storage instance with a temporary directory.
func setupTestStorage(t *testing.T) *Storage {
	t.Helper()

	tmpDir := t.TempDir()
	t.Setenv("SHELLIFY_CONFIG_DIR", tmpDir)

	storage, err := NewStorage()
	require.NoError(t, err)

	return storage
}

func TestNewStorage(t *testing.T) {
	t.Run("creates storage and directories", func(t *testing.T) {
		storage := setupTestStorage(t)
		assert.NotNil(t, storage)
		assert.DirExists(t, storage.paths.ConfigDir)
		assert.DirExists(t, storage.paths.ProjectsDir)
	})
}

func TestStorage_Config(t *testing.T) {
	storage := setupTestStorage(t)

	t.Run("returns default config when file doesn't exist", func(t *testing.T) {
		config, err := storage.LoadConfig()
		require.NoError(t, err)
		assert.Equal(t, domain.DefaultConfig(), config)
	})

	t.Run("saves and loads config", func(t *testing.T) {
		config := domain.Config{
			DefaultEditor:    "vim",
			DefaultTerminal:  "alacritty",
			CLIIntegration:   true,
			MaximizeTerminal: false,
			ServerPort:       8080,
		}

		err := storage.SaveConfig(config)
		require.NoError(t, err)

		loaded, err := storage.LoadConfig()
		require.NoError(t, err)
		assert.Equal(t, config, loaded)
	})
}

func TestStorage_Projects(t *testing.T) {
	storage := setupTestStorage(t)

	t.Run("list empty projects", func(t *testing.T) {
		projects, err := storage.ListProjects()
		require.NoError(t, err)
		assert.Empty(t, projects)
	})

	t.Run("create project", func(t *testing.T) {
		project := domain.Project{
			Name:          "test-project",
			Description:   "A test project",
			SessionPrefix: "test",
		}

		created, err := storage.CreateProject(project)
		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, project.Name, created.Name)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)
	})

	t.Run("get project", func(t *testing.T) {
		project := domain.Project{
			Name:          "get-test",
			SessionPrefix: "get",
		}

		created, err := storage.CreateProject(project)
		require.NoError(t, err)

		retrieved, err := storage.GetProject(created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Name, retrieved.Name)
	})

	t.Run("get non-existent project", func(t *testing.T) {
		_, err := storage.GetProject("non-existent-id")
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})

	t.Run("update project", func(t *testing.T) {
		project := domain.Project{
			Name:          "update-test",
			SessionPrefix: "upd",
		}

		created, err := storage.CreateProject(project)
		require.NoError(t, err)

		created.Name = "updated-name"
		created.Description = "Updated description"

		err = storage.UpdateProject(created)
		require.NoError(t, err)

		retrieved, err := storage.GetProject(created.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated-name", retrieved.Name)
		assert.Equal(t, "Updated description", retrieved.Description)
	})

	t.Run("delete project", func(t *testing.T) {
		project := domain.Project{
			Name:          "delete-test",
			SessionPrefix: "del",
		}

		created, err := storage.CreateProject(project)
		require.NoError(t, err)

		err = storage.DeleteProject(created.ID)
		require.NoError(t, err)

		_, err = storage.GetProject(created.ID)
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})

	t.Run("delete non-existent project", func(t *testing.T) {
		err := storage.DeleteProject("non-existent-id")
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})

	t.Run("list projects sorted by UpdatedAt", func(t *testing.T) {
		// Create multiple projects
		for _, name := range []string{"project-a", "project-b", "project-c"} {
			_, err := storage.CreateProject(domain.Project{
				Name:          name,
				SessionPrefix: name[:1],
			})
			require.NoError(t, err)
		}

		projects, err := storage.ListProjects()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(projects), 3)

		// Verify sorted by UpdatedAt (newest first)
		for i := 0; i < len(projects)-1; i++ {
			assert.True(t, projects[i].UpdatedAt.After(projects[i+1].UpdatedAt) ||
				projects[i].UpdatedAt.Equal(projects[i+1].UpdatedAt))
		}
	})
}

func TestStorage_Sessions(t *testing.T) {
	storage := setupTestStorage(t)

	// Create a project first
	project, err := storage.CreateProject(domain.Project{
		Name:          "session-test-project",
		SessionPrefix: "stp",
	})
	require.NoError(t, err)

	t.Run("list empty sessions", func(t *testing.T) {
		sessions, err := storage.ListSessions(project.ID)
		require.NoError(t, err)
		assert.Empty(t, sessions)
	})

	t.Run("create session", func(t *testing.T) {
		session := domain.Session{
			Name:              "test-session",
			SessionName:       "test-session",
			Description:       "A test session",
			TargetMultiplexer: domain.MultiplexerTmux,
			Windows: []domain.Window{
				{
					ID:   "win-1",
					Name: "main",
					Panes: []domain.Pane{
						{ID: "pane-1", Name: "shell"},
					},
				},
			},
		}

		created, err := storage.CreateSession(project.ID, session)
		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, session.Name, created.Name)
		assert.NotZero(t, created.CreatedAt)
	})

	t.Run("create session in non-existent project", func(t *testing.T) {
		session := domain.Session{Name: "test"}
		_, err := storage.CreateSession("non-existent", session)
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})

	t.Run("get session", func(t *testing.T) {
		session := domain.Session{
			Name:              "get-test-session",
			SessionName:       "get-test",
			TargetMultiplexer: domain.MultiplexerZellij,
		}

		created, err := storage.CreateSession(project.ID, session)
		require.NoError(t, err)

		retrieved, err := storage.GetSession(project.ID, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Name, retrieved.Name)
	})

	t.Run("get non-existent session", func(t *testing.T) {
		_, err := storage.GetSession(project.ID, "non-existent")
		assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	})

	t.Run("update session", func(t *testing.T) {
		session := domain.Session{
			Name:              "update-test-session",
			SessionName:       "update-test",
			TargetMultiplexer: domain.MultiplexerTmux,
		}

		created, err := storage.CreateSession(project.ID, session)
		require.NoError(t, err)

		created.Name = "updated-session"
		created.Description = "Updated"

		err = storage.UpdateSession(project.ID, created)
		require.NoError(t, err)

		retrieved, err := storage.GetSession(project.ID, created.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated-session", retrieved.Name)
	})

	t.Run("delete session", func(t *testing.T) {
		session := domain.Session{
			Name:              "delete-test-session",
			SessionName:       "delete-test",
			TargetMultiplexer: domain.MultiplexerTmux,
		}

		created, err := storage.CreateSession(project.ID, session)
		require.NoError(t, err)

		err = storage.DeleteSession(project.ID, created.ID)
		require.NoError(t, err)

		_, err = storage.GetSession(project.ID, created.ID)
		assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	})

	t.Run("delete non-existent session", func(t *testing.T) {
		err := storage.DeleteSession(project.ID, "non-existent")
		assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	})
}
