package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/domain"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var sessionImportCmd = &cobra.Command{
	Use:   "import <file.yaml>",
	Short: "Import a session from a YAML file",
	Long: `Import a session configuration from a YAML file into a project.

Use '-' as the filename to read from stdin.

Examples:
  sfy session import ./my-session.yaml -p my-project
  sfy session import ./my-session.yaml -p my-project --name override-name
  cat my-session.yaml | sfy session import - -p my-project`,
	Args: cobra.ExactArgs(1),
	RunE: runSessionImport,
}

func init() {
	sessionCmd.AddCommand(sessionImportCmd)
	sessionImportCmd.Flags().StringP("project", "p", "", "Project name or ID (required)")
	sessionImportCmd.Flags().String("name", "", "Override the session name from the YAML")

	_ = sessionImportCmd.MarkFlagRequired("project")
	_ = sessionImportCmd.RegisterFlagCompletionFunc("project", completeProjectNames)
}

func runSessionImport(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	projectFlag, _ := cmd.Flags().GetString("project")
	nameFlag, _ := cmd.Flags().GetString("name")

	project, err := projectSvc.GetProject(projectFlag)
	if err != nil {
		return err
	}

	var yamlData []byte
	filePath := args[0]

	if filePath == "-" {
		yamlData, err = io.ReadAll(os.Stdin)
	} else {
		yamlData, err = os.ReadFile(filePath)
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	session, err := YAMLToSession(yamlData)
	if err != nil {
		return err
	}

	if strings.TrimSpace(nameFlag) != "" {
		session.Name = strings.TrimSpace(nameFlag)
	}

	input := application.CreateSessionInput{
		Name:             session.Name,
		Description:      session.Description,
		WorkingDirectory: session.WorkingDirectory,
		Multiplexer:      session.TargetMultiplexer,
		Environment:      session.Environment,
		PreCommands:      session.PreCommands,
		PostCommands:     session.PostCommands,
		Windows:          domainWindowsToWindowInputs(session.Windows),
		DefaultWindowID:  session.DefaultWindowID,
	}

	created, err := sessionSvc.CreateSession(project.ID, input)
	if err != nil {
		return err
	}

	fmt.Printf("Session imported: %s (id: %s)\n", created.Name, truncateID(created.ID))
	return nil
}

func domainWindowsToWindowInputs(windows []domain.Window) []application.WindowInput {
	result := make([]application.WindowInput, len(windows))
	for i, w := range windows {
		result[i] = application.WindowInput{
			Name:             w.Name,
			WorkingDirectory: w.WorkingDirectory,
			RootDirection:    w.RootDirection,
			Panes:            domainPanesToPaneInputs(w.Panes),
		}
	}
	return result
}
