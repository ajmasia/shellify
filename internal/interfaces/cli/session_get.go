package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var sessionGetCmd = &cobra.Command{
	Use:     "get [id|name]",
	Aliases: []string{"show", "info"},
	Short:   "Get session details",
	Long: `Get details of a session by ID or name.

If no session is specified, interactive mode will prompt for selection.

Examples:
  sfy session get                                 # Interactive mode
  sfy session get my-session
  sfy session get my-session -p my-project
  sfy session get abc12345 --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionGet,
}

func init() {
	sessionCmd.AddCommand(sessionGetCmd)
	sessionGetCmd.Flags().StringP("project", "p", "", "Project name or ID (for disambiguation)")
}

func runSessionGet(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	sessionID, projectID, err := resolveSession(cmd, args, projectSvc, sessionSvc)
	if err != nil {
		return err
	}

	session, err := sessionSvc.GetSession(projectID, sessionID)
	if err != nil {
		return err
	}

	project, _ := projectSvc.GetProject(projectID)

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(session)
	}

	fmt.Printf("Name:        %s\n", session.Name)
	fmt.Printf("ID:          %s\n", session.ID)
	fmt.Printf("Session:     %s\n", session.SessionName)
	fmt.Printf("Project:     %s\n", project.Name)
	if session.Description != "" {
		fmt.Printf("Description: %s\n", session.Description)
	}
	fmt.Printf("Multiplexer: %s\n", session.TargetMultiplexer)
	if session.WorkingDirectory != "" {
		fmt.Printf("Directory:   %s\n", session.WorkingDirectory)
	}
	fmt.Printf("Windows:     %d\n", len(session.Windows))
	fmt.Printf("Created:     %s\n", session.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:     %s\n", formatTime(session.UpdatedAt))

	if len(session.Windows) > 0 {
		fmt.Println("\nWindows:")
		for i, w := range session.Windows {
			cmd := "(no command)"
			if len(w.Panes) > 0 && w.Panes[0].Command != "" {
				cmd = w.Panes[0].Command
			}
			fmt.Printf("  %d. %s: %s\n", i+1, w.Name, cmd)
		}
	}

	return nil
}
