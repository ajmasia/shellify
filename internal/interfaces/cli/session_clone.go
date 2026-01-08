package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionCloneCmd = &cobra.Command{
	Use:   "clone [id|name]",
	Short: "Clone a session",
	Long: `Clone an existing session.

Creates a copy of the session with a new name.
If no session is specified, interactive mode will prompt for selection.

Examples:
  sfy session clone                              # Interactive mode
  sfy session clone my-session                   # Clone with auto-generated name
  sfy session clone my-session -n new-session    # Clone with specific name
  sfy session clone my-session -p project        # Clone with specific project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionClone,
}

func init() {
	sessionCmd.AddCommand(sessionCloneCmd)
	sessionCloneCmd.Flags().StringP("project", "p", "", "Project name or ID (for disambiguation)")
	sessionCloneCmd.Flags().StringP("name", "n", "", "New session name")
}

func runSessionClone(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	projectFlag, _ := cmd.Flags().GetString("project")
	nameFlag, _ := cmd.Flags().GetString("name")
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

	// Get new name interactively if not provided
	newName := nameFlag
	if newName == "" {
		session, err := sessionSvc.GetSession(projectID, sessionID)
		if err != nil {
			return err
		}
		defaultName := session.Name + " (copy)"
		newName, err = tui.PromptText("New session name", defaultName)
		if err != nil {
			return err
		}
	}

	cloned, err := sessionSvc.CloneSession(projectID, sessionID, newName)
	if err != nil {
		return err
	}

	if jsonOutput {
		return printJSON(cloned)
	}

	fmt.Printf("Session cloned: %s (%s)\n", cloned.Name, truncateID(cloned.ID))
	return nil
}
