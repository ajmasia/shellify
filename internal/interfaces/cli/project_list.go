package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var projectListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List all projects",
	RunE:    runProjectList,
}

func init() {
	projectCmd.AddCommand(projectListCmd)
}

func runProjectList(cmd *cobra.Command, _ []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	svc := application.NewProjectService(store)
	projects, err := svc.ListProjects()
	if err != nil {
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(projects)
	}

	if len(projects) == 0 {
		fmt.Println("No projects found. Create one with: sfy project create <name>")
		return nil
	}

	headers := []string{"ID", "NAME", "SESSIONS", "UPDATED"}
	rows := make([][]string, len(projects))

	for i, p := range projects {
		count, _ := svc.CountSessions(p.ID)
		rows[i] = []string{
			truncateID(p.ID),
			p.Name,
			fmt.Sprintf("%d", count),
			formatTime(p.UpdatedAt),
		}
	}

	printTable(headers, rows)
	return nil
}
