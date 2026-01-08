package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/domain"
)

// Response represents a standard API response.
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JSON sends a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Success sends a 200 OK response with data.
func Success(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response with data.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error sends an error response with the given status code.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

// BadRequest sends a 400 Bad Request error.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// NotFound sends a 404 Not Found error.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// Conflict sends a 409 Conflict error.
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, message)
}

// InternalError sends a 500 Internal Server Error.
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}

// DecodeJSON decodes JSON from request body into v.
func DecodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// URLParam extracts a URL parameter from the request.
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// MapDomainError maps domain errors to appropriate HTTP responses.
func MapDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProjectNotFound):
		NotFound(w, "project not found")
	case errors.Is(err, domain.ErrSessionNotFound):
		NotFound(w, "session not found")
	case errors.Is(err, domain.ErrProjectNameEmpty):
		BadRequest(w, "project name cannot be empty")
	case errors.Is(err, domain.ErrSessionNameEmpty):
		BadRequest(w, "session name cannot be empty")
	case errors.Is(err, domain.ErrProjectAlreadyExists):
		Conflict(w, "project with this name already exists")
	case errors.Is(err, domain.ErrSessionAlreadyExists):
		Conflict(w, "session with this name already exists")
	case errors.Is(err, domain.ErrInvalidMultiplexer):
		BadRequest(w, "invalid multiplexer type")
	case errors.Is(err, domain.ErrInvalidDirection):
		BadRequest(w, "invalid direction")
	case errors.Is(err, domain.ErrNoWindows):
		BadRequest(w, "session must have at least one window")
	default:
		var hasSessionsErr *application.ProjectHasSessionsError
		if errors.As(err, &hasSessionsErr) {
			Conflict(w, fmt.Sprintf("project has %d sessions, delete them first or use force=true", hasSessionsErr.Count))
			return
		}
		InternalError(w, err.Error())
	}
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheck returns a handler for the health check endpoint.
func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Success(w, HealthResponse{Status: "ok"})
	}
}
