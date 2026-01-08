package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplexerType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		m        MultiplexerType
		expected bool
	}{
		{"tmux is valid", MultiplexerTmux, true},
		{"zellij is valid", MultiplexerZellij, true},
		{"empty is invalid", MultiplexerType(""), false},
		{"unknown is invalid", MultiplexerType("screen"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.m.IsValid())
		})
	}
}

func TestMultiplexerType_String(t *testing.T) {
	assert.Equal(t, "tmux", MultiplexerTmux.String())
	assert.Equal(t, "zellij", MultiplexerZellij.String())
}

func TestDirection_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		d        Direction
		expected bool
	}{
		{"horizontal is valid", DirectionHorizontal, true},
		{"vertical is valid", DirectionVertical, true},
		{"empty is invalid", Direction(""), false},
		{"unknown is invalid", Direction("diagonal"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.d.IsValid())
		})
	}
}

func TestDirection_String(t *testing.T) {
	assert.Equal(t, "horizontal", DirectionHorizontal.String())
	assert.Equal(t, "vertical", DirectionVertical.String())
}
