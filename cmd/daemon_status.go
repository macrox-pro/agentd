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
