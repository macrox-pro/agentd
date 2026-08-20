package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dispatchCmd)
}

var dispatchRoutesJSON bool

var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Inspect dispatch routing",
	Long: `Debug and operations commands for the Dispatch Engine.

Shows compiled routes after config merge — match order, mode, and targets.`,
}

func init() {
	dispatchCmd.AddCommand(dispatchRoutesCmd)
	dispatchRoutesCmd.Flags().BoolVar(&dispatchRoutesJSON, "json", false, "JSON output")
}

var dispatchRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List compiled dispatch routes",
	Long: `Print dispatch routes from the running daemon snapshot.

Answers "why did this hook use parallel mode?" without reading source.
Not invoked by coding agents.`,
	Example: `  agentd dispatch routes
  agentd dispatch routes --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
