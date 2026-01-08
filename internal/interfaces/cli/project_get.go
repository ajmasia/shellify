package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var projectGetCmd = &cobra.Command{
	Use:     "get <id|name>",
	Aliases: []string{"show", "info"},
	Short:   "Get project details",
	Long: `Get details of a project by ID or name.

Examples:
  sfy project get my-project
  sfy project get abc12345
  sfy project get my-project --json`,
	Args: cobra.ExactArgs(1),
	RunE: runProjectGet,
}

func init() {
	projectCmd.AddCommand(projectGetCmd)
}

func runProjectGet(cmd *cobra.Command, args []string) error {
	idOrName := args[0]

	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	svc := application.NewProjectService(store)
	project, err := svc.GetProject(idOrName)
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
