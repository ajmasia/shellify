package handlers

import (
	"net/http"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/domain"
)

// SessionHandler handles session-related HTTP requests.
type SessionHandler struct {
	sessionSvc  *application.SessionService
	launcherSvc *application.LauncherService
}

// NewSessionHandler creates a new SessionHandler.
func NewSessionHandler(
	sessionSvc *application.SessionService,
	launcherSvc *application.LauncherService,
) *SessionHandler {
	return &SessionHandler{
		sessionSvc:  sessionSvc,
		launcherSvc: launcherSvc,
	}
}

// SessionResponse is the response for a session with project ID.
type SessionResponse struct {
	ProjectID string         `json:"projectId"`
	Session   domain.Session `json:"session"`
}

// CreateSessionRequest is the request body for creating a session.
type CreateSessionRequest struct {
	ProjectID string                       `json:"projectId"`
	Session   CreateSessionSessionInput    `json:"session"`
}

// CreateSessionSessionInput is the session input for creation.
type CreateSessionSessionInput struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description,omitempty"`
	WorkingDirectory string                 `json:"workingDirectory,omitempty"`
	Multiplexer      domain.MultiplexerType `json:"targetMultiplexer"`
	Environment      map[string]string      `json:"environment,omitempty"`
	PreCommands      []string               `json:"preCommands,omitempty"`
	PostCommands     []string               `json:"postCommands,omitempty"`
	Windows          []WindowInput          `json:"windows"`
	DefaultWindowID  string                 `json:"defaultWindowId,omitempty"`
}

// WindowInput is the window input for creation.
type WindowInput struct {
	Name             string      `json:"name"`
	WorkingDirectory string      `json:"workingDirectory,omitempty"`
	Command          string      `json:"command,omitempty"`
	Panes            []PaneInput `json:"panes,omitempty"`
}

// PaneInput is the pane input for creation.
type PaneInput struct {
	Name             string  `json:"name,omitempty"`
	Command          string  `json:"command,omitempty"`
	WorkingDirectory string  `json:"workingDirectory,omitempty"`
	Size             float64 `json:"size,omitempty"`
}

// CloneSessionRequest is the request body for cloning a session.
type CloneSessionRequest struct {
	Name string `json:"name"`
}

// SessionStatusResponse is the response for session status.
type SessionStatusResponse struct {
	Running     bool   `json:"running"`
	Attached    bool   `json:"attached"`
	SessionName string `json:"sessionName"`
	Multiplexer string `json:"multiplexer"`
}

// ListAll returns all sessions across all projects.
func (h *SessionHandler) ListAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions, err := h.sessionSvc.ListAllSessions()
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, sessions)
	}
}

// Get returns a single session.
func (h *SessionHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		session, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, SessionResponse{
			ProjectID: projectID,
			Session:   session,
		})
	}
}

// Create creates a new session.
func (h *SessionHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSessionRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "invalid JSON")
			return
		}

		if req.ProjectID == "" {
			BadRequest(w, "projectId is required")
			return
		}

		input := application.CreateSessionInput{
			Name:             req.Session.Name,
			Description:      req.Session.Description,
			WorkingDirectory: req.Session.WorkingDirectory,
			Multiplexer:      req.Session.Multiplexer,
			Environment:      req.Session.Environment,
			PreCommands:      req.Session.PreCommands,
			PostCommands:     req.Session.PostCommands,
			Windows:          convertWindowInputs(req.Session.Windows),
			DefaultWindowID:  req.Session.DefaultWindowID,
		}

		session, err := h.sessionSvc.CreateSession(req.ProjectID, input)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Created(w, SessionResponse{
			ProjectID: req.ProjectID,
			Session:   session,
		})
	}
}

// Update updates a session.
func (h *SessionHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		// Find session to get project ID
		_, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		var req struct {
			Name             *string                 `json:"name,omitempty"`
			Description      *string                 `json:"description,omitempty"`
			WorkingDirectory *string                 `json:"workingDirectory,omitempty"`
			Multiplexer      *domain.MultiplexerType `json:"targetMultiplexer,omitempty"`
		}
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "invalid JSON")
			return
		}

		input := application.UpdateSessionInput{
			Name:             req.Name,
			Description:      req.Description,
			WorkingDirectory: req.WorkingDirectory,
			Multiplexer:      req.Multiplexer,
		}

		session, err := h.sessionSvc.UpdateSession(projectID, sessionID, input)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, SessionResponse{
			ProjectID: projectID,
			Session:   session,
		})
	}
}

// Delete deletes a session.
func (h *SessionHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		// Find session to get project ID
		_, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		if err := h.sessionSvc.DeleteSession(projectID, sessionID); err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, map[string]string{"deleted": sessionID})
	}
}

// Clone clones a session.
func (h *SessionHandler) Clone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		// Find session to get project ID
		_, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		var req CloneSessionRequest
		_ = DecodeJSON(r, &req)

		newName := req.Name
		if newName == "" {
			newName = "Copy"
		}

		cloned, err := h.sessionSvc.CloneSession(projectID, sessionID, newName)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Created(w, SessionResponse{
			ProjectID: projectID,
			Session:   cloned,
		})
	}
}

// ListRunning returns all running sessions.
func (h *SessionHandler) ListRunning() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		running, err := h.launcherSvc.ListRunningSessions()
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, running)
	}
}

// Launch launches a session in detached mode.
func (h *SessionHandler) Launch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		// Find session to get project ID
		session, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		if err := h.launcherSvc.LaunchSessionDetached(projectID, sessionID); err != nil {
			InternalError(w, err.Error())
			return
		}

		Success(w, map[string]any{
			"launched":    true,
			"sessionName": session.SessionName,
			"multiplexer": session.TargetMultiplexer,
		})
	}
}

// Stop stops a running session.
func (h *SessionHandler) Stop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		// Find session to get project ID
		session, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		if err := h.launcherSvc.StopSession(projectID, sessionID); err != nil {
			InternalError(w, err.Error())
			return
		}

		Success(w, map[string]any{
			"stopped":     true,
			"sessionName": session.SessionName,
		})
	}
}

// Attach attaches to a running session.
func (h *SessionHandler) Attach() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		// Find session to get project ID
		session, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		if err := h.launcherSvc.AttachSession(projectID, sessionID); err != nil {
			InternalError(w, err.Error())
			return
		}

		Success(w, map[string]any{
			"attached":    true,
			"sessionName": session.SessionName,
			"multiplexer": session.TargetMultiplexer,
		})
	}
}

// Status returns the status of a session.
func (h *SessionHandler) Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := URLParam(r, "sessionId")

		// Find session to get project ID
		_, projectID, err := h.sessionSvc.FindSessionByID(sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		status, err := h.launcherSvc.GetSessionStatus(projectID, sessionID)
		if err != nil {
			MapDomainError(w, err)
			return
		}

		Success(w, SessionStatusResponse{
			Running:     status.Running,
			Attached:    status.Attached,
			SessionName: status.SessionName,
			Multiplexer: string(status.Multiplexer),
		})
	}
}

// convertWindowInputs converts handler WindowInput to application WindowInput.
func convertWindowInputs(windows []WindowInput) []application.WindowInput {
	result := make([]application.WindowInput, len(windows))
	for i, w := range windows {
		result[i] = application.WindowInput{
			Name:             w.Name,
			WorkingDirectory: w.WorkingDirectory,
			Command:          w.Command,
			Panes:            convertPaneInputs(w.Panes),
		}
	}
	return result
}

// convertPaneInputs converts handler PaneInput to application PaneInput.
func convertPaneInputs(panes []PaneInput) []application.PaneInput {
	result := make([]application.PaneInput, len(panes))
	for i, p := range panes {
		result[i] = application.PaneInput{
			Name:             p.Name,
			Command:          p.Command,
			WorkingDirectory: p.WorkingDirectory,
			Size:             p.Size,
		}
	}
	return result
}
