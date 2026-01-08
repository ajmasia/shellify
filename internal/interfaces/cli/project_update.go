package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var projectUpdateCmd = &cobra.Command{
	Use:     "update [id|name]",
	Aliases: []string{"edit"},
	Short:   "Update a project",
	Long: `Update a project's name or description.

If no project is specified, interactive mode will prompt for selection.
If no flags are provided, interactive mode will prompt for which fields to update.

Examples:
  sfy project update                              # Interactive mode (select + choose fields)
  sfy project update my-project                   # Interactive (choose fields)
  sfy project update my-project --name "New Name"
  sfy project update my-project --description "New description"
  sfy project update my-project -n "New" -d "Updated" --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectUpdate,
}

func init() {
	projectCmd.AddCommand(projectUpdateCmd)
	projectUpdateCmd.Flags().StringP("name", "n", "", "New project name")
	projectUpdateCmd.Flags().StringP("description", "d", "", "New project description")
}

func runProjectUpdate(cmd *cobra.Command, args []string) error {
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

	var name, description *string

	nameChanged := cmd.Flags().Changed("name")
	descChanged := cmd.Flags().Changed("description")

	if !nameChanged && !descChanged {
		// Interactive mode: ask what to update
		current, err := svc.GetProject(projectID)
		if err != nil {
			return err
		}

		updateName, updateDesc, err := tui.SelectUpdateFields()
		if err != nil {
			return err
		}

		if !updateName && !updateDesc {
			fmt.Println("No fields selected.")
			return nil
		}

		if updateName {
			newName, err := tui.PromptText("New name", current.Name)
			if err != nil {
				return err
			}
			if newName != current.Name {
				name = &newName
			}
		}

		if updateDesc {
			newDesc, err := tui.PromptText("New description", current.Description)
			if err != nil {
				return err
			}
			if newDesc != current.Description {
				description = &newDesc
			}
		}

		if name == nil && description == nil {
			fmt.Println("No changes made.")
			return nil
		}
	} else {
		if nameChanged {
			n, _ := cmd.Flags().GetString("name")
			name = &n
		}

		if descChanged {
			d, _ := cmd.Flags().GetString("description")
			description = &d
		}
	}

	project, err := svc.UpdateProject(projectID, name, description)
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
