package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(trajectoryCmd)
	trajectoryCmd.AddCommand(trajectoryStatsCmd)
}

var trajectoryCmd = &cobra.Command{
	Use:   "trajectory",
	Short: "Trajectory statistics and rollups",
	Long: `Inspect daemon-lifetime trajectory counters from the running daemon.

Requires trajectory.enabled and trajectory.statistics. Offline session ledgers use session stats.`,
	Example: `  agentd trajectory stats
  agentd trajectory stats --provider claude-code --json`,
}
