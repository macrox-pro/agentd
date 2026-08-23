package cmd

import (
	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/daemon"
)

var daemonStatusJSON bool

func init() {
	daemonStatusCmd.Flags().BoolVar(&daemonStatusJSON, "json", false, "print status as JSON")
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show whether the service is running",
	Long: `Print whether the agentd service is running and basic health details.

When the service is running, the output includes that process's build version
from the Status RPC. That can differ from this CLI binary after an upgrade;
use "agentd version" for this binary without contacting the daemon.

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
		return daemon.WriteStatus(cmd.OutOrStdout(), rep, daemonStatusJSON)
	},
}
