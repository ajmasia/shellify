package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionGenerateCmd = &cobra.Command{
	Use:     "generate [id|name]",
	Aliases: []string{"gen"},
	Short:   "Generate session script or layout",
	Long: `Generate a script or layout file for a session.

For tmux sessions: generates an executable bash script (.sh)
For zellij sessions: generates a KDL layout file (.kdl)

If no output path is specified, the script is printed to stdout.
When writing to a file, tmux scripts are made executable (chmod +x).

Examples:
  sfy session generate                              # Interactive mode, output to stdout
  sfy session generate my-session                   # Output to stdout
  sfy session generate my-session -o ./start.sh    # Write to file
  sfy session generate my-session -p my-project    # Specify project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionGenerate,
}

func init() {
	sessionCmd.AddCommand(sessionGenerateCmd)
	sessionGenerateCmd.Flags().StringP("project", "p", "", "Project name or ID (for disambiguation)")
	sessionGenerateCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}

func runSessionGenerate(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)
	generatorSvc := application.NewGeneratorService(store, store)

	projectFlag, _ := cmd.Flags().GetString("project")
	outputFlag, _ := cmd.Flags().GetString("output")

	var projectID, sessionID string

	if len(args) == 0 {
		// Interactive mode: select session from all projects
		allSessions, err := sessionSvc.ListAllSessions()
		if err != nil {
			return err
		}

		if len(allSessions) == 0 {
			return fmt.Errorf("no sessions found")
		}

		options := make([]tui.SessionOption, len(allSessions))
		for i, s := range allSessions {
			options[i] = tui.SessionOption{
				SessionID:   s.Session.ID,
				ProjectID:   s.ProjectID,
				ProjectName: s.ProjectName,
				Session:     s.Session,
			}
		}

		sessionID, projectID, err = tui.SelectSessionFromAll(options)
		if err != nil {
			return err
		}
	} else {
		sessionID = args[0]

		if projectFlag != "" {
			project, err := projectSvc.GetProject(projectFlag)
			if err != nil {
				return err
			}
			projectID = project.ID
		} else {
			// Try to find the session across all projects
			allSessions, err := sessionSvc.ListAllSessions()
			if err != nil {
				return err
			}

			var matches []application.SessionWithProject
			for _, s := range allSessions {
				if s.Session.ID == sessionID || s.Session.Name == sessionID {
					matches = append(matches, s)
				}
			}

			if len(matches) == 0 {
				return fmt.Errorf("session not found: %s", sessionID)
			}

			if len(matches) > 1 {
				return fmt.Errorf("multiple sessions found with name '%s'. Use -p to specify project", sessionID)
			}

			projectID = matches[0].ProjectID
			sessionID = matches[0].Session.ID
		}
	}

	// Generate the script
	if outputFlag != "" {
		// Write to file
		filePath, err := generatorSvc.GenerateSessionToFile(projectID, sessionID, outputFlag)
		if err != nil {
			return fmt.Errorf("generating session: %w", err)
		}

		session, _ := sessionSvc.GetSession(projectID, sessionID)
		fmt.Printf("Generated %s script: %s\n", session.TargetMultiplexer, filePath)
		fmt.Printf("\nTo run:\n  %s\n", filePath)
	} else {
		// Output to stdout
		result, err := generatorSvc.GenerateSession(projectID, sessionID)
		if err != nil {
			return fmt.Errorf("generating session: %w", err)
		}

		_, _ = fmt.Fprint(os.Stdout, result.Content)
	}

	return nil
}
