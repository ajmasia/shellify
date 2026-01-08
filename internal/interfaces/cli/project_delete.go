package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var projectDeleteCmd = &cobra.Command{
	Use:     "delete <id|name>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a project",
	Long: `Delete a project and all its sessions.

Use --force to delete a project that has sessions.

Examples:
  sfy project delete my-project
  sfy project delete my-project --force
  sfy project delete abc12345 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runProjectDelete,
}

func init() {
	projectCmd.AddCommand(projectDeleteCmd)
	projectDeleteCmd.Flags().BoolP("force", "f", false, "Force delete even if project has sessions")
}

func runProjectDelete(cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	force, _ := cmd.Flags().GetBool("force")

	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	svc := application.NewProjectService(store)

	// Get project first to show name in output
	project, err := svc.GetProject(idOrName)
	if err != nil {
		return err
	}

	sessionCount, _ := svc.CountSessions(project.ID)

	err = svc.DeleteProject(project.ID, force)
	if err != nil {
		var sessionsErr *application.ProjectHasSessionsError
		if errors.As(err, &sessionsErr) {
			return fmt.Errorf("project has %d session(s). Use --force to delete", sessionsErr.Count)
		}
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(map[string]any{
			"deleted":  project.ID,
			"name":     project.Name,
			"sessions": sessionCount,
		})
	}

	if sessionCount > 0 {
		fmt.Printf("Deleted project: %s (%d session(s))\n", project.Name, sessionCount)
	} else {
		fmt.Printf("Deleted project: %s\n", project.Name)
	}

	return nil
}
