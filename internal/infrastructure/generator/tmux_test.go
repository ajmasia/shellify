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

		// Verify working directory (~ expanded to $HOME)
		assert.Contains(t, output, `WORKING_DIR="$HOME/projects"`)

		// Verify session creation with inline command
		assert.Contains(t, output, `tmux new-session -d -s "proj_dev" -n "editor"`)
		assert.Contains(t, output, `"$SHELL -i -c 'nvim .; exec $SHELL'"`)

		// Verify no send-keys for commands
		assert.NotContains(t, output, `tmux send-keys`)

		// Verify attach
		assert.Contains(t, output, `tmux attach-session -t "proj_dev"`)
	})

	t.Run("generates script with multiple panes", func(t *testing.T) {
		session := createMultiPaneSession()
		output, err := gen.Generate(session)
		require.NoError(t, err)

		// Verify split command with inline command for second pane (no command → no inline)
		assert.Contains(t, output, "tmux split-window -h")

		// First pane gets inline command via new-session
		assert.Contains(t, output, `"$SHELL -i -c 'nvim .; exec $SHELL'"`)

		// No send-keys
		assert.NotContains(t, output, `tmux send-keys`)
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
		assert.NotContains(t, output, `tmux setenv`)
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

func TestBuildInlineCommand(t *testing.T) {
	t.Run("returns empty for empty command", func(t *testing.T) {
		assert.Equal(t, "", buildInlineCommand(""))
	})

	t.Run("wraps simple command", func(t *testing.T) {
		result := buildInlineCommand("nvim .")
		assert.Equal(t, `"$SHELL -i -c 'nvim .; exec $SHELL'"`, result)
	})

	t.Run("escapes single quotes with apostrophe trick", func(t *testing.T) {
		result := buildInlineCommand("fettly-show-banner 'Genially Mono'")
		assert.Equal(t, `"$SHELL -i -c 'fettly-show-banner '\''Genially Mono'\''; exec $SHELL'"`, result)
	})

	t.Run("escapes double quotes", func(t *testing.T) {
		result := buildInlineCommand(`fettly-show-banner "Genially Mono"`)
		assert.Equal(t, `"$SHELL -i -c 'fettly-show-banner \"Genially Mono\"; exec $SHELL'"`, result)
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

		assert.Contains(t, output, `-c "$WORKING_DIR/src"`)
		assert.NotContains(t, output, `tmux send-keys`)
	})

	t.Run("handles absolute working directory", func(t *testing.T) {
		session := createSimpleSession()
		session.Windows[0].Panes[0].WorkingDirectory = "/tmp/test"
		output, err := gen.Generate(session)
		require.NoError(t, err)

		assert.Contains(t, output, `-c "/tmp/test"`)
		assert.NotContains(t, output, `tmux send-keys`)
	})

	t.Run("handles home-relative working directory", func(t *testing.T) {
		session := createSimpleSession()
		session.Windows[0].Panes[0].WorkingDirectory = "~/other"
		output, err := gen.Generate(session)
		require.NoError(t, err)

		assert.Contains(t, output, `-c "$HOME/other"`)
		assert.NotContains(t, output, `tmux send-keys`)
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
