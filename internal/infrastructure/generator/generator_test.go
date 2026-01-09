package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

func TestNew(t *testing.T) {
	t.Run("creates TmuxGenerator for tmux", func(t *testing.T) {
		gen, err := New(domain.MultiplexerTmux)
		require.NoError(t, err)
		assert.Equal(t, "tmux", gen.Name())
		assert.Equal(t, ".sh", gen.FileExtension())
	})

	t.Run("creates ZellijGenerator for zellij", func(t *testing.T) {
		gen, err := New(domain.MultiplexerZellij)
		require.NoError(t, err)
		assert.Equal(t, "zellij", gen.Name())
		assert.Equal(t, ".sh", gen.FileExtension())
	})

	t.Run("returns error for unsupported multiplexer", func(t *testing.T) {
		_, err := New("unknown")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported multiplexer")
	})
}

// createSimpleSession creates a session with a single window and single pane.
func createSimpleSession() domain.Session {
	return domain.Session{
		ID:                "session-1",
		Name:              "dev",
		SessionName:       "proj_dev",
		Description:       "Development session",
		WorkingDirectory:  "~/projects",
		TargetMultiplexer: domain.MultiplexerTmux,
		DefaultWindowID:   "w1",
		Windows: []domain.Window{
			{
				ID:            "w1",
				Name:          "editor",
				RootDirection: domain.DirectionHorizontal,
				Panes: []domain.Pane{
					{
						ID:      "p1",
						Name:    "editor",
						Command: "nvim .",
						Size:    100,
					},
				},
			},
		},
	}
}

// createMultiPaneSession creates a session with multiple panes.
func createMultiPaneSession() domain.Session {
	return domain.Session{
		ID:                "session-2",
		Name:              "dev",
		SessionName:       "proj_dev",
		WorkingDirectory:  "~/projects",
		TargetMultiplexer: domain.MultiplexerTmux,
		DefaultWindowID:   "w1",
		Windows: []domain.Window{
			{
				ID:            "w1",
				Name:          "main",
				RootDirection: domain.DirectionHorizontal,
				Panes: []domain.Pane{
					{
						ID:        "root",
						Size:      100,
						Direction: domain.DirectionHorizontal,
						Children: []domain.Pane{
							{
								ID:      "p1",
								Name:    "editor",
								Command: "nvim .",
								Size:    60,
							},
							{
								ID:      "p2",
								Name:    "terminal",
								Command: "",
								Size:    40,
							},
						},
					},
				},
			},
		},
	}
}

// createMultiWindowSession creates a session with multiple windows.
func createMultiWindowSession() domain.Session {
	return domain.Session{
		ID:                "session-3",
		Name:              "fullstack",
		SessionName:       "proj_fullstack",
		WorkingDirectory:  "~/app",
		TargetMultiplexer: domain.MultiplexerTmux,
		DefaultWindowID:   "w1",
		Windows: []domain.Window{
			{
				ID:            "w1",
				Name:          "code",
				RootDirection: domain.DirectionHorizontal,
				Panes: []domain.Pane{
					{
						ID:      "p1",
						Name:    "editor",
						Command: "nvim .",
						Size:    100,
					},
				},
			},
			{
				ID:            "w2",
				Name:          "server",
				RootDirection: domain.DirectionVertical,
				Panes: []domain.Pane{
					{
						ID:        "root",
						Size:      100,
						Direction: domain.DirectionVertical,
						Children: []domain.Pane{
							{
								ID:               "p1",
								Name:             "frontend",
								Command:          "npm run dev",
								WorkingDirectory: "frontend",
								Size:             50,
							},
							{
								ID:               "p2",
								Name:             "backend",
								Command:          "go run .",
								WorkingDirectory: "backend",
								Size:             50,
							},
						},
					},
				},
			},
		},
	}
}
