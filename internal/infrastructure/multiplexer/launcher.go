package multiplexer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/ajmasia/shellify/internal/domain"
)

// Launcher handles session launching for different multiplexers.
type Launcher struct {
	tempDir  string
	terminal string
}

// NewLauncher creates a new Launcher instance.
func NewLauncher() (*Launcher, error) {
	tempDir := filepath.Join(os.TempDir(), "shellify-sessions")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("creating temp directory: %w", err)
	}
	return &Launcher{
		tempDir:  tempDir,
		terminal: DetectTerminal(),
	}, nil
}

// LaunchTmux starts a new tmux session using the provided script content.
// If the session already exists, it attaches to it instead.
// This replaces the current process.
func (l *Launcher) LaunchTmux(sessionName, scriptContent string) error {
	if l.IsTmuxRunning(sessionName) {
		return l.AttachTmux(sessionName)
	}

	scriptPath := filepath.Join(l.tempDir, sessionName+".sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("writing script: %w", err)
	}

	return syscall.Exec("/bin/bash", []string{"bash", scriptPath}, os.Environ())
}

// LaunchZellij starts a new zellij session using the provided script content.
// The script is a bash wrapper that creates a temporary KDL layout and launches zellij.
// If the session already exists, it attaches to it instead.
// This replaces the current process.
func (l *Launcher) LaunchZellij(sessionName, scriptContent string) error {
	if l.IsZellijRunning(sessionName) {
		return l.AttachZellij(sessionName)
	}

	scriptPath := filepath.Join(l.tempDir, sessionName+".sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("writing script: %w", err)
	}

	return syscall.Exec("/bin/bash", []string{"bash", scriptPath}, os.Environ())
}

// LaunchTmuxDetached starts a tmux session in a new terminal window.
// Used by server to avoid replacing the server process.
func (l *Launcher) LaunchTmuxDetached(sessionName, scriptContent string) error {
	if l.terminal == "" {
		return fmt.Errorf("no terminal emulator found")
	}

	if l.IsTmuxRunning(sessionName) {
		return OpenTerminalWithCommand(l.terminal, "tmux", "attach-session", "-t", sessionName)
	}

	scriptPath := filepath.Join(l.tempDir, sessionName+".sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("writing script: %w", err)
	}

	return OpenTerminalWithCommand(l.terminal, "bash", scriptPath)
}

// LaunchZellijDetached starts a zellij session in a new terminal window.
// The script is a bash wrapper that creates a temporary KDL layout and launches zellij.
// Used by server to avoid replacing the server process.
func (l *Launcher) LaunchZellijDetached(sessionName, scriptContent string) error {
	if l.terminal == "" {
		return fmt.Errorf("no terminal emulator found")
	}

	if l.IsZellijRunning(sessionName) {
		return OpenTerminalWithCommand(l.terminal, "zellij", "attach", sessionName)
	}

	scriptPath := filepath.Join(l.tempDir, sessionName+".sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("writing script: %w", err)
	}

	return OpenTerminalWithCommand(l.terminal, "bash", scriptPath)
}

// AttachTmux attaches to an existing tmux session.
// This replaces the current process.
func (l *Launcher) AttachTmux(sessionName string) error {
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("tmux not found: %w", err)
	}
	return syscall.Exec(tmuxPath, []string{"tmux", "attach-session", "-t", sessionName}, os.Environ())
}

// AttachZellij attaches to an existing zellij session.
// This replaces the current process.
func (l *Launcher) AttachZellij(sessionName string) error {
	zellijPath, err := exec.LookPath("zellij")
	if err != nil {
		return fmt.Errorf("zellij not found: %w", err)
	}
	return syscall.Exec(zellijPath, []string{"zellij", "attach", sessionName}, os.Environ())
}

// StopTmux kills a tmux session.
func (l *Launcher) StopTmux(sessionName string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StopZellij kills a zellij session.
func (l *Launcher) StopZellij(sessionName string) error {
	cmd := exec.Command("zellij", "kill-session", sessionName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsTmuxRunning checks if a tmux session is currently active.
func (l *Launcher) IsTmuxRunning(sessionName string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", sessionName)
	return cmd.Run() == nil
}

// IsZellijRunning checks if a zellij session is currently active.
func (l *Launcher) IsZellijRunning(sessionName string) bool {
	cmd := exec.Command("zellij", "list-sessions")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	cleanOutput := stripAnsiCodes(string(output))
	return containsSession(cleanOutput, sessionName)
}

// IsTmuxAttached checks if a tmux session has any clients attached.
func (l *Launcher) IsTmuxAttached(sessionName string) bool {
	cmd := exec.Command("tmux", "list-clients", "-t", sessionName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

// IsZellijAttached always returns false as zellij doesn't provide this info easily.
func (l *Launcher) IsZellijAttached(_ string) bool {
	return false
}

// IsRunning checks if a session is currently running for the given multiplexer.
func (l *Launcher) IsRunning(sessionName string, multiplexer domain.MultiplexerType) bool {
	switch multiplexer {
	case domain.MultiplexerTmux:
		return l.IsTmuxRunning(sessionName)
	case domain.MultiplexerZellij:
		return l.IsZellijRunning(sessionName)
	default:
		return false
	}
}

// IsAttached checks if a session has clients attached.
func (l *Launcher) IsAttached(sessionName string, multiplexer domain.MultiplexerType) bool {
	switch multiplexer {
	case domain.MultiplexerTmux:
		return l.IsTmuxAttached(sessionName)
	case domain.MultiplexerZellij:
		return l.IsZellijAttached(sessionName)
	default:
		return false
	}
}

// containsSession checks if a session name appears in zellij list-sessions output.
func containsSession(output, sessionName string) bool {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Skip exited sessions
		if strings.Contains(line, "EXITED") {
			continue
		}
		// Check exact match or prefix match (session might have status info)
		if line == sessionName || strings.HasPrefix(line, sessionName+" ") {
			return true
		}
	}
	return false
}

// stripAnsiCodes removes ANSI escape sequences from a string.
func stripAnsiCodes(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && !((s[j] >= 'A' && s[j] <= 'Z') || (s[j] >= 'a' && s[j] <= 'z')) {
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}
