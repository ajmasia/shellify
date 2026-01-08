package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/domain"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionUpdateCmd = &cobra.Command{
	Use:     "update [id|name]",
	Aliases: []string{"edit"},
	Short:   "Update a session",
	Long: `Update a session's properties.

If no session is specified, interactive mode will prompt for selection.
If no flags are provided, interactive mode will prompt for which fields to update.

Examples:
  sfy session update                              # Interactive mode
  sfy session update my-session                   # Interactive field selection
  sfy session update my-session --name "New Name"
  sfy session update my-session -p proj --description "Updated"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionUpdate,
}

func init() {
	sessionCmd.AddCommand(sessionUpdateCmd)
	sessionUpdateCmd.Flags().StringP("project", "p", "", "Project name or ID")
	sessionUpdateCmd.Flags().StringP("name", "n", "", "New session name")
	sessionUpdateCmd.Flags().StringP("description", "d", "", "New description")
	sessionUpdateCmd.Flags().StringP("multiplexer", "m", "", "New multiplexer (tmux|zellij)")
	sessionUpdateCmd.Flags().StringP("working-dir", "w", "", "New working directory")
}

func runSessionUpdate(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	projectFlag, _ := cmd.Flags().GetString("project")

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

	// Get current session for interactive mode
	current, err := sessionSvc.GetSession(projectID, sessionID)
	if err != nil {
		return err
	}

	var input application.UpdateSessionInput

	nameChanged := cmd.Flags().Changed("name")
	descChanged := cmd.Flags().Changed("description")
	muxChanged := cmd.Flags().Changed("multiplexer")
	wdChanged := cmd.Flags().Changed("working-dir")

	if !nameChanged && !descChanged && !muxChanged && !wdChanged {
		// Interactive mode: ask what to update
		updateName, updateDesc, updateWD, updateMux, err := tui.SelectSessionUpdateFields()
		if err != nil {
			return err
		}

		if !updateName && !updateDesc && !updateWD && !updateMux {
			fmt.Println("No fields selected.")
			return nil
		}

		if updateName {
			newName, err := tui.PromptText("New name", current.Name)
			if err != nil {
				return err
			}
			if newName != current.Name {
				input.Name = &newName
			}
		}

		if updateDesc {
			newDesc, err := tui.PromptText("New description", current.Description)
			if err != nil {
				return err
			}
			if newDesc != current.Description {
				input.Description = &newDesc
			}
		}

		if updateWD {
			newWD, err := tui.PromptText("New working directory", current.WorkingDirectory)
			if err != nil {
				return err
			}
			if newWD != current.WorkingDirectory {
				input.WorkingDirectory = &newWD
			}
		}

		if updateMux {
			newMux, err := tui.SelectMultiplexer()
			if err != nil {
				return err
			}
			if newMux != current.TargetMultiplexer {
				input.Multiplexer = &newMux
			}
		}

		if input.Name == nil && input.Description == nil && input.WorkingDirectory == nil && input.Multiplexer == nil {
			fmt.Println("No changes made.")
			return nil
		}
	} else {
		if nameChanged {
			n, _ := cmd.Flags().GetString("name")
			input.Name = &n
		}

		if descChanged {
			d, _ := cmd.Flags().GetString("description")
			input.Description = &d
		}

		if wdChanged {
			w, _ := cmd.Flags().GetString("working-dir")
			input.WorkingDirectory = &w
		}

		if muxChanged {
			m, _ := cmd.Flags().GetString("multiplexer")
			mux := domain.MultiplexerType(m)
			input.Multiplexer = &mux
		}
	}

	session, err := sessionSvc.UpdateSession(projectID, sessionID, input)
	if err != nil {
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(session)
	}

	fmt.Printf("Session updated: %s\n", session.Name)
	return nil
}
