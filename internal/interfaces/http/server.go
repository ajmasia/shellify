package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

// Config holds HTTP server configuration.
type Config struct {
	Port      int
	Host      string
	StaticDir string
}

// DefaultConfig returns the default server configuration.
func DefaultConfig() Config {
	return Config{
		Port: 3777,
		Host: "localhost",
	}
}

// Server represents the HTTP server.
type Server struct {
	config      Config
	router      *chi.Mux
	httpServer  *http.Server
	projectSvc  *application.ProjectService
	sessionSvc  *application.SessionService
	launcherSvc *application.LauncherService
	storage     *storage.Storage
}

// New creates a new Server instance.
func New(
	config Config,
	projectSvc *application.ProjectService,
	sessionSvc *application.SessionService,
	launcherSvc *application.LauncherService,
	store *storage.Storage,
) *Server {
	s := &Server{
		config:      config,
		projectSvc:  projectSvc,
		sessionSvc:  sessionSvc,
		launcherSvc: launcherSvc,
		storage:     store,
	}

	s.setupRouter()
	return s
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on http://%s", addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// Addr returns the server address.
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
}
