package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/multiplexer"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionStopCmd = &cobra.Command{
	Use:     "stop [id|name]",
	Aliases: []string{"kill"},
	Short:   "Stop a running session",
	Long: `Stop a running terminal multiplexer session.

This kills the session and all its windows and panes.
If no session is specified, interactive mode will prompt for selection.

Examples:
  sfy session stop                        # Interactive mode
  sfy session stop my-session             # Stop by name
  sfy session stop my-session -p project  # Stop with specific project
  sfy session stop my-session --force     # Skip confirmation`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionStop,
}

func init() {
	sessionCmd.AddCommand(sessionStopCmd)
	sessionStopCmd.Flags().StringP("project", "p", "", "Project name or ID (for disambiguation)")
	sessionStopCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
}

func runSessionStop(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)
	generatorSvc := application.NewGeneratorService(store, store)
	launcherSvc, err := application.NewLauncherService(store, store, generatorSvc)
	if err != nil {
		return err
	}

	projectFlag, _ := cmd.Flags().GetString("project")
	forceFlag, _ := cmd.Flags().GetBool("force")

	var projectID, sessionID string

	if len(args) == 0 {
		// Interactive mode: select from running sessions
		allSessions, err := sessionSvc.ListAllSessions()
		if err != nil {
			return err
		}

		if len(allSessions) == 0 {
			return fmt.Errorf("no sessions found")
		}

		// Filter to show only running sessions
		launcher, err := multiplexer.NewLauncher()
		if err != nil {
			return err
		}

		var runningOptions []tui.SessionOption
		for _, s := range allSessions {
			fullSession, err := sessionSvc.GetSession(s.ProjectID, s.Session.ID)
			if err != nil {
				continue
			}
			if launcher.IsRunning(fullSession.SessionName, fullSession.TargetMultiplexer) {
				runningOptions = append(runningOptions, tui.SessionOption{
					SessionID:   s.Session.ID,
					ProjectID:   s.ProjectID,
					ProjectName: s.ProjectName,
					Session:     s.Session,
				})
			}
		}

		if len(runningOptions) == 0 {
			return fmt.Errorf("no running sessions found")
		}

		sessionID, projectID, err = tui.SelectSessionFromAll(runningOptions)
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

	// Get session info for feedback
	session, err := sessionSvc.GetSession(projectID, sessionID)
	if err != nil {
		return err
	}

	// Confirm before stopping
	if !forceFlag {
		confirmed, err := tui.Confirm(fmt.Sprintf("Stop session '%s'?", session.Name))
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Aborted")
			return nil
		}
	}

	if err := launcherSvc.StopSession(projectID, sessionID); err != nil {
		return err
	}

	fmt.Printf("Session '%s' stopped\n", session.Name)
	return nil
}
