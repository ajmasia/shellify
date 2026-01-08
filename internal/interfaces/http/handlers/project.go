package handlers

import (
	"net/http"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/domain"
)

// ProjectHandler handles project-related HTTP requests.
type ProjectHandler struct {
	projectSvc *application.ProjectService
	sessionSvc *application.SessionService
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(
	projectSvc *application.ProjectService,
	sessionSvc *application.SessionService,
) *ProjectHandler {
	return &ProjectHandler{
		projectSvc: projectSvc,
		sessionSvc: sessionSvc,
	}
}

// ProjectWithCount extends Project with session count.
type ProjectWithCount struct {
	domain.Project
	SessionCount int `json:"sessionCount"`
}

// CreateProjectRequest is the request body for creating a project.
type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateProjectRequest is the request body for updating a project.
type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// RestoreRequest is the request body for restoring sessions.
type RestoreRequest struct {
	Sessions []domain.Session `json:"sessions"`
}

// List returns all projects with session counts.
func (h *ProjectHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := h.projectSvc.ListProjects()
		if err != nil {
			MapDomainError(w, err)
			return
		}

		result := make([]ProjectWithCount, len(projects))
		for i, p := range projects {
			count, _ := h.projectSvc.CountSessions(p.ID)
			result[i] = ProjectWithCount{
				Project:      p,
				SessionCount: count,
			}
		}

		Success(w, result)
	}
}

// Get returns a single project.
func (h *ProjectHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := URLParam(r, "projectId")

		project, err := h.projectSvc.GetProject(projectID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		count, _ := h.projectSvc.CountSessions(project.ID)

		Success(w, ProjectWithCount{
			Project:      project,
			SessionCount: count,
		})
	}
}

// Create creates a new project.
func (h *ProjectHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProjectRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "invalid JSON")
			return
		}

		project, err := h.projectSvc.CreateProject(req.Name, req.Description)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Created(w, project)
	}
}

// Update updates an existing project.
func (h *ProjectHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := URLParam(r, "projectId")

		var req UpdateProjectRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "invalid JSON")
			return
		}

		project, err := h.projectSvc.UpdateProject(projectID, req.Name, req.Description)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, project)
	}
}

// Delete deletes a project.
func (h *ProjectHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := URLParam(r, "projectId")
		force := r.URL.Query().Get("force") == "true"

		if err := h.projectSvc.DeleteProject(projectID, force); err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, map[string]string{"deleted": projectID})
	}
}

// ListSessions returns all sessions for a project.
func (h *ProjectHandler) ListSessions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := URLParam(r, "projectId")

		// Verify project exists
		project, err := h.projectSvc.GetProject(projectID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		sessions, err := h.sessionSvc.ListSessions(project.ID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, sessions)
	}
}

// Backup creates a backup of a project with all its sessions.
func (h *ProjectHandler) Backup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := URLParam(r, "projectId")

		project, err := h.projectSvc.GetProject(projectID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		sessions, err := h.sessionSvc.ListSessions(project.ID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		// Get full session data
		var fullSessions []domain.Session
		for _, s := range sessions {
			session, err := h.sessionSvc.GetSession(project.ID, s.ID)
			if err != nil {
				continue
			}
			fullSessions = append(fullSessions, session)
		}

		backup := map[string]any{
			"project":  project,
			"sessions": fullSessions,
		}

		Success(w, backup)
	}
}

// Restore restores sessions into a project.
func (h *ProjectHandler) Restore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := URLParam(r, "projectId")

		// Verify project exists
		project, err := h.projectSvc.GetProject(projectID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		var req RestoreRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "invalid JSON")
			return
		}

		var restored []string
		for _, session := range req.Sessions {
			input := application.CreateSessionInput{
				Name:             session.Name,
				Description:      session.Description,
				WorkingDirectory: session.WorkingDirectory,
				Multiplexer:      session.TargetMultiplexer,
				Environment:      session.Environment,
				PreCommands:      session.PreCommands,
				PostCommands:     session.PostCommands,
				Windows:          convertWindows(session.Windows),
			}

			created, err := h.sessionSvc.CreateSession(project.ID, input)
			if err != nil {
				continue
			}
			restored = append(restored, created.ID)
		}

		Success(w, map[string]any{
			"restored": len(restored),
			"ids":      restored,
		})
	}
}

// convertWindows converts domain.Window slice to WindowInput slice.
func convertWindows(windows []domain.Window) []application.WindowInput {
	result := make([]application.WindowInput, len(windows))
	for i, w := range windows {
		result[i] = application.WindowInput{
			Name:             w.Name,
			WorkingDirectory: w.WorkingDirectory,
			RootDirection:    w.RootDirection,
			Panes:            convertPanes(w.Panes),
		}
	}
	return result
}

// convertPanes converts domain.Pane slice to PaneInput slice.
func convertPanes(panes []domain.Pane) []application.PaneInput {
	result := make([]application.PaneInput, len(panes))
	for i, p := range panes {
		result[i] = application.PaneInput{
			Name:             p.Name,
			WorkingDirectory: p.WorkingDirectory,
			Command:          p.Command,
			Size:             p.Size,
			Direction:        string(p.Direction),
		}
	}
	return result
}
