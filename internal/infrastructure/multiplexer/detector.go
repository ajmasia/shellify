// Package multiplexer provides infrastructure for multiplexer detection and session management.
package multiplexer

import (
	"os/exec"

	"github.com/ajmasia/shellify/internal/domain"
)

// MultiplexerAvailability holds info about available multiplexers.
type MultiplexerAvailability struct {
	TmuxAvailable   bool
	ZellijAvailable bool
	TmuxPath        string
	ZellijPath      string
}

// CheckMultiplexers detects which multiplexers are available on the system.
func CheckMultiplexers() MultiplexerAvailability {
	avail := MultiplexerAvailability{}

	if path, err := exec.LookPath("tmux"); err == nil {
		avail.TmuxAvailable = true
		avail.TmuxPath = path
	}

	if path, err := exec.LookPath("zellij"); err == nil {
		avail.ZellijAvailable = true
		avail.ZellijPath = path
	}

	return avail
}

// HasAnyMultiplexer returns true if at least one multiplexer is available.
func (a MultiplexerAvailability) HasAnyMultiplexer() bool {
	return a.TmuxAvailable || a.ZellijAvailable
}

// HasMultiplexer checks if a specific multiplexer is available.
func (a MultiplexerAvailability) HasMultiplexer(m domain.MultiplexerType) bool {
	switch m {
	case domain.MultiplexerTmux:
		return a.TmuxAvailable
	case domain.MultiplexerZellij:
		return a.ZellijAvailable
	default:
		return false
	}
}

// GetPath returns the path for a specific multiplexer.
func (a MultiplexerAvailability) GetPath(m domain.MultiplexerType) string {
	switch m {
	case domain.MultiplexerTmux:
		return a.TmuxPath
	case domain.MultiplexerZellij:
		return a.ZellijPath
	default:
		return ""
	}
}

// GetInstallInstructions returns installation instructions for a multiplexer.
func GetInstallInstructions(m domain.MultiplexerType) string {
	switch m {
	case domain.MultiplexerTmux:
		return `tmux is not installed. Install it using:
  - Debian/Ubuntu: sudo apt install tmux
  - Fedora/RHEL:   sudo dnf install tmux
  - macOS:         brew install tmux
  - Arch:          sudo pacman -S tmux`
	case domain.MultiplexerZellij:
		return `zellij is not installed. Install it using:
  - Cargo:         cargo install zellij
  - macOS:         brew install zellij
  - Arch:          sudo pacman -S zellij
  - Binary:        https://github.com/zellij-org/zellij/releases`
	default:
		return "Unknown multiplexer"
	}
}
