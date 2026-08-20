package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/daemon"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/transport"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start and manage the agentd background service",
	Long: `Commands to start, stop, check, and reload the agentd background service.

Start the service once per user, then use "agentd hook" from your agent
settings. Management commands talk to the running service.`,
}

var (
	daemonForeground  bool
	daemonStopTimeout string
	daemonStatusJSON  bool
)

func init() {
	daemonCmd.AddCommand(daemonStartCmd, daemonStopCmd, daemonStatusCmd, daemonReloadCmd)

	daemonStartCmd.Flags().BoolVar(&daemonForeground, "foreground", false, "keep the service attached to this terminal")

	daemonStopCmd.Flags().StringVar(&daemonStopTimeout, "timeout", "10s", "how long to wait for a clean shutdown")

	daemonStatusCmd.Flags().BoolVar(&daemonStatusJSON, "json", false, "print status as JSON")
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
			Version:    "dev",
		})
	},
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running agentd service",
	Long: `Ask the running agentd service to shut down cleanly.

The command waits up to --timeout for the service to exit. If the service
does not respond, stop may force termination.`,
	Example: `  agentd daemon stop
  agentd daemon stop --timeout 30s`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		timeout, err := time.ParseDuration(daemonStopTimeout)
		if err != nil {
			return fmt.Errorf("invalid --timeout: %w", err)
		}
		return daemon.Stop(cmd.Context(), resolveSocket(), timeout)
	},
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show whether the service is running",
	Long: `Print whether the agentd service is running and basic health details.

Use --json for machine-readable output in scripts. For configuration contents,
use "agentd config show".`,
	Example: `  agentd daemon status
  agentd daemon status --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		rep, err := daemon.Status(cmd.Context(), resolveSocket())
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()
		if daemonStatusJSON {
			payload := map[string]any{
				"running": rep.Running,
				"socket":  rep.Socket,
			}
			if rep.Running {
				payload["version"] = rep.Version
				payload["started_at"] = rep.StartedAt.UTC().Format(time.RFC3339)
				payload["generation"] = rep.Generation
				payload["fingerprint"] = rep.Fingerprint
				payload["async_queue_depth"] = rep.AsyncQueueDepth
				payload["compiled_route_count"] = rep.CompiledRouteCount
			}
			return json.NewEncoder(out).Encode(payload)
		}
		if !rep.Running {
			fmt.Fprintln(out, "agentd: not running")
			return nil
		}
		fmt.Fprintf(out, "agentd: running (version %s, generation %d)\n",
			rep.Version, rep.Generation)
		return nil
	},
}

var daemonReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload settings from disk",
	Long: `Reload agentd settings from the config file without restarting the service.

Use this after editing your config by hand when you want changes applied
immediately.`,
	Example: `  agentd daemon reload`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		cli, err := hookclient.Dial(cmd.Context(), resolveSocket())
		if err != nil {
			return fmt.Errorf("daemon not running: %w", err)
		}
		defer cli.Close()
		resp, err := cli.Reload(cmd.Context())
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "reloaded generation=%d fingerprint=%s\n",
			resp.GetConfig().GetGeneration(), resp.GetConfig().GetFingerprint())
		return nil
	},
}

func resolveSocket() string {
	if socketPath != "" {
		return socketPath
	}
	return transport.DefaultSocketPath()
}

func resolveConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".agentd.yaml")
}
