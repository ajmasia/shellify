// Package domain contains the core business entities and rules.
package domain

// MultiplexerType represents the target terminal multiplexer.
type MultiplexerType string

const (
	MultiplexerTmux   MultiplexerType = "tmux"
	MultiplexerZellij MultiplexerType = "zellij"
)

// IsValid returns true if the multiplexer type is valid.
func (m MultiplexerType) IsValid() bool {
	return m == MultiplexerTmux || m == MultiplexerZellij
}

// String returns the string representation of the multiplexer type.
func (m MultiplexerType) String() string {
	return string(m)
}

// Direction represents the split direction for panes.
type Direction string

const (
	DirectionHorizontal Direction = "horizontal"
	DirectionVertical   Direction = "vertical"
)

// IsValid returns true if the direction is valid.
func (d Direction) IsValid() bool {
	return d == DirectionHorizontal || d == DirectionVertical
}

// String returns the string representation of the direction.
func (d Direction) String() string {
	return string(d)
}
