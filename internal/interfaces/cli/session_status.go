package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionStatusCmd = &cobra.Command{
	Use:   "status [id|name]",
	Short: "Show session status",
	Long: `Show the running status of a session.

Displays whether a session is running and if it has clients attached.
If no session is specified, interactive mode will prompt for selection.

Examples:
  sfy session status                        # Interactive mode
  sfy session status my-session             # Status by name
  sfy session status my-session -p project  # Status with specific project
  sfy session status my-session --json      # JSON output`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionStatus,
}

func init() {
	sessionCmd.AddCommand(sessionStatusCmd)
	sessionStatusCmd.Flags().StringP("project", "p", "", "Project name or ID (for disambiguation)")
	sessionStatusCmd.Flags().Bool("json", false, "Output in JSON format")
}

func runSessionStatus(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)
	generatorSvc := application.NewGeneratorService(store, store)
	launcherSvc, err := application.NewLauncherService(store, store, generatorSvc)
	if err != nil {
		return err
	}

	projectFlag, _ := cmd.Flags().GetString("project")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	var projectID, sessionID string

	if len(args) == 0 {
		// Interactive mode: select session from all projects
		allSessions, err := sessionSvc.ListAllSessions()
		if err != nil {
			return err
		}

		if len(allSessions) == 0 {
			return fmt.Errorf("no sessions found")
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
			return err
		}
	} else {
		sessionID = args[0]

		if projectFlag != "" {
			project, err := projectSvc.GetProject(projectFlag)
			if err != nil {
				return err
			}
			projectID = project.ID
		} else {
			// Try to find the session across all projects
			allSessions, err := sessionSvc.ListAllSessions()
			if err != nil {
				return err
			}

			var matches []application.SessionWithProject
			for _, s := range allSessions {
				if s.Session.ID == sessionID || s.Session.Name == sessionID {
					matches = append(matches, s)
				}
			}

			if len(matches) == 0 {
				return fmt.Errorf("session not found: %s", sessionID)
			}

			if len(matches) > 1 {
				return fmt.Errorf("multiple sessions found with name '%s'. Use -p to specify project", sessionID)
			}

			projectID = matches[0].ProjectID
			sessionID = matches[0].Session.ID
		}
	}

	status, err := launcherSvc.GetSessionStatus(projectID, sessionID)
	if err != nil {
		return err
	}

	// Get session info
	session, err := sessionSvc.GetSession(projectID, sessionID)
	if err != nil {
		return err
	}

	if jsonOutput {
		return printJSON(map[string]any{
			"name":        session.Name,
			"sessionName": status.SessionName,
			"multiplexer": status.Multiplexer,
			"running":     status.Running,
			"attached":    status.Attached,
		})
	}

	// Text output
	fmt.Printf("Session: %s\n", session.Name)
	fmt.Printf("Name:    %s\n", status.SessionName)
	fmt.Printf("Type:    %s\n", status.Multiplexer)

	if status.Running {
		fmt.Printf("Status:  Running\n")
		if status.Attached {
			fmt.Printf("Clients: Attached\n")
		} else {
			fmt.Printf("Clients: Detached\n")
		}
	} else {
		fmt.Printf("Status:  Stopped\n")
	}

	return nil
}
