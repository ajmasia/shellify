package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/multiplexer"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionAttachCmd = &cobra.Command{
	Use:   "attach [id|name]",
	Short: "Attach to a running session",
	Long: `Attach to a running terminal multiplexer session.

Only sessions that are currently running can be attached to.
If no session is specified, interactive mode will prompt for selection.

Examples:
  sfy session attach                        # Interactive mode
  sfy session attach my-session             # Attach by name
  sfy session attach my-session -p project  # Attach with specific project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSessionAttach,
}

func init() {
	sessionCmd.AddCommand(sessionAttachCmd)
	sessionAttachCmd.Flags().StringP("project", "p", "", "Project name or ID (for disambiguation)")

	// Shell completions
	_ = sessionAttachCmd.RegisterFlagCompletionFunc("project", completeProjectNames)
	sessionAttachCmd.ValidArgsFunction = completeSessionNames
}

func runSessionAttach(cmd *cobra.Command, args []string) error {
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
		// Interactive mode: select from running sessions
		var allSessions []application.SessionWithProject

		if projectFlag != "" {
			// Filter by project
			project, err := projectSvc.GetProject(projectFlag)
			if err != nil {
				return err
			}
			sessions, err := sessionSvc.ListSessions(project.ID)
			if err != nil {
				return err
			}
			for _, s := range sessions {
				allSessions = append(allSessions, application.SessionWithProject{
					Session:     s,
					ProjectID:   project.ID,
					ProjectName: project.Name,
				})
			}
		} else {
			allSessions, err = sessionSvc.ListAllSessions()
			if err != nil {
				return err
			}
		}

		if len(allSessions) == 0 {
			if projectFlag != "" {
				return fmt.Errorf("no sessions found in project '%s'", projectFlag)
			}
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
			if projectFlag != "" {
				return fmt.Errorf("no running sessions found in project '%s'", projectFlag)
			}
			return fmt.Errorf("no running sessions found")
		}

		sessionID, projectID, err = tui.SelectSessionFromAll(runningOptions)
		if err != nil {
			return err
		}
	} else {
		sessionID, projectID, err = resolveSession(cmd, args, projectSvc, sessionSvc)
		if err != nil {
			return err
		}
	}

	// Get session info for feedback
	session, err := sessionSvc.GetSession(projectID, sessionID)
	if err != nil {
		return err
	}

	fmt.Printf("Attaching to session '%s' (%s)...\n", session.Name, session.TargetMultiplexer)

	// This replaces the current process
	return launcherSvc.AttachSession(projectID, sessionID)
}
