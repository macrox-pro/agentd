package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/daemon"
)

var daemonStopTimeout string

func init() {
	daemonStopCmd.Flags().StringVar(&daemonStopTimeout, "timeout", "10s", "how long to wait for a clean shutdown")
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
