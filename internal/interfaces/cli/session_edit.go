package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/domain"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

var sessionEditCmd = &cobra.Command{
	Use:   "edit [id|name]",
	Short: "Edit a session",
	Long: `Edit a session using the visual GUI editor or in $EDITOR as YAML.

By default, opens a clean YAML representation of the session in $EDITOR.
Use --gui to open the visual layout editor in the browser instead.

Examples:
  sfy session edit my-session              # Open YAML in $EDITOR
  sfy session edit my-session --gui        # Visual mode: open GUI editor
  sfy session edit -p my-project --gui     # GUI with project context`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionEdit,
}

func init() {
	sessionCmd.AddCommand(sessionEditCmd)
	sessionEditCmd.Flags().StringP("project", "p", "", "Project name or ID")
	sessionEditCmd.Flags().Bool("gui", false, "Open visual GUI editor in browser")

	_ = sessionEditCmd.RegisterFlagCompletionFunc("project", completeProjectNames)
	sessionEditCmd.ValidArgsFunction = completeSessionNames
}

func runSessionEdit(cmd *cobra.Command, args []string) error {
	guiMode, _ := cmd.Flags().GetBool("gui")
	if guiMode {
		return openGUIForSessionEdit(cmd, args)
	}
	return openYAMLEdit(cmd, args)
}

func openGUIForSessionEdit(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)

	sessionID, _, err := resolveSession(cmd, args, projectSvc, sessionSvc)
	if err != nil {
		return err
	}

	started, err := ensureServerRunning()
	if err != nil {
		return err
	}

	if started {
		fmt.Println("Server started at http://localhost:3777")
		fmt.Println("Use 'sfy server stop' when done")
		fmt.Println()
	}

	url := fmt.Sprintf("http://localhost:3777/sessions/%s/edit", sessionID)
	fmt.Printf("Opening: %s\n", url)
	openURL(url)

	return nil
}

func openYAMLEdit(cmd *cobra.Command, args []string) error {
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

	tmpFile, err := os.CreateTemp("", "sfy-session-*.yaml")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if _, err := tmpFile.Write(yamlData); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := openInEditor(tmpPath); err != nil {
		return err
	}

	edited, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("reading edited file: %w", err)
	}

	updated, err := YAMLToSession(edited)
	if err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	// Preserve the original ID and session name derivation
	updated.ID = session.ID
	updated.SessionName = session.SessionName

	input := yamlSessionToUpdateInput(updated)
	if _, err := sessionSvc.UpdateSession(projectID, sessionID, input); err != nil {
		return err
	}

	fmt.Println("Session updated.")
	return nil
}

func openInEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}
	editor = os.ExpandEnv(editor)

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// yamlSessionToUpdateInput converts a parsed domain.Session to an UpdateSessionInput.
func yamlSessionToUpdateInput(s domain.Session) application.UpdateSessionInput {
	windows := make([]application.WindowInput, len(s.Windows))
	for i, w := range s.Windows {
		windows[i] = application.WindowInput{
			Name:             w.Name,
			WorkingDirectory: w.WorkingDirectory,
			RootDirection:    w.RootDirection,
			Panes:            domainPanesToPaneInputs(w.Panes),
		}
	}

	return application.UpdateSessionInput{
		Name:             &s.Name,
		Description:      &s.Description,
		WorkingDirectory: &s.WorkingDirectory,
		Multiplexer:      &s.TargetMultiplexer,
		Environment:      s.Environment,
		PreCommands:      s.PreCommands,
		PostCommands:     s.PostCommands,
		Windows:          windows,
		DefaultWindowID:  &s.DefaultWindowID,
	}
}

func domainPanesToPaneInputs(panes []domain.Pane) []application.PaneInput {
	result := make([]application.PaneInput, len(panes))
	for i, p := range panes {
		result[i] = application.PaneInput{
			ID:               p.ID,
			Name:             p.Name,
			Command:          p.Command,
			WorkingDirectory: p.WorkingDirectory,
			Size:             p.Size,
			Direction:        p.Direction,
			Children:         domainPanesToPaneInputs(p.Children),
		}
	}
	return result
}
