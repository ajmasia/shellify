package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var sessionExportCmd = &cobra.Command{
	Use:   "export [id|name]",
	Short: "Export a session as YAML",
	Long: `Export a session configuration as human-readable YAML.

The output can be redirected to a file, committed to a dotfiles repository,
and later restored on any machine with 'sfy session import'.

Examples:
  sfy session export my-session                        # Print YAML to stdout
  sfy session export my-session -o ./my-session.yaml  # Write to file
  sfy session export my-session -p my-project         # Disambiguate by project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionExport,
}

func init() {
	sessionCmd.AddCommand(sessionExportCmd)
	sessionExportCmd.Flags().StringP("project", "p", "", "Project name or ID")
	sessionExportCmd.Flags().StringP("output", "o", "", "Write output to file instead of stdout")

	_ = sessionExportCmd.RegisterFlagCompletionFunc("project", completeProjectNames)
	sessionExportCmd.ValidArgsFunction = completeSessionNames
}

func runSessionExport(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	sessionID, projectID, err := resolveSession(cmd, args, projectSvc, sessionSvc)
	if err != nil {
		return err
	}

	session, err := sessionSvc.GetSession(projectID, sessionID)
	if err != nil {
		return err
	}

	yamlData, err := SessionToYAML(session)
	if err != nil {
		return fmt.Errorf("serializing session: %w", err)
	}

	outputFlag, _ := cmd.Flags().GetString("output")
	if outputFlag == "" {
		fmt.Print(string(yamlData))
		return nil
	}

	if err := os.WriteFile(outputFlag, yamlData, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Session exported to %s\n", outputFlag)
	return nil
}
