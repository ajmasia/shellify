// Package storage provides filesystem-based persistence for projects and sessions.
package storage

import (
	"os"
	"path/filepath"
)

const (
	appName        = "shellify"
	configFileName = "config.json"
	projectsDir    = "projects"
	sessionsDir    = "sessions"
	projectFile    = "project.json"
)

// Paths holds the resolved paths for application data storage.
type Paths struct {
	ConfigDir   string
	ConfigFile  string
	ProjectsDir string
}

// NewPaths creates a new Paths instance with resolved directories.
// Priority: SHELLIFY_CONFIG_DIR > XDG_CONFIG_HOME > ~/.config
func NewPaths() (*Paths, error) {
	configDir := resolveConfigDir()

	p := &Paths{
		ConfigDir:   configDir,
		ConfigFile:  filepath.Join(configDir, configFileName),
		ProjectsDir: filepath.Join(configDir, projectsDir),
	}

	return p, nil
}

// resolveConfigDir determines the configuration directory based on environment.
func resolveConfigDir() string {
	// Check for explicit override
	if dir := os.Getenv("SHELLIFY_CONFIG_DIR"); dir != "" {
		return dir
	}

	// Check XDG_CONFIG_HOME
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, appName)
	}

	// Fall back to ~/.config
	home, err := os.UserHomeDir()
	if err != nil {
		// Last resort fallback
		return filepath.Join(".", ".config", appName)
	}

	return filepath.Join(home, ".config", appName)
}

// ProjectDir returns the directory path for a specific project.
func (p *Paths) ProjectDir(projectID string) string {
	return filepath.Join(p.ProjectsDir, projectID)
}

// ProjectFile returns the path to a project's metadata file.
func (p *Paths) ProjectFile(projectID string) string {
	return filepath.Join(p.ProjectDir(projectID), projectFile)
}

// SessionsDir returns the sessions directory for a specific project.
func (p *Paths) SessionsDir(projectID string) string {
	return filepath.Join(p.ProjectDir(projectID), sessionsDir)
}

// SessionFile returns the path to a session file.
func (p *Paths) SessionFile(projectID, sessionID string) string {
	return filepath.Join(p.SessionsDir(projectID), sessionID+".json")
}

// EnsureDir creates a directory and all parent directories if they don't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
