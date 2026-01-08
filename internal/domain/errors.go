package domain

import "errors"

// Domain errors.
var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrSessionNotFound      = errors.New("session not found")
	ErrInvalidMultiplexer   = errors.New("invalid multiplexer type")
	ErrInvalidDirection     = errors.New("invalid direction")
	ErrProjectNameEmpty     = errors.New("project name cannot be empty")
	ErrSessionNameEmpty     = errors.New("session name cannot be empty")
	ErrProjectAlreadyExists = errors.New("project already exists")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrNoWindows            = errors.New("session must have at least one window")
)
