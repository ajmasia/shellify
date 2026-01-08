package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPane_IsContainer(t *testing.T) {
	t.Run("pane with children is container", func(t *testing.T) {
		pane := Pane{
			ID:   "1",
			Name: "container",
			Children: []Pane{
				{ID: "2", Name: "child1"},
				{ID: "3", Name: "child2"},
			},
		}
		assert.True(t, pane.IsContainer())
	})

	t.Run("pane without children is not container", func(t *testing.T) {
		pane := Pane{
			ID:   "1",
			Name: "leaf",
		}
		assert.False(t, pane.IsContainer())
	})
}

func TestPane_IsLeaf(t *testing.T) {
	t.Run("pane without children is leaf", func(t *testing.T) {
		pane := Pane{
			ID:   "1",
			Name: "leaf",
		}
		assert.True(t, pane.IsLeaf())
	})

	t.Run("pane with children is not leaf", func(t *testing.T) {
		pane := Pane{
			ID:   "1",
			Name: "container",
			Children: []Pane{
				{ID: "2", Name: "child"},
			},
		}
		assert.False(t, pane.IsLeaf())
	})
}

func TestPane_CountLeaves(t *testing.T) {
	t.Run("single leaf pane returns 1", func(t *testing.T) {
		pane := Pane{ID: "1", Name: "leaf"}
		assert.Equal(t, 1, pane.CountLeaves())
	})

	t.Run("container with two leaves returns 2", func(t *testing.T) {
		pane := Pane{
			ID:   "1",
			Name: "container",
			Children: []Pane{
				{ID: "2", Name: "leaf1"},
				{ID: "3", Name: "leaf2"},
			},
		}
		assert.Equal(t, 2, pane.CountLeaves())
	})

	t.Run("nested containers count all leaves", func(t *testing.T) {
		// Structure:
		// root
		// ├── child1 (leaf)
		// └── child2 (container)
		//     ├── grandchild1 (leaf)
		//     └── grandchild2 (leaf)
		pane := Pane{
			ID:   "root",
			Name: "root",
			Children: []Pane{
				{ID: "child1", Name: "leaf"},
				{
					ID:   "child2",
					Name: "container",
					Children: []Pane{
						{ID: "grandchild1", Name: "leaf"},
						{ID: "grandchild2", Name: "leaf"},
					},
				},
			},
		}
		assert.Equal(t, 3, pane.CountLeaves())
	})
}
