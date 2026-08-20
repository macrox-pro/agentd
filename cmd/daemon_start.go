package cmd

import (
	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/daemon"
	"github.com/macrox-pro/agentd/internal/version"
)

var daemonForeground bool

func init() {
	daemonStartCmd.Flags().BoolVar(&daemonForeground, "foreground", false, "keep the service attached to this terminal")
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the agentd background service",
	Long: `Start the agentd background service for this user.

By default the service runs in the background and start waits until the
service answers health checks. Use --foreground to keep it attached to the
terminal (useful while developing or under a process manager).

Only one instance should run per user. If a service is already running, start
reports an error instead of replacing it.`,
	Example: `  agentd daemon start
  agentd daemon start --foreground`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		return daemon.Start(cmd.Context(), daemon.StartOptions{
			Socket:     resolveSocket(),
			ConfigPath: resolveConfigPath(),
			Foreground: daemonForeground,
			Version:    version.Version,
		})
	},
}
