package cli

import (
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"s"},
	Short:   "Manage sessions",
	Long:    "Create, list, and manage terminal multiplexer sessions.",
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
