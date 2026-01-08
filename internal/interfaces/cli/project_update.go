package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var projectUpdateCmd = &cobra.Command{
	Use:     "update <id|name>",
	Aliases: []string{"edit"},
	Short:   "Update a project",
	Long: `Update a project's name or description.

Examples:
  sfy project update my-project --name "New Name"
  sfy project update my-project --description "New description"
  sfy project update my-project -n "New" -d "Updated" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runProjectUpdate,
}

func init() {
	projectCmd.AddCommand(projectUpdateCmd)
	projectUpdateCmd.Flags().StringP("name", "n", "", "New project name")
	projectUpdateCmd.Flags().StringP("description", "d", "", "New project description")
}

func runProjectUpdate(cmd *cobra.Command, args []string) error {
	idOrName := args[0]

	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	var name, description *string

	if cmd.Flags().Changed("name") {
		n, _ := cmd.Flags().GetString("name")
		name = &n
	}

	if cmd.Flags().Changed("description") {
		d, _ := cmd.Flags().GetString("description")
		description = &d
	}

	if name == nil && description == nil {
		return fmt.Errorf("at least one of --name or --description is required")
	}

	svc := application.NewProjectService(store)
	project, err := svc.UpdateProject(idOrName, name, description)
	if err != nil {
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(project)
	}

	fmt.Printf("Project updated: %s\n", project.Name)
	return nil
}
