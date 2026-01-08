package multiplexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajmasia/shellify/internal/domain"
)

func TestMultiplexerAvailability_HasAnyMultiplexer(t *testing.T) {
	tests := []struct {
		name     string
		avail    MultiplexerAvailability
		expected bool
	}{
		{
			name:     "no multiplexers",
			avail:    MultiplexerAvailability{},
			expected: false,
		},
		{
			name: "only tmux",
			avail: MultiplexerAvailability{
				TmuxAvailable: true,
			},
			expected: true,
		},
		{
			name: "only zellij",
			avail: MultiplexerAvailability{
				ZellijAvailable: true,
			},
			expected: true,
		},
		{
			name: "both available",
			avail: MultiplexerAvailability{
				TmuxAvailable:   true,
				ZellijAvailable: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.avail.HasAnyMultiplexer()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMultiplexerAvailability_HasMultiplexer(t *testing.T) {
	avail := MultiplexerAvailability{
		TmuxAvailable:   true,
		ZellijAvailable: false,
		TmuxPath:        "/usr/bin/tmux",
	}

	assert.True(t, avail.HasMultiplexer(domain.MultiplexerTmux))
	assert.False(t, avail.HasMultiplexer(domain.MultiplexerZellij))
	assert.False(t, avail.HasMultiplexer(domain.MultiplexerType("unknown")))
}

func TestMultiplexerAvailability_GetPath(t *testing.T) {
	avail := MultiplexerAvailability{
		TmuxAvailable:   true,
		ZellijAvailable: true,
		TmuxPath:        "/usr/bin/tmux",
		ZellijPath:      "/usr/bin/zellij",
	}

	assert.Equal(t, "/usr/bin/tmux", avail.GetPath(domain.MultiplexerTmux))
	assert.Equal(t, "/usr/bin/zellij", avail.GetPath(domain.MultiplexerZellij))
	assert.Equal(t, "", avail.GetPath(domain.MultiplexerType("unknown")))
}

func TestGetInstallInstructions(t *testing.T) {
	tmuxInstructions := GetInstallInstructions(domain.MultiplexerTmux)
	assert.Contains(t, tmuxInstructions, "tmux is not installed")
	assert.Contains(t, tmuxInstructions, "apt install")

	zellijInstructions := GetInstallInstructions(domain.MultiplexerZellij)
	assert.Contains(t, zellijInstructions, "zellij is not installed")
	assert.Contains(t, zellijInstructions, "cargo install")

	unknownInstructions := GetInstallInstructions(domain.MultiplexerType("unknown"))
	assert.Equal(t, "Unknown multiplexer", unknownInstructions)
}
