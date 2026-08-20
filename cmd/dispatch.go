package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dispatchCmd)
	dispatchCmd.AddCommand(dispatchRoutesCmd)
}

var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Inspect how hook events are routed",
	Long: `Inspect how the running service routes hook events.

Use these commands when debugging why an event took a particular path.`,
	Example: `  agentd dispatch routes
  agentd dispatch routes --json`,
}
