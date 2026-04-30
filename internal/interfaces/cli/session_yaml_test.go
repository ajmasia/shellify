package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajmasia/shellify/internal/domain"
)

// sessionForTest builds a minimal domain.Session for round-trip tests.
func sessionForTest() domain.Session {
	return domain.Session{
		ID:                "test-id",
		Name:              "my-session",
		Description:       "test session",
		TargetMultiplexer: domain.MultiplexerTmux,
		WorkingDirectory:  "~/projects",
	}
}

func TestYAMLRoundTrip_SingleWindowSinglePane(t *testing.T) {
	s := sessionForTest()
	s.Windows = []domain.Window{
		{
			ID:            "w1",
			Name:          "editor",
			RootDirection: domain.DirectionHorizontal,
			Panes: []domain.Pane{
				{ID: "w1-p1", Name: "nvim", Command: "nvim .", Size: 100},
			},
		},
	}
	s.DefaultWindowID = "w1"

	data, err := SessionToYAML(s)
	require.NoError(t, err)

	got, err := YAMLToSession(data)
	require.NoError(t, err)

	assert.Equal(t, s.Name, got.Name)
	assert.Equal(t, s.Description, got.Description)
	assert.Equal(t, s.TargetMultiplexer, got.TargetMultiplexer)
	assert.Equal(t, s.WorkingDirectory, got.WorkingDirectory)
	require.Len(t, got.Windows, 1)
	assert.Equal(t, "editor", got.Windows[0].Name)
	require.Len(t, got.Windows[0].Panes, 1)
	leaf := got.Windows[0].Panes[0]
	assert.True(t, leaf.IsLeaf())
	assert.Equal(t, "nvim .", leaf.Command)
}

func TestYAMLRoundTrip_TwoWindowsHorizontalSplit(t *testing.T) {
	s := sessionForTest()
	s.Windows = []domain.Window{
		{
			ID:            "w1",
			Name:          "editor",
			RootDirection: domain.DirectionHorizontal,
			Panes: []domain.Pane{
				{
					ID:        "w1-root",
					Size:      100,
					Direction: domain.DirectionHorizontal,
					Children: []domain.Pane{
						{ID: "w1-p1", Name: "nvim", Command: "nvim .", Size: 65},
						{ID: "w1-p2", Name: "term", Size: 35},
					},
				},
			},
		},
		{
			ID:            "w2",
			Name:          "git",
			RootDirection: domain.DirectionHorizontal,
			Panes: []domain.Pane{
				{ID: "w2-p1", Command: "lazygit", Size: 100},
			},
		},
	}
	s.DefaultWindowID = "w1"

	data, err := SessionToYAML(s)
	require.NoError(t, err)

	got, err := YAMLToSession(data)
	require.NoError(t, err)

	require.Len(t, got.Windows, 2)
	assert.Equal(t, "editor", got.Windows[0].Name)
	assert.Equal(t, "git", got.Windows[1].Name)

	// editor window: root container with 2 children
	editorRoot := got.Windows[0].Panes[0]
	assert.True(t, editorRoot.IsContainer())
	require.Len(t, editorRoot.Children, 2)
	assert.Equal(t, "nvim .", editorRoot.Children[0].Command)
	assert.InDelta(t, 65.0, editorRoot.Children[0].Size, 0.1)
	assert.InDelta(t, 35.0, editorRoot.Children[1].Size, 0.1)

	// git window: single leaf
	gitPane := got.Windows[1].Panes[0]
	assert.True(t, gitPane.IsLeaf())
	assert.Equal(t, "lazygit", gitPane.Command)
}

func TestYAMLRoundTrip_NestedPaneTree(t *testing.T) {
	// editor (65%) | [test (50%) / dev (50%)] (35%) vertical
	s := sessionForTest()
	s.Windows = []domain.Window{
		{
			ID:            "w1",
			Name:          "workspace",
			RootDirection: domain.DirectionHorizontal,
			Panes: []domain.Pane{
				{
					ID:        "w1-root",
					Size:      100,
					Direction: domain.DirectionHorizontal,
					Children: []domain.Pane{
						{ID: "w1-p1", Command: "nvim .", Size: 65},
						{
							ID:        "w1-p2",
							Size:      35,
							Direction: domain.DirectionVertical,
							Children: []domain.Pane{
								{ID: "w1-p2-p1", Command: "npm test", Size: 50},
								{ID: "w1-p2-p2", Command: "npm run dev", Size: 50},
							},
						},
					},
				},
			},
		},
	}
	s.DefaultWindowID = "w1"

	data, err := SessionToYAML(s)
	require.NoError(t, err)

	got, err := YAMLToSession(data)
	require.NoError(t, err)

	require.Len(t, got.Windows, 1)
	root := got.Windows[0].Panes[0]
	require.True(t, root.IsContainer())
	require.Len(t, root.Children, 2)

	assert.Equal(t, "nvim .", root.Children[0].Command)

	nested := root.Children[1]
	assert.True(t, nested.IsContainer())
	assert.Equal(t, domain.DirectionVertical, nested.Direction)
	require.Len(t, nested.Children, 2)
	assert.Equal(t, "npm test", nested.Children[0].Command)
	assert.Equal(t, "npm run dev", nested.Children[1].Command)
}

func TestYAMLRoundTrip_EnvAndCommands(t *testing.T) {
	s := sessionForTest()
	s.Environment = map[string]string{"NODE_ENV": "development", "DEBUG": "true"}
	s.PreCommands = []string{"source ~/.nvm/nvm.sh"}
	s.PostCommands = []string{"echo done"}
	s.Windows = []domain.Window{
		{
			ID:   "w1",
			Name: "shell",
			Panes: []domain.Pane{
				{ID: "w1-p1", Size: 100},
			},
		},
	}
	s.DefaultWindowID = "w1"

	data, err := SessionToYAML(s)
	require.NoError(t, err)

	got, err := YAMLToSession(data)
	require.NoError(t, err)

	assert.Equal(t, s.Environment, got.Environment)
	assert.Equal(t, s.PreCommands, got.PreCommands)
	assert.Equal(t, s.PostCommands, got.PostCommands)
}

func TestYAMLToSession_DefaultsMultiplexerToTmux(t *testing.T) {
	input := `
name: test
windows:
  - name: main
    panes:
      - command: bash
`
	got, err := YAMLToSession([]byte(input))
	require.NoError(t, err)
	assert.Equal(t, domain.MultiplexerTmux, got.TargetMultiplexer)
}

func TestYAMLToSession_RejectsUnsupportedMultiplexer(t *testing.T) {
	input := `
name: test
multiplexer: zellij
windows:
  - name: main
    panes:
      - command: bash
`
	_, err := YAMLToSession([]byte(input))
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "zellij"))
}

func TestYAMLToSession_ErrorOnEmptyName(t *testing.T) {
	input := `
windows:
  - name: main
    panes:
      - command: bash
`
	_, err := YAMLToSession([]byte(input))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestYAMLToSession_ErrorOnNoWindows(t *testing.T) {
	input := `name: test`
	_, err := YAMLToSession([]byte(input))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "window")
}

func TestYAMLToSession_ErrorOnUnknownField(t *testing.T) {
	input := `
name: test
unknownField: bad
windows:
  - name: main
    panes:
      - command: bash
`
	_, err := YAMLToSession([]byte(input))
	require.Error(t, err)
}

func TestYAMLToSession_SizeDefaultsToEvenSplit(t *testing.T) {
	input := `
name: test
windows:
  - name: main
    panes:
      - command: nvim
      - command: bash
      - command: htop
`
	got, err := YAMLToSession([]byte(input))
	require.NoError(t, err)

	root := got.Windows[0].Panes[0]
	require.True(t, root.IsContainer())
	for _, child := range root.Children {
		assert.InDelta(t, 33.33, child.Size, 0.1)
	}
}

func TestYAMLToSession_DefaultWindowByName(t *testing.T) {
	input := `
name: test
defaultWindow: second
windows:
  - name: first
    panes:
      - command: bash
  - name: second
    panes:
      - command: vim
`
	got, err := YAMLToSession([]byte(input))
	require.NoError(t, err)
	assert.Equal(t, got.Windows[1].ID, got.DefaultWindowID)
}

func TestYAMLToSession_ErrorOnInvalidDefaultWindow(t *testing.T) {
	input := `
name: test
defaultWindow: nonexistent
windows:
  - name: main
    panes:
      - command: bash
`
	_, err := YAMLToSession([]byte(input))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
}

func TestYAMLStableIDs(t *testing.T) {
	// Verify that the same YAML always produces the same IDs.
	input := `
name: test
windows:
  - name: editor
    panes:
      - command: nvim .
        size: 70
      - command: bash
        size: 30
`
	got1, err := YAMLToSession([]byte(input))
	require.NoError(t, err)
	got2, err := YAMLToSession([]byte(input))
	require.NoError(t, err)

	assert.Equal(t, got1.Windows[0].ID, got2.Windows[0].ID)
	root1 := got1.Windows[0].Panes[0]
	root2 := got2.Windows[0].Panes[0]
	assert.Equal(t, root1.ID, root2.ID)
	assert.Equal(t, root1.Children[0].ID, root2.Children[0].ID)
}
