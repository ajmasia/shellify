package generator

import (
	"fmt"

	"github.com/ajmasia/shellify/internal/domain"
)

// Generator defines the interface for generating multiplexer scripts/layouts.
type Generator interface {
	// Generate creates a script or layout from a session configuration.
	Generate(session domain.Session) (string, error)
	// Name returns the generator name (e.g., "tmux", "zellij").
	Name() string
	// FileExtension returns the file extension for generated files (e.g., ".sh", ".kdl").
	FileExtension() string
}

// New creates a new Generator for the specified multiplexer type.
func New(multiplexer domain.MultiplexerType) (Generator, error) {
	switch multiplexer {
	case domain.MultiplexerTmux:
		return &TmuxGenerator{}, nil
	case domain.MultiplexerZellij:
		return &ZellijGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported multiplexer: %s", multiplexer)
	}
}
