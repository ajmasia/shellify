package domain

import "time"

// Project represents a collection of related sessions.
type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	SessionPrefix string    `json:"sessionPrefix"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ProjectMetadata is an alias for Project since all fields are metadata.
type ProjectMetadata = Project
