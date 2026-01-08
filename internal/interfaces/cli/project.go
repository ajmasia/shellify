package cli

import (
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"p"},
	Short:   "Manage projects",
	Long:    "Create, list, and manage projects. Projects are collections of sessions.",
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
