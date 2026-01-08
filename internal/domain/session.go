package domain

import "time"

// Session represents a multiplexer session configuration.
type Session struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	SessionName       string            `json:"sessionName"`
	Description       string            `json:"description,omitempty"`
	WorkingDirectory  string            `json:"workingDirectory,omitempty"`
	Environment       map[string]string `json:"environment,omitempty"`
	PreCommands       []string          `json:"preCommands,omitempty"`
	PostCommands      []string          `json:"postCommands,omitempty"`
	Windows           []Window          `json:"windows"`
	DefaultWindowID   string            `json:"defaultWindowId,omitempty"`
	TargetMultiplexer MultiplexerType   `json:"targetMultiplexer"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

// SessionMetadata contains lightweight session information for listing.
type SessionMetadata struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description,omitempty"`
	TargetMultiplexer MultiplexerType `json:"targetMultiplexer,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}

// ToMetadata extracts metadata from a full session.
func (s *Session) ToMetadata() SessionMetadata {
	return SessionMetadata{
		ID:                s.ID,
		Name:              s.Name,
		Description:       s.Description,
		TargetMultiplexer: s.TargetMultiplexer,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

// CountWindows returns the number of windows in the session.
func (s *Session) CountWindows() int {
	return len(s.Windows)
}

// CountPanes returns the total number of panes across all windows.
func (s *Session) CountPanes() int {
	count := 0
	for i := range s.Windows {
		count += s.Windows[i].CountPanes()
	}
	return count
}
