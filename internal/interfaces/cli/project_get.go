package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var projectGetCmd = &cobra.Command{
	Use:     "get [id|name]",
	Aliases: []string{"show", "info"},
	Short:   "Get project details",
	Long: `Get details of a project by ID or name.

If no project is specified, interactive mode will prompt for selection.

Examples:
  sfy project get                                 # Interactive mode
  sfy project get my-project
  sfy project get abc12345
  sfy project get my-project --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectGet,
}

func init() {
	projectCmd.AddCommand(projectGetCmd)
}

func runProjectGet(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	svc := application.NewProjectService(store)

	var projectID string

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
	} else {
		projectID = args[0]
	}

	project, err := svc.GetProject(projectID)
	if err != nil {
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(project)
	}

	sessionCount, _ := svc.CountSessions(project.ID)

	fmt.Printf("Name:        %s\n", project.Name)
	fmt.Printf("ID:          %s\n", project.ID)
	if project.Description != "" {
		fmt.Printf("Description: %s\n", project.Description)
	}
	fmt.Printf("Prefix:      %s\n", project.SessionPrefix)
	fmt.Printf("Sessions:    %d\n", sessionCount)
	fmt.Printf("Created:     %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:     %s\n", formatTime(project.UpdatedAt))

	return nil
}
