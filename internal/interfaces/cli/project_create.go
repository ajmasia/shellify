package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var projectCreateCmd = &cobra.Command{
	Use:     "create <name>",
	Aliases: []string{"new", "add"},
	Short:   "Create a new project",
	Long: `Create a new project with the given name.

Examples:
  sfy project create my-project
  sfy project create "My Project" --description "A cool project"
  sfy project create my-project -d "Description" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runProjectCreate,
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)
	projectCreateCmd.Flags().StringP("description", "d", "", "Project description")
}

func runProjectCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	description, _ := cmd.Flags().GetString("description")

	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	svc := application.NewProjectService(store)
	project, err := svc.CreateProject(name, description)
	if err != nil {
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(project)
	}

	fmt.Printf("Project created: %s (%s)\n", project.Name, truncateID(project.ID))
	return nil
}
