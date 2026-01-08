package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/multiplexer"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	"github.com/ajmasia/shellify/internal/interfaces/tui"
)

var sessionKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill running sessions",
	Long: `Kill one or more running multiplexer sessions.

An interactive multi-selector will be shown to choose which sessions to kill.
This is useful for stopping multiple sessions at once.

Examples:
  sfy session kill    # Select and kill multiple sessions`,
	RunE: runSessionKill,
}

func init() {
	sessionCmd.AddCommand(sessionKillCmd)
}

func runSessionKill(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	sessionSvc := application.NewSessionService(store, store)
	generatorSvc := application.NewGeneratorService(store, store)
	launcherSvc, err := application.NewLauncherService(store, store, generatorSvc)
	if err != nil {
		return err
	}

	launcher, err := multiplexer.NewLauncher()
	if err != nil {
		return err
	}

	// Get all sessions
	allSessions, err := sessionSvc.ListAllSessions()
	if err != nil {
		return err
	}

	if len(allSessions) == 0 {
		return fmt.Errorf("no sessions found")
	}

	// Filter to show only running sessions
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
		fmt.Println("No running sessions found")
		return nil
	}

	// Show multi-select
	selected, err := tui.SelectMultipleSessions(runningOptions)
	if err != nil {
		return err
	}

	// Kill selected sessions
	for _, opt := range selected {
		if err := launcherSvc.StopSession(opt.ProjectID, opt.SessionID); err != nil {
			fmt.Printf("Error killing %s: %v\n", opt.Session.Name, err)
		} else {
			fmt.Printf("Killed: %s\n", opt.Session.Name)
		}
	}

	return nil
}
