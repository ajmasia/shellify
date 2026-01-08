package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/domain"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionCreateCmd = &cobra.Command{
	Use:     "create [name]",
	Aliases: []string{"new", "add"},
	Short:   "Create a new session",
	Long: `Create a new session in a project.

If no name is provided, interactive mode will prompt for session details.
If no project is specified, interactive mode will prompt for project selection.

Examples:
  sfy session create                              # Full interactive mode
  sfy session create -p my-project                # Interactive with project
  sfy session create my-session -p my-project     # Direct mode
  sfy session create dev -p proj -m tmux -w ~/dev`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionCreate,
}

func init() {
	sessionCmd.AddCommand(sessionCreateCmd)
	sessionCreateCmd.Flags().StringP("project", "p", "", "Project name or ID")
	sessionCreateCmd.Flags().StringP("description", "d", "", "Session description")
	sessionCreateCmd.Flags().StringP("multiplexer", "m", "tmux", "Target multiplexer (tmux|zellij)")
	sessionCreateCmd.Flags().StringP("working-dir", "w", "", "Working directory")
}

func runSessionCreate(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	// Step 1: Resolve project
	projectFlag, _ := cmd.Flags().GetString("project")
	var project domain.Project

	if projectFlag == "" {
		// Interactive: select project
		projects, err := projectSvc.ListProjects()
		if err != nil {
			return err
		}

		if len(projects) == 0 {
			return fmt.Errorf("no projects found. Create a project first with: sfy project create")
		}

		selected, err := tui.SelectProject(projects)
		if err != nil {
			return err
		}
		project = *selected
	} else {
		project, err = projectSvc.GetProject(projectFlag)
		if err != nil {
			return err
		}
	}

	// Step 2: Get session details
	var input application.CreateSessionInput

	if len(args) == 0 {
		// Interactive mode
		result, err := tui.SessionPrompt()
		if err != nil {
			return err
		}

		input = application.CreateSessionInput{
			Name:             result.Name,
			Description:      result.Description,
			WorkingDirectory: result.WorkingDirectory,
			Multiplexer:      result.Multiplexer,
			Windows:          make([]application.WindowInput, len(result.Windows)),
		}

		for i, w := range result.Windows {
			input.Windows[i] = application.WindowInput{
				Name:    w.Name,
				Command: w.Command,
			}
		}
	} else {
		// Direct mode
		input.Name = args[0]
		input.Description, _ = cmd.Flags().GetString("description")
		input.WorkingDirectory, _ = cmd.Flags().GetString("working-dir")

		mux, _ := cmd.Flags().GetString("multiplexer")
		input.Multiplexer = domain.MultiplexerType(mux)

		// Default window if none specified
		input.Windows = []application.WindowInput{
			{Name: "main", Command: ""},
		}
	}

	session, err := sessionSvc.CreateSession(project.ID, input)
	if err != nil {
		return err
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		return printJSON(session)
	}

	fmt.Printf("Session created: %s (%s) in project '%s'\n", session.Name, truncateID(session.ID), project.Name)
	fmt.Printf("Session name: %s\n", session.SessionName)
	fmt.Printf("Windows: %d\n", len(session.Windows))

	return nil
}
