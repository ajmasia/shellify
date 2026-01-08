package http

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ajmasia/shellify/internal/interfaces/http/handlers"
)

func (s *Server) setupRouter() {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(CORSMiddleware())

	// API routes
	r.Route("/api", s.apiRoutes)

	// Static file serving
	if s.config.StaticDir != "" {
		s.setupStaticHandler(r)
	}

	s.router = r
}

func (s *Server) apiRoutes(r chi.Router) {
	projectHandler := handlers.NewProjectHandler(s.projectSvc, s.sessionSvc)
	sessionHandler := handlers.NewSessionHandler(s.sessionSvc, s.launcherSvc)
	settingsHandler := handlers.NewSettingsHandler(s.storage)

	// API Documentation
	r.Get("/docs", handlers.APIDocs())

	// Health check
	r.Get("/health", handlers.HealthCheck())

	// Projects
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", projectHandler.List())
		r.Post("/", projectHandler.Create())
		r.Route("/{projectId}", func(r chi.Router) {
			r.Get("/", projectHandler.Get())
			r.Put("/", projectHandler.Update())
			r.Delete("/", projectHandler.Delete())
			r.Get("/sessions", projectHandler.ListSessions())
			r.Post("/backup", projectHandler.Backup())
			r.Post("/restore", projectHandler.Restore())
		})
	})

	// Sessions
	r.Route("/sessions", func(r chi.Router) {
		r.Get("/", sessionHandler.ListAll())
		r.Post("/", sessionHandler.Create())
		r.Get("/running", sessionHandler.ListRunning())
		r.Route("/{sessionId}", func(r chi.Router) {
			r.Get("/", sessionHandler.Get())
			r.Put("/", sessionHandler.Update())
			r.Delete("/", sessionHandler.Delete())
			r.Post("/clone", sessionHandler.Clone())
			r.Post("/launch", sessionHandler.Launch())
			r.Post("/attach", sessionHandler.Attach())
			r.Post("/stop", sessionHandler.Stop())
			r.Get("/status", sessionHandler.Status())
		})
	})

	// Settings
	r.Route("/settings", func(r chi.Router) {
		r.Get("/", settingsHandler.Get())
		r.Put("/", settingsHandler.Update())
	})
}

func (s *Server) setupStaticHandler(r *chi.Mux) {
	staticDir := s.config.StaticDir

	// Check if directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: Static directory %s does not exist", staticDir)
		return
	}

	// Create file server
	fileServer := http.FileServer(http.Dir(staticDir))

	// Serve static assets
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/", fileServer).ServeHTTP(w, r)
	})

	// Serve favicon and other root static files
	r.Get("/vite.svg", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	// SPA catch-all handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// Don't serve index.html for API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Check if the file exists
		filePath := filepath.Join(staticDir, r.URL.Path)
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Serve index.html for SPA routing
		indexPath := filepath.Join(staticDir, "index.html")
		http.ServeFile(w, r, indexPath)
	})
}
