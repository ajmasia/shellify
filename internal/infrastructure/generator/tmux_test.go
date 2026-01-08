package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

func TestTmuxGenerator_Generate(t *testing.T) {
	gen := &TmuxGenerator{}

	t.Run("generates script for simple session", func(t *testing.T) {
		session := createSimpleSession()
		output, err := gen.Generate(session)
		require.NoError(t, err)

		// Verify script header
		assert.Contains(t, output, "#!/bin/bash")
		assert.Contains(t, output, "# Session: proj_dev")
		assert.Contains(t, output, "set -e")

		// Verify working directory
		assert.Contains(t, output, `WORKING_DIR="~/projects"`)

		// Verify session creation
		assert.Contains(t, output, `tmux new-session -d -s "proj_dev" -n "editor"`)

		// Verify pane command
		assert.Contains(t, output, `tmux send-keys`)
		assert.Contains(t, output, `"nvim ." Enter`)

		// Verify attach
		assert.Contains(t, output, `tmux attach-session -t "proj_dev"`)
	})

	t.Run("generates script with multiple panes", func(t *testing.T) {
		session := createMultiPaneSession()
		output, err := gen.Generate(session)
		require.NoError(t, err)

		// Verify split command
		assert.Contains(t, output, "tmux split-window -h")

		// Verify both panes referenced
		assert.Contains(t, output, "editor")
		assert.Contains(t, output, "nvim .")
	})

	t.Run("generates script with multiple windows", func(t *testing.T) {
		session := createMultiWindowSession()
		output, err := gen.Generate(session)
		require.NoError(t, err)

		// Verify new window creation
		assert.Contains(t, output, `tmux new-window -t "proj_fullstack" -n "server"`)

		// Verify window comments
		assert.Contains(t, output, "# Window 1: code")
		assert.Contains(t, output, "# Window 2: server")
	})

	t.Run("handles empty working directory", func(t *testing.T) {
		session := createSimpleSession()
		session.WorkingDirectory = ""
		output, err := gen.Generate(session)
		require.NoError(t, err)

		assert.Contains(t, output, `WORKING_DIR="$HOME"`)
	})

	t.Run("handles environment variables", func(t *testing.T) {
		session := createSimpleSession()
		session.Environment = map[string]string{
			"NODE_ENV": "development",
			"DEBUG":    "true",
		}
		output, err := gen.Generate(session)
		require.NoError(t, err)

		assert.Contains(t, output, `export NODE_ENV="development"`)
		assert.Contains(t, output, `export DEBUG="true"`)
		assert.Contains(t, output, `tmux setenv`)
	})

	t.Run("handles pre and post commands", func(t *testing.T) {
		session := createSimpleSession()
		session.PreCommands = []string{"echo 'Starting...'"}
		session.PostCommands = []string{"echo 'Done!'"}
		output, err := gen.Generate(session)
		require.NoError(t, err)

		assert.Contains(t, output, "# Pre-session commands")
		assert.Contains(t, output, "echo 'Starting...'")
		assert.Contains(t, output, "# Post-session commands")
		assert.Contains(t, output, "echo 'Done!'")
	})

	t.Run("handles session with no windows", func(t *testing.T) {
		session := domain.Session{
			SessionName: "empty",
			Windows:     []domain.Window{},
		}
		output, err := gen.Generate(session)
		require.NoError(t, err)

		assert.Contains(t, output, `echo "No windows defined"`)
		assert.Contains(t, output, "exit 1")
	})
}

func TestTmuxGenerator_Name(t *testing.T) {
	gen := &TmuxGenerator{}
	assert.Equal(t, "tmux", gen.Name())
}

func TestTmuxGenerator_FileExtension(t *testing.T) {
	gen := &TmuxGenerator{}
	assert.Equal(t, ".sh", gen.FileExtension())
}

func TestIsAbsolutePath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/home/user", true},
		{"/", true},
		{"~", true},
		{"~/projects", true},
		{"$HOME", true},
		{"$HOME/projects", true},
		{"relative/path", false},
		{"./path", false},
		{"path", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isAbsolutePath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTmuxGenerator_PaneWorkingDirectory(t *testing.T) {
	gen := &TmuxGenerator{}

	t.Run("handles relative working directory", func(t *testing.T) {
		session := createSimpleSession()
		session.Windows[0].Panes[0].WorkingDirectory = "src"
		output, err := gen.Generate(session)
		require.NoError(t, err)

		// Should prefix with $WORKING_DIR
		assert.Contains(t, output, `cd $WORKING_DIR/src`)
	})

	t.Run("handles absolute working directory", func(t *testing.T) {
		session := createSimpleSession()
		session.Windows[0].Panes[0].WorkingDirectory = "/tmp/test"
		output, err := gen.Generate(session)
		require.NoError(t, err)

		// Should use path as-is
		assert.Contains(t, output, `cd /tmp/test`)
	})

	t.Run("handles home-relative working directory", func(t *testing.T) {
		session := createSimpleSession()
		session.Windows[0].Panes[0].WorkingDirectory = "~/other"
		output, err := gen.Generate(session)
		require.NoError(t, err)

		// Should use path as-is
		assert.Contains(t, output, `cd ~/other`)
	})
}

func TestTmuxGenerator_OutputIsValidBash(t *testing.T) {
	gen := &TmuxGenerator{}
	session := createMultiWindowSession()
	output, err := gen.Generate(session)
	require.NoError(t, err)

	// Basic validation that output looks like valid bash
	lines := strings.Split(output, "\n")

	// First line should be shebang
	assert.Equal(t, "#!/bin/bash", lines[0])

	// Should not have any unclosed quotes (basic check)
	quoteCount := strings.Count(output, `"`)
	assert.Equal(t, 0, quoteCount%2, "unbalanced quotes in output")
}
