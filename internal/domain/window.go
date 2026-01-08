package domain

// Window represents a window within a multiplexer session.
type Window struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	WorkingDirectory string            `json:"workingDirectory,omitempty"`
	Environment      map[string]string `json:"environment,omitempty"`
	Panes            []Pane            `json:"panes"`
	RootDirection    Direction         `json:"rootDirection"`
}

// CountPanes returns the total number of leaf panes in the window.
func (w *Window) CountPanes() int {
	count := 0
	for i := range w.Panes {
		count += w.Panes[i].CountLeaves()
	}
	return count
}
