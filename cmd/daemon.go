package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.AddCommand(daemonStartCmd, daemonStopCmd, daemonStatusCmd, daemonReloadCmd, daemonEnableCmd, daemonDisableCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start and manage the agentd background service",
	Long: `Commands to start, stop, check, reload, and configure login autostart for
the agentd background service.

Start the service once per user, then use "agentd hook" from your agent
settings. Use "agentd daemon enable" to start agentd automatically when
you log in.`,
	Example: `  agentd daemon start
  agentd daemon enable
  agentd daemon status`,
}
