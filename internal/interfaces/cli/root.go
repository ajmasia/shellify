// Package cli provides the command-line interface for shellify.
package cli

import (
	"github.com/spf13/cobra"
)

// Version information set via ldflags at build time.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "sfy",
	Short: "Shellify - Terminal multiplexer session manager",
	Long: `Shellify (sfy) is a CLI tool for managing tmux and zellij sessions.

Create, configure, and launch terminal multiplexer sessions with ease.
Use the CLI for quick operations or start the web UI for visual editing.

Examples:
  sfy version             Show version information
  sfy project list        List all projects
  sfy session list        List all sessions
  sfy session launch dev  Launch a session`,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
