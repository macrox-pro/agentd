package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the user-level agentd daemon",
	Long: `Lifecycle commands for the agentd daemon process.

The daemon holds the gRPC server, ConfigStore, Dispatch Engine, and async
queue. Hook events are handled by "agentd hook", not by daemon subcommands.`,
}

var (
	daemonForeground bool
	daemonStopTimeout string
	daemonStatusJSON bool
)

func init() {
	daemonCmd.AddCommand(daemonStartCmd, daemonStopCmd, daemonStatusCmd, daemonReloadCmd)

	daemonStartCmd.Flags().BoolVar(&daemonForeground, "foreground", false, "run in foreground (no detach)")

	daemonStopCmd.Flags().StringVar(&daemonStopTimeout, "timeout", "10s", "shutdown wait timeout")

	daemonStatusCmd.Flags().BoolVar(&daemonStatusJSON, "json", false, "JSON output")
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the user-level daemon",
	Long: `Start the agentd daemon (one instance per user).

By default the process detaches from the terminal. Use --foreground for
development or when supervised by systemd/launchd.

The daemon listens on a Unix domain socket (Linux/macOS) or named pipe
(Windows) for gRPC from hook and management CLI commands.`,
	Example: `  agentd daemon start
  agentd daemon start --foreground`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the daemon gracefully",
	Long: `Request graceful shutdown of the running daemon.

Drains in-flight hook Invoke calls and the async queue, then removes the
socket and PID files. Uses gRPC Shutdown when available; SIGTERM as fallback.`,
	Example: `  agentd daemon stop
  agentd daemon stop --timeout 30s`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show daemon runtime status",
	Long: `Print daemon health, uptime, config generation, and queue depth.

Use --json for machine-readable output in scripts and CI.

This shows runtime state, not declarative config (see "agentd config show").`,
	Example: `  agentd daemon status
  agentd daemon status --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}

var daemonReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload configuration from disk",
	Long: `Force a config re-merge from user and project YAML files.

Normally fsnotify reloads config automatically. Use this when the watcher
is unavailable (NFS, some containers) or after bulk edits.`,
	Example: `  agentd daemon reload`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
