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

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
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
// Returns true if the server was just started, false if it was already running.
func ensureServerRunning() (bool, error) {
	// Check if server is already running
	if isServerRunning() {
		return false, nil
	}

	// Start server in daemon mode
	exe, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %w", err)
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
	// Suppress daemon output - it prints its own messages
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to start server: %w", err)
	}

	// Wait a moment for server to be ready
	time.Sleep(500 * time.Millisecond)

	return true, nil
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

// resolveSession resolves a session ID from args or interactive selection.
// If projectFlag is provided, filters sessions by that project.
// Returns sessionID, projectID, and error.
func resolveSession(cmd *cobra.Command, args []string, projectSvc *application.ProjectService, sessionSvc *application.SessionService) (sessionID, projectID string, err error) {
	projectFlag, _ := cmd.Flags().GetString("project")

	if len(args) == 0 {
		// Interactive mode: select session
		var allSessions []application.SessionWithProject

		if projectFlag != "" {
			// Filter by project
			project, err := projectSvc.GetProject(projectFlag)
			if err != nil {
				return "", "", err
			}
			sessions, err := sessionSvc.ListSessions(project.ID)
			if err != nil {
				return "", "", err
			}
			for _, s := range sessions {
				allSessions = append(allSessions, application.SessionWithProject{
					Session:     s,
					ProjectID:   project.ID,
					ProjectName: project.Name,
				})
			}
		} else {
			// List all sessions from all projects
			allSessions, err = sessionSvc.ListAllSessions()
			if err != nil {
				return "", "", err
			}
		}

		if len(allSessions) == 0 {
			if projectFlag != "" {
				return "", "", fmt.Errorf("no sessions found in project '%s'", projectFlag)
			}
			return "", "", fmt.Errorf("no sessions found")
		}

		options := make([]tui.SessionOption, len(allSessions))
		for i, s := range allSessions {
			options[i] = tui.SessionOption{
				SessionID:   s.Session.ID,
				ProjectID:   s.ProjectID,
				ProjectName: s.ProjectName,
				Session:     s.Session,
			}
		}

		sessionID, projectID, err = tui.SelectSessionFromAll(options)
		if err != nil {
			return "", "", err
		}
		return sessionID, projectID, nil
	}

	sessionID = args[0]

	if projectFlag != "" {
		project, err := projectSvc.GetProject(projectFlag)
		if err != nil {
			return "", "", err
		}
		projectID = project.ID
	} else {
		// Try to find the session across all projects
		allSessions, err := sessionSvc.ListAllSessions()
		if err != nil {
			return "", "", err
		}

		var matches []application.SessionWithProject
		for _, s := range allSessions {
			if s.Session.ID == sessionID || s.Session.Name == sessionID {
				matches = append(matches, s)
			}
		}

		if len(matches) == 0 {
			return "", "", fmt.Errorf("session not found: %s", sessionID)
		}

		if len(matches) > 1 {
			return "", "", fmt.Errorf("multiple sessions found with name '%s'. Use -p to specify project", sessionID)
		}

		projectID = matches[0].ProjectID
		sessionID = matches[0].Session.ID
	}

	return sessionID, projectID, nil
}
