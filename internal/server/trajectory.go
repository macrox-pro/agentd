package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

type trajectoryService struct {
	agentdv1.UnimplementedTrajectoryServiceServer
	opts Options
}

func (s *trajectoryService) Statistics(_ context.Context, req *agentdv1.StatisticsRequest) (*agentdv1.StatisticsResponse, error) {
	if s.opts.Store == nil {
		return statistics.Response(s.opts.StartedAt, statistics.StatisticsRollup{}), nil
	}
	snap := s.opts.Store.Current()
	if err := statistics.Gate(snap.Trajectory); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}
	var rollup statistics.StatisticsRollup
	if s.opts.Collector != nil {
		rollup = s.opts.Collector.Snapshot(req.GetProvider())
	}
	return statistics.Response(s.opts.StartedAt, rollup), nil
}
