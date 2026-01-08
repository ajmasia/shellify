package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPaths(t *testing.T) {
	t.Run("creates paths successfully", func(t *testing.T) {
		paths, err := NewPaths()
		require.NoError(t, err)
		assert.NotEmpty(t, paths.ConfigDir)
		assert.NotEmpty(t, paths.ConfigFile)
		assert.NotEmpty(t, paths.ProjectsDir)
	})
}

func TestResolveConfigDir(t *testing.T) {
	t.Run("uses SHELLIFY_CONFIG_DIR when set", func(t *testing.T) {
		t.Setenv("SHELLIFY_CONFIG_DIR", "/custom/path")
		t.Setenv("XDG_CONFIG_HOME", "/xdg/path")

		result := resolveConfigDir()
		assert.Equal(t, "/custom/path", result)
	})

	t.Run("uses XDG_CONFIG_HOME when SHELLIFY_CONFIG_DIR not set", func(t *testing.T) {
		t.Setenv("SHELLIFY_CONFIG_DIR", "")
		t.Setenv("XDG_CONFIG_HOME", "/xdg/path")

		// Unset SHELLIFY_CONFIG_DIR completely
		err := os.Unsetenv("SHELLIFY_CONFIG_DIR")
		require.NoError(t, err)

		result := resolveConfigDir()
		assert.Equal(t, "/xdg/path/shellify", result)
	})

	t.Run("falls back to ~/.config when neither set", func(t *testing.T) {
		t.Setenv("SHELLIFY_CONFIG_DIR", "")
		t.Setenv("XDG_CONFIG_HOME", "")

		err := os.Unsetenv("SHELLIFY_CONFIG_DIR")
		require.NoError(t, err)
		err = os.Unsetenv("XDG_CONFIG_HOME")
		require.NoError(t, err)

		result := resolveConfigDir()
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, ".config", "shellify")
		assert.Equal(t, expected, result)
	})
}

func TestPaths_PathMethods(t *testing.T) {
	paths := &Paths{
		ConfigDir:   "/config",
		ConfigFile:  "/config/config.json",
		ProjectsDir: "/config/projects",
	}

	t.Run("ProjectDir", func(t *testing.T) {
		result := paths.ProjectDir("proj-123")
		assert.Equal(t, "/config/projects/proj-123", result)
	})

	t.Run("ProjectFile", func(t *testing.T) {
		result := paths.ProjectFile("proj-123")
		assert.Equal(t, "/config/projects/proj-123/project.json", result)
	})

	t.Run("SessionsDir", func(t *testing.T) {
		result := paths.SessionsDir("proj-123")
		assert.Equal(t, "/config/projects/proj-123/sessions", result)
	})

	t.Run("SessionFile", func(t *testing.T) {
		result := paths.SessionFile("proj-123", "sess-456")
		assert.Equal(t, "/config/projects/proj-123/sessions/sess-456.json", result)
	})
}

func TestEnsureDir(t *testing.T) {
	t.Run("creates nested directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedDir := filepath.Join(tmpDir, "a", "b", "c")

		err := EnsureDir(nestedDir)
		require.NoError(t, err)

		info, err := os.Stat(nestedDir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("succeeds if directory exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		err := EnsureDir(tmpDir)
		assert.NoError(t, err)
	})
}
