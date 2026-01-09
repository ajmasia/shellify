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

	force, _ := cmd.Flags().GetBool("force")

	sessionID, projectID, err := resolveSession(cmd, args, projectSvc, sessionSvc)
	if err != nil {
		return err
	}

	// Get session and project details for confirmation message
	session, err := sessionSvc.GetSession(projectID, sessionID)
	if err != nil {
		return err
	}
	project, err := projectSvc.GetProject(projectID)
	if err != nil {
		return err
	}

	// Confirmation prompt unless --force is used
	if !force {
		message := fmt.Sprintf("Delete session '%s' from project '%s'?", session.Name, project.Name)

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
			"name":    session.Name,
			"project": project.Name,
		})
	}

	fmt.Printf("Deleted session: %s (from %s)\n", session.Name, project.Name)
	return nil
}
