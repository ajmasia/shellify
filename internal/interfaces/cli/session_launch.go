package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/multiplexer"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionLaunchCmd = &cobra.Command{
	Use:     "launch [id|name]",
	Aliases: []string{"run", "start"},
	Short:   "Launch a session",
	Long: `Launch a terminal multiplexer session.

If the session is already running, it will attach to it instead.
If no session is specified, interactive mode will prompt for selection.

Examples:
  sfy session launch                        # Interactive mode
  sfy session launch my-session             # Launch by name
  sfy session launch my-session -p project  # Launch with specific project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionLaunch,
}

func init() {
	sessionCmd.AddCommand(sessionLaunchCmd)
	sessionLaunchCmd.Flags().StringP("project", "p", "", "Project name or ID (for disambiguation)")
}

func runSessionLaunch(cmd *cobra.Command, args []string) error {
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

		// Filter to show only stopped sessions
		launcher, err := multiplexer.NewLauncher()
		if err != nil {
			return err
		}

		var stoppedOptions []tui.SessionOption
		for _, s := range allSessions {
			fullSession, err := sessionSvc.GetSession(s.ProjectID, s.Session.ID)
			if err != nil {
				continue
			}
			if !launcher.IsRunning(fullSession.SessionName, fullSession.TargetMultiplexer) {
				stoppedOptions = append(stoppedOptions, tui.SessionOption{
					SessionID:   s.Session.ID,
					ProjectID:   s.ProjectID,
					ProjectName: s.ProjectName,
					Session:     s.Session,
				})
			}
		}

		if len(stoppedOptions) == 0 {
			// All sessions are running, show all and let launch handle attach
			for _, s := range allSessions {
				stoppedOptions = append(stoppedOptions, tui.SessionOption{
					SessionID:   s.Session.ID,
					ProjectID:   s.ProjectID,
					ProjectName: s.ProjectName,
					Session:     s.Session,
				})
			}
		}

		sessionID, projectID, err = tui.SelectSessionFromAll(stoppedOptions)
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

	// Check if multiplexer is available
	avail := multiplexer.CheckMultiplexers()
	if !avail.HasMultiplexer(session.TargetMultiplexer) {
		return fmt.Errorf("%s is not installed.\n\n%s",
			session.TargetMultiplexer,
			multiplexer.GetInstallInstructions(session.TargetMultiplexer))
	}

	fmt.Printf("Launching session '%s' (%s)...\n", session.Name, session.TargetMultiplexer)

	// This replaces the current process
	return launcherSvc.LaunchSession(projectID, sessionID)
}
