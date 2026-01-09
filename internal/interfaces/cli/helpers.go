package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"
)

// printJSON outputs data as formatted JSON.
func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// printTable outputs data in a formatted table.
func printTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print headers
	for i, h := range headers {
		if i > 0 {
			_, _ = fmt.Fprint(w, "\t")
		}
		_, _ = fmt.Fprint(w, h)
	}
	_, _ = fmt.Fprintln(w)

	// Print rows
	for _, row := range rows {
		for i, col := range row {
			if i > 0 {
				_, _ = fmt.Fprint(w, "\t")
			}
			_, _ = fmt.Fprint(w, col)
		}
		_, _ = fmt.Fprintln(w)
	}

	_ = w.Flush()
}

// formatTime returns a human-readable relative time string.
func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("2006-01-02")
	}
}

// truncateID returns the first 8 characters of an ID.
func truncateID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8]
}

// generateSessionName creates a session name from project prefix and session name.
func generateSessionName(prefix, name string) string {
	sessionName := strings.ToLower(strings.TrimSpace(name))
	sessionName = strings.ReplaceAll(sessionName, " ", "-")

	if prefix != "" {
		return prefix + "_" + sessionName
	}
	return sessionName
}

// ensureServerRunning checks if the server is running and starts it if not.
func ensureServerRunning() error {
	// Check if server is already running
	if isServerRunning() {
		return nil
	}

	// Start server in daemon mode
	fmt.Println("Starting server...")
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Build command args
	args := []string{"server", "-d"}

	// Try to find gui/dist directory for development
	guiDistPaths := []string{
		"gui/dist",
		"../gui/dist",
		"../../gui/dist",
	}
	for _, path := range guiDistPaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			args = append(args, "-s", path)
			break
		}
	}

	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Wait a moment for server to be ready
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Server running at http://localhost:3777")
	fmt.Println("Use 'sfy server stop' when done")
	fmt.Println()
	return nil
}

// isServerRunning checks if the server process is running.
func isServerRunning() bool {
	pid, err := readPIDFile()
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix, FindProcess always succeeds, send signal 0 to check
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false
	}

	return true
}
