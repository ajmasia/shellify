package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/multiplexer"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var sessionListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List sessions",
	Long: `List all sessions, optionally filtered by project.

Examples:
  sfy session list                    # List all sessions
  sfy session list -p my-project      # List sessions in project
  sfy session list --json             # JSON output`,
	RunE: runSessionList,
}

func init() {
	sessionCmd.AddCommand(sessionListCmd)
	sessionListCmd.Flags().StringP("project", "p", "", "Filter by project name or ID")

	// Shell completions
	_ = sessionListCmd.RegisterFlagCompletionFunc("project", completeProjectNames)
}

func runSessionList(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	projectFilter, _ := cmd.Flags().GetString("project")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// Create launcher to check running status
	launcher, err := multiplexer.NewLauncher()
	if err != nil {
		return err
	}

	if projectFilter != "" {
		// List sessions for specific project
		project, err := projectSvc.GetProject(projectFilter)
		if err != nil {
			return err
		}

		sessions, err := sessionSvc.ListSessions(project.ID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(sessions)
		}

		if len(sessions) == 0 {
			fmt.Printf("No sessions in project '%s'\n", project.Name)
			return nil
		}

		headers := []string{"ID", "NAME", "MULTIPLEXER", "UPDATED", "STATUS"}
		rows := make([][]string, len(sessions))
		for i, s := range sessions {
			// Check if session is running
			fullSession, _ := sessionSvc.GetSession(project.ID, s.ID)
			status := "stopped"
			if launcher.IsRunning(fullSession.SessionName, fullSession.TargetMultiplexer) {
				status = "running"
			}

			rows[i] = []string{
				truncateID(s.ID),
				s.Name,
				string(s.TargetMultiplexer),
				formatTime(s.UpdatedAt),
				status,
			}
		}

		fmt.Printf("Sessions in '%s':\n\n", project.Name)
		printTable(headers, rows)
		return nil
	}

	// List all sessions from all projects
	allSessions, err := sessionSvc.ListAllSessions()
	if err != nil {
		return err
	}

	if jsonOutput {
		return printJSON(allSessions)
	}

	if len(allSessions) == 0 {
		fmt.Println("No sessions found")
		return nil
	}

	headers := []string{"ID", "NAME", "PROJECT", "MULTIPLEXER", "UPDATED", "STATUS"}
	rows := make([][]string, len(allSessions))
	for i, s := range allSessions {
		// Check if session is running
		fullSession, _ := sessionSvc.GetSession(s.ProjectID, s.Session.ID)
		status := "stopped"
		if launcher.IsRunning(fullSession.SessionName, fullSession.TargetMultiplexer) {
			status = "running"
		}

		rows[i] = []string{
			truncateID(s.Session.ID),
			s.Session.Name,
			s.ProjectName,
			string(s.Session.TargetMultiplexer),
			formatTime(s.Session.UpdatedAt),
			status,
		}
	}

	printTable(headers, rows)
	return nil
}
