package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var projectDeleteCmd = &cobra.Command{
	Use:     "delete [id|name]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a project",
	Long: `Delete a project and all its sessions.

If no project is specified, interactive mode will prompt for selection.
Prompts for confirmation unless --force is used.

Examples:
  sfy project delete                              # Interactive mode
  sfy project delete my-project
  sfy project delete my-project --force
  sfy project delete abc12345 --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectDelete,
}

func init() {
	projectCmd.AddCommand(projectDeleteCmd)
	projectDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation and force delete")
}

func runProjectDelete(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	svc := application.NewProjectService(store)

	var projectID, projectName string
	var sessionCount int

	if len(args) == 0 {
		// Interactive mode: select project
		projects, err := svc.ListProjects()
		if err != nil {
			return err
		}

		if len(projects) == 0 {
			return fmt.Errorf("no projects found")
		}

		selected, err := tui.SelectProject(projects)
		if err != nil {
			return err
		}

		projectID = selected.ID
		projectName = selected.Name
		sessionCount, _ = svc.CountSessions(projectID)
	} else {
		// Direct mode: get project by id or name
		project, err := svc.GetProject(args[0])
		if err != nil {
			return err
		}

		projectID = project.ID
		projectName = project.Name
		sessionCount, _ = svc.CountSessions(projectID)
	}

	// Confirmation prompt unless --force is used
	if !force {
		message := fmt.Sprintf("Delete project '%s'?", projectName)
		if sessionCount > 0 {
			message = fmt.Sprintf("Delete project '%s' with %d session(s)?", projectName, sessionCount)
		}

		confirmed, err := tui.Confirm(message)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	err = svc.DeleteProject(projectID, true)
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
			"deleted":  projectID,
			"name":     projectName,
			"sessions": sessionCount,
		})
	}

	if sessionCount > 0 {
		fmt.Printf("Deleted project: %s (%d session(s))\n", projectName, sessionCount)
	} else {
		fmt.Printf("Deleted project: %s\n", projectName)
	}

	return nil
}
