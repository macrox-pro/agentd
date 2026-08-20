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
	Short: "Inspect how hook events are routed",
	Long: `Inspect how the running service routes hook events.

Use these commands when debugging why an event took a particular path.`,
}

func init() {
	dispatchCmd.AddCommand(dispatchRoutesCmd)
	dispatchRoutesCmd.Flags().BoolVar(&dispatchRoutesJSON, "json", false, "print routes as JSON")
}

var dispatchRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List active hook routes",
	Long: `List the active hook routes loaded by the running service.

Shows match order and whether each route waits for a decision or runs in the
background. For operators and debugging; agents do not call this command.`,
	Example: `  agentd dispatch routes
  agentd dispatch routes --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
