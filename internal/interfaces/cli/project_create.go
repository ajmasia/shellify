package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var projectCreateCmd = &cobra.Command{
	Use:     "create [name]",
	Aliases: []string{"new", "add"},
	Short:   "Create a new project",
	Long: `Create a new project with the given name.

If no name is provided, interactive mode will prompt for project details.

Examples:
  sfy project create                              # Interactive mode
  sfy project create my-project
  sfy project create "My Project" --description "A cool project"
  sfy project create my-project -d "Description" --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectCreate,
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)
	projectCreateCmd.Flags().StringP("description", "d", "", "Project description")
}

func runProjectCreate(cmd *cobra.Command, args []string) error {
	var name, description string

	if len(args) == 0 {
		// Interactive mode
		var err error
		name, description, err = tui.ProjectPrompt()
		if err != nil {
			return err
		}
	} else {
		name = args[0]
		description, _ = cmd.Flags().GetString("description")
	}

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
