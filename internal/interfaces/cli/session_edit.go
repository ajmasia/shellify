package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionEditCmd = &cobra.Command{
	Use:   "edit [id|name]",
	Short: "Edit a session",
	Long: `Edit a session using the visual GUI editor or directly in $EDITOR.

By default, opens the session JSON file in your $EDITOR for direct editing.
Use --gui flag to open the visual layout editor in your browser.

Examples:
  sfy session edit my-session              # Open JSON in $EDITOR
  sfy session edit my-session --gui        # Visual mode: open GUI editor
  sfy session edit -p my-project --gui     # GUI with project context`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionEdit,
}

func init() {
	sessionCmd.AddCommand(sessionEditCmd)
	sessionEditCmd.Flags().StringP("project", "p", "", "Project name or ID")
	sessionEditCmd.Flags().Bool("gui", false, "Open visual GUI editor in browser")

	// Shell completions
	_ = sessionEditCmd.RegisterFlagCompletionFunc("project", completeProjectNames)
	sessionEditCmd.ValidArgsFunction = completeSessionNames
}

func runSessionEdit(cmd *cobra.Command, args []string) error {
	guiMode, _ := cmd.Flags().GetBool("gui")

	if guiMode {
		return openGUIForSessionEdit(cmd, args)
	}

	return openAdvancedEdit(cmd, args)
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

	// Ensure server is running
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

func openAdvancedEdit(cmd *cobra.Command, args []string) error {
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

	// Show warning
	fmt.Println("\033[33m[WARN]\033[0m Advanced editing mode")
	fmt.Println("This will open the session JSON file directly. Invalid changes may break the session.")
	fmt.Println("I recommend using the visual editor: sfy session edit --gui")
	fmt.Println()

	confirmed, err := tui.Confirm("Continue?")
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Aborted")
		return nil
	}

	// Get the path to the session file
	paths, err := storage.NewPaths()
	if err != nil {
		return err
	}
	sessionFile := paths.SessionFile(projectID, sessionID)

	// Open in editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}
	// Expand environment variables in the path
	editor = os.ExpandEnv(editor)

	editorCmd := exec.Command(editor, sessionFile)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	fmt.Printf("Opening %s...\n", sessionFile)

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	fmt.Println("Session file saved.")
	return nil
}
