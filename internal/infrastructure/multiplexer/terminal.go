package multiplexer

import (
	"fmt"
	"os/exec"
)

// supportedTerminals is the ordered list of terminal emulators to check.
var supportedTerminals = []string{
	"alacritty",
	"kitty",
	"wezterm",
	"gnome-terminal",
	"konsole",
	"xfce4-terminal",
	"xterm",
}

// DetectTerminal finds an available terminal emulator.
// Returns empty string if no terminal is found.
func DetectTerminal() string {
	for _, term := range supportedTerminals {
		if _, err := exec.LookPath(term); err == nil {
			return term
		}
	}
	return ""
}

// OpenTerminalWithCommand opens a new terminal window and executes the given command.
func OpenTerminalWithCommand(terminal, command string, args ...string) error {
	if terminal == "" {
		return fmt.Errorf("no terminal emulator specified")
	}

	fullArgs := append([]string{command}, args...)
	fullCmd := buildCommandString(command, args)

	var cmd *exec.Cmd

	switch terminal {
	case "alacritty":
		cmdArgs := append([]string{"-e"}, fullArgs...)
		cmd = exec.Command("alacritty", cmdArgs...)
	case "kitty":
		cmd = exec.Command("kitty", fullArgs...)
	case "wezterm":
		cmdArgs := append([]string{"start", "--"}, fullArgs...)
		cmd = exec.Command("wezterm", cmdArgs...)
	case "gnome-terminal":
		cmdArgs := append([]string{"--"}, fullArgs...)
		cmd = exec.Command("gnome-terminal", cmdArgs...)
	case "konsole":
		cmdArgs := append([]string{"-e"}, fullArgs...)
		cmd = exec.Command("konsole", cmdArgs...)
	case "xfce4-terminal":
		cmd = exec.Command("xfce4-terminal", "-e", fullCmd)
	case "xterm":
		cmd = exec.Command("xterm", "-e", fullCmd)
	default:
		return fmt.Errorf("unsupported terminal: %s", terminal)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting terminal: %w", err)
	}

	return nil
}

// buildCommandString builds a space-separated command string.
func buildCommandString(command string, args []string) string {
	result := command
	for _, arg := range args {
		result += " " + arg
	}
	return result
}
