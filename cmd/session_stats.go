package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

var (
	sessionStatsProvider string
	sessionStatsJSON     bool
)

func init() {
	sessionStatsCmd.Flags().StringVar(&sessionStatsProvider, "provider", "", "provider id (required)")
	sessionStatsCmd.Flags().BoolVar(&sessionStatsJSON, "json", false, "print JSON stats")
	_ = sessionStatsCmd.MarkFlagRequired("provider")
}

var sessionStatsCmd = &cobra.Command{
	Use:   "stats SESSION_ID",
	Short: "Show statistics for one session ledger",
	Long: `Scan a local session JSONL ledger and print aggregated statistics.

Requires trajectory.enabled and trajectory.statistics in config. Does not require a running daemon.`,
	Example: `  agentd session stats s1 --provider claude-code
  agentd session stats s1 --provider claude-code --json`,
	Args: cobra.ExactArgs(1),
	RunE: runSessionStats,
}

func runSessionStats(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	prov, err := provider.Parse(sessionStatsProvider)
	if err != nil {
		return err
	}
	stats, err := statistics.Load(ctx, statistics.SessionOptions{
		ConfigPath:   resolveConfigPath(),
		SessionsRoot: trajectory.DefaultSessionsDir(),
		Provider:     string(prov),
		SessionID:    args[0],
	})
	if err != nil {
		switch {
		case errors.Is(err, trajectory.ErrSessionNotFound):
			return fmt.Errorf("session %q not found", args[0])
		case errors.Is(err, statistics.ErrDisabled), errors.Is(err, statistics.ErrStatsOff):
			return mapStatisticsGateErr(err)
		default:
			return err
		}
	}
	return statistics.WriteSession(cmd.OutOrStdout(), stats, sessionStatsJSON)
}
