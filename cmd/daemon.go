package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.AddCommand(daemonStartCmd, daemonStopCmd, daemonStatusCmd, daemonReloadCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start and manage the agentd background service",
	Long: `Commands to start, stop, check, and reload the agentd background service.

Start the service once per user, then use "agentd hook" from your agent
settings. Management commands talk to the running service.`,
	Example: `  agentd daemon start
  agentd daemon status`,
}
