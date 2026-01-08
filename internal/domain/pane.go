package domain

// Pane represents a terminal pane within a window.
// Panes can be nested to create complex layouts via the Children field.
type Pane struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Command          string            `json:"command,omitempty"`
	WorkingDirectory string            `json:"workingDirectory,omitempty"`
	Environment      map[string]string `json:"environment,omitempty"`
	Size             float64           `json:"size"`
	Children         []Pane            `json:"children,omitempty"`
	Direction        Direction         `json:"direction,omitempty"`
}

// IsContainer returns true if the pane contains child panes.
func (p *Pane) IsContainer() bool {
	return len(p.Children) > 0
}

// IsLeaf returns true if the pane is a terminal pane (no children).
func (p *Pane) IsLeaf() bool {
	return len(p.Children) == 0
}

// CountLeaves returns the total number of leaf panes in the tree.
func (p *Pane) CountLeaves() int {
	if p.IsLeaf() {
		return 1
	}

	count := 0
	for i := range p.Children {
		count += p.Children[i].CountLeaves()
	}
	return count
}
