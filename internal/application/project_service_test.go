package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

// mockProjectRepository is a test double for domain.ProjectRepository.
type mockProjectRepository struct {
	projects map[string]domain.Project
	sessions map[string]int
}

func newMockRepo() *mockProjectRepository {
	return &mockProjectRepository{
		projects: make(map[string]domain.Project),
		sessions: make(map[string]int),
	}
}

func (m *mockProjectRepository) ListProjects() ([]domain.Project, error) {
	result := make([]domain.Project, 0, len(m.projects))
	for _, p := range m.projects {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockProjectRepository) GetProject(id string) (domain.Project, error) {
	if p, ok := m.projects[id]; ok {
		return p, nil
	}
	return domain.Project{}, domain.ErrProjectNotFound
}

func (m *mockProjectRepository) GetProjectByName(name string) (domain.Project, error) {
	for _, p := range m.projects {
		if p.Name == name {
			return p, nil
		}
	}
	return domain.Project{}, domain.ErrProjectNotFound
}

func (m *mockProjectRepository) CreateProject(project domain.Project) (domain.Project, error) {
	project.ID = "test-id-" + project.Name
	m.projects[project.ID] = project
	return project, nil
}

func (m *mockProjectRepository) UpdateProject(project domain.Project) error {
	if _, ok := m.projects[project.ID]; !ok {
		return domain.ErrProjectNotFound
	}
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) DeleteProject(id string) error {
	if _, ok := m.projects[id]; !ok {
		return domain.ErrProjectNotFound
	}
	delete(m.projects, id)
	return nil
}

func (m *mockProjectRepository) ProjectExists(id string) bool {
	_, ok := m.projects[id]
	return ok
}

func (m *mockProjectRepository) ProjectExistsByName(name string) bool {
	for _, p := range m.projects {
		if p.Name == name {
			return true
		}
	}
	return false
}

func (m *mockProjectRepository) CountSessions(projectID string) (int, error) {
	return m.sessions[projectID], nil
}

func TestProjectService_CreateProject(t *testing.T) {
	t.Run("creates project successfully", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		project, err := svc.CreateProject("My Project", "A test project")
		require.NoError(t, err)
		assert.Equal(t, "My Project", project.Name)
		assert.Equal(t, "A test project", project.Description)
		assert.Equal(t, "myproject", project.SessionPrefix)
		assert.NotEmpty(t, project.ID)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		_, err := svc.CreateProject("", "description")
		assert.ErrorIs(t, err, domain.ErrProjectNameEmpty)
	})

	t.Run("fails with whitespace-only name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		_, err := svc.CreateProject("   ", "description")
		assert.ErrorIs(t, err, domain.ErrProjectNameEmpty)
	})

	t.Run("fails with duplicate name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		_, err := svc.CreateProject("Test", "first")
		require.NoError(t, err)

		_, err = svc.CreateProject("Test", "second")
		assert.ErrorIs(t, err, domain.ErrProjectAlreadyExists)
	})

	t.Run("trims whitespace from name and description", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		project, err := svc.CreateProject("  My Project  ", "  description  ")
		require.NoError(t, err)
		assert.Equal(t, "My Project", project.Name)
		assert.Equal(t, "description", project.Description)
	})
}

func TestProjectService_GetProject(t *testing.T) {
	t.Run("gets project by ID", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test", "")
		found, err := svc.GetProject(created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("gets project by name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test Project", "")
		found, err := svc.GetProject("Test Project")

		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("returns error for non-existent project", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		_, err := svc.GetProject("non-existent")
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	t.Run("updates project name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Old Name", "")
		newName := "New Name"
		updated, err := svc.UpdateProject(created.ID, &newName, nil)

		require.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, "newname", updated.SessionPrefix)
	})

	t.Run("updates project description", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test", "old desc")
		newDesc := "new description"
		updated, err := svc.UpdateProject(created.ID, nil, &newDesc)

		require.NoError(t, err)
		assert.Equal(t, "new description", updated.Description)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test", "")
		emptyName := ""
		_, err := svc.UpdateProject(created.ID, &emptyName, nil)

		assert.ErrorIs(t, err, domain.ErrProjectNameEmpty)
	})

	t.Run("fails with duplicate name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		_, _ = svc.CreateProject("Project A", "")
		created, _ := svc.CreateProject("Project B", "")

		dupName := "Project A"
		_, err := svc.UpdateProject(created.ID, &dupName, nil)

		assert.ErrorIs(t, err, domain.ErrProjectAlreadyExists)
	})

	t.Run("allows same name (no change)", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test", "")
		sameName := "Test"
		_, err := svc.UpdateProject(created.ID, &sameName, nil)

		assert.NoError(t, err)
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	t.Run("deletes project without sessions", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test", "")
		err := svc.DeleteProject(created.ID, false)

		require.NoError(t, err)
		_, err = svc.GetProject(created.ID)
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})

	t.Run("fails to delete project with sessions", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test", "")
		repo.sessions[created.ID] = 3

		err := svc.DeleteProject(created.ID, false)

		var sessionsErr *ProjectHasSessionsError
		assert.ErrorAs(t, err, &sessionsErr)
		assert.Equal(t, 3, sessionsErr.Count)
	})

	t.Run("force deletes project with sessions", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		created, _ := svc.CreateProject("Test", "")
		repo.sessions[created.ID] = 3

		err := svc.DeleteProject(created.ID, true)
		require.NoError(t, err)
	})

	t.Run("returns error for non-existent project", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewProjectService(repo)

		err := svc.DeleteProject("non-existent", false)
		assert.ErrorIs(t, err, domain.ErrProjectNotFound)
	})
}

func TestGeneratePrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "MyProject", "myproject"},
		{"with spaces", "My Project", "myproject"},
		{"with numbers", "Project123", "project123"},
		{"with special chars", "My-Project_v2!", "myprojectv2"},
		{"uppercase", "UPPERCASE", "uppercase"},
		{"mixed", "Hello World 2024!", "helloworld2024"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generatePrefix(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
