package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionDeleteCmd = &cobra.Command{
	Use:     "delete [id|name]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a session",
	Long: `Delete a session.

If no session is specified, interactive mode will prompt for selection.
Prompts for confirmation unless --force is used.

Examples:
  sfy session delete                              # Interactive mode
  sfy session delete my-session
  sfy session delete my-session -p my-project
  sfy session delete my-session --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionDelete,
}

func init() {
	sessionCmd.AddCommand(sessionDeleteCmd)
	sessionDeleteCmd.Flags().StringP("project", "p", "", "Project name or ID")
	sessionDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
}

func runSessionDelete(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	projectFlag, _ := cmd.Flags().GetString("project")
	force, _ := cmd.Flags().GetBool("force")

	var projectID, sessionID, sessionName, projectName string

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

		// Find the selected session info
		for _, opt := range options {
			if opt.SessionID == sessionID {
				sessionName = opt.Session.Name
				projectName = opt.ProjectName
				break
			}
		}
	} else {
		sessionID = args[0]

		if projectFlag != "" {
			project, err := projectSvc.GetProject(projectFlag)
			if err != nil {
				return err
			}
			projectID = project.ID
			projectName = project.Name
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
			projectName = matches[0].ProjectName
			sessionID = matches[0].Session.ID
			sessionName = matches[0].Session.Name
		}

		// Get session name if not already set
		if sessionName == "" {
			session, err := sessionSvc.GetSession(projectID, sessionID)
			if err != nil {
				return err
			}
			sessionName = session.Name
		}
	}

	// Confirmation prompt unless --force is used
	if !force {
		message := fmt.Sprintf("Delete session '%s' from project '%s'?", sessionName, projectName)

		confirmed, err := tui.Confirm(message)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	err = sessionSvc.DeleteSession(projectID, sessionID)
	if err != nil {
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(map[string]any{
			"deleted": sessionID,
			"name":    sessionName,
			"project": projectName,
		})
	}

	fmt.Printf("Deleted session: %s (from %s)\n", sessionName, projectName)
	return nil
}
