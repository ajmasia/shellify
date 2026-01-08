package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSession_ToMetadata(t *testing.T) {
	now := time.Now()
	session := Session{
		ID:                "session-1",
		Name:              "dev-session",
		Description:       "Development session",
		TargetMultiplexer: MultiplexerTmux,
		CreatedAt:         now,
		UpdatedAt:         now,
		Windows: []Window{
			{ID: "win-1", Name: "editor"},
		},
	}

	metadata := session.ToMetadata()

	assert.Equal(t, session.ID, metadata.ID)
	assert.Equal(t, session.Name, metadata.Name)
	assert.Equal(t, session.Description, metadata.Description)
	assert.Equal(t, session.TargetMultiplexer, metadata.TargetMultiplexer)
	assert.Equal(t, session.CreatedAt, metadata.CreatedAt)
	assert.Equal(t, session.UpdatedAt, metadata.UpdatedAt)
}

func TestSession_CountWindows(t *testing.T) {
	t.Run("empty session has 0 windows", func(t *testing.T) {
		session := Session{}
		assert.Equal(t, 0, session.CountWindows())
	})

	t.Run("session with 3 windows", func(t *testing.T) {
		session := Session{
			Windows: []Window{
				{ID: "1", Name: "win1"},
				{ID: "2", Name: "win2"},
				{ID: "3", Name: "win3"},
			},
		}
		assert.Equal(t, 3, session.CountWindows())
	})
}

func TestSession_CountPanes(t *testing.T) {
	t.Run("empty session has 0 panes", func(t *testing.T) {
		session := Session{}
		assert.Equal(t, 0, session.CountPanes())
	})

	t.Run("session with multiple windows and panes", func(t *testing.T) {
		// Window 1: 2 panes
		// Window 2: 1 pane with 2 children (2 leaves)
		// Total: 4 panes
		session := Session{
			Windows: []Window{
				{
					ID:   "win1",
					Name: "window1",
					Panes: []Pane{
						{ID: "p1", Name: "pane1"},
						{ID: "p2", Name: "pane2"},
					},
				},
				{
					ID:   "win2",
					Name: "window2",
					Panes: []Pane{
						{
							ID:   "p3",
							Name: "container",
							Children: []Pane{
								{ID: "p4", Name: "child1"},
								{ID: "p5", Name: "child2"},
							},
						},
					},
				},
			},
		}
		assert.Equal(t, 4, session.CountPanes())
	})
}

func TestWindow_CountPanes(t *testing.T) {
	t.Run("empty window has 0 panes", func(t *testing.T) {
		window := Window{}
		assert.Equal(t, 0, window.CountPanes())
	})

	t.Run("window with simple panes", func(t *testing.T) {
		window := Window{
			Panes: []Pane{
				{ID: "1", Name: "pane1"},
				{ID: "2", Name: "pane2"},
			},
		}
		assert.Equal(t, 2, window.CountPanes())
	})

	t.Run("window with nested panes", func(t *testing.T) {
		window := Window{
			Panes: []Pane{
				{ID: "1", Name: "leaf"},
				{
					ID:   "2",
					Name: "container",
					Children: []Pane{
						{ID: "3", Name: "child1"},
						{ID: "4", Name: "child2"},
					},
				},
			},
		}
		assert.Equal(t, 3, window.CountPanes())
	})
}
