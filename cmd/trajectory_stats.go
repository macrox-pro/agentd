package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

var (
	trajectoryStatsProvider string
	trajectoryStatsJSON     bool
)

func init() {
	trajectoryStatsCmd.Flags().StringVar(&trajectoryStatsProvider, "provider", "", "optional provider filter")
	trajectoryStatsCmd.Flags().BoolVar(&trajectoryStatsJSON, "json", false, "print JSON response")
}

var trajectoryStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show daemon-lifetime trajectory statistics",
	Long: `Print in-memory trajectory counters from the running daemon.

Requires trajectory.enabled and trajectory.statistics in config, and a running daemon.`,
	Example: `  agentd trajectory stats
  agentd trajectory stats --provider claude-code --json`,
	RunE: runTrajectoryStats,
}

func runTrajectoryStats(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	prov, err := provider.ParseFilter(trajectoryStatsProvider, cmd.Flags().Changed("provider"))
	if err != nil {
		return err
	}
	store, err := config.Load(ctx, resolveConfigPath())
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if err := statistics.Gate(store.Current().Trajectory); err != nil {
		return mapStatisticsGateErr(err)
	}
	cli, err := hookclient.Dial(ctx, resolveSocket())
	if err != nil {
		return fmt.Errorf("daemon not running: %w", err)
	}
	defer cli.Close()

	req := &agentdv1.StatisticsRequest{}
	if prov != "" {
		p, err := prov.Proto()
		if err != nil {
			return err
		}
		req.Provider = p
	}
	resp, err := cli.Statistics(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			return mapStatisticsRPCGateErr(st.Message())
		}
		return fmt.Errorf("daemon not running: %w", err)
	}
	return statistics.WriteRollup(cmd.OutOrStdout(), resp, trajectoryStatsJSON)
}

func mapStatisticsGateErr(err error) error {
	switch {
	case errors.Is(err, statistics.ErrDisabled):
		return fmt.Errorf("trajectory is disabled; enable with: agentd config enable trajectory")
	case errors.Is(err, statistics.ErrStatsOff):
		return fmt.Errorf("trajectory statistics is disabled; enable with: agentd config enable trajectory-statistics")
	default:
		return err
	}
}

func mapStatisticsRPCGateErr(msg string) error {
	if strings.Contains(msg, statistics.ErrDisabled.Error()) {
		return fmt.Errorf("trajectory is disabled; enable with: agentd config enable trajectory")
	}
	if strings.Contains(msg, statistics.ErrStatsOff.Error()) {
		return fmt.Errorf("trajectory statistics is disabled; enable with: agentd config enable trajectory-statistics")
	}
	return fmt.Errorf("statistics unavailable: %s", msg)
}
