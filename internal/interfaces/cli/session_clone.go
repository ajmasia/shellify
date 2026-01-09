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

	// Shell completions
	_ = sessionCloneCmd.RegisterFlagCompletionFunc("project", completeProjectNames)
	sessionCloneCmd.ValidArgsFunction = completeSessionNames
}

func runSessionClone(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	nameFlag, _ := cmd.Flags().GetString("name")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	sessionID, projectID, err := resolveSession(cmd, args, projectSvc, sessionSvc)
	if err != nil {
		return err
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
