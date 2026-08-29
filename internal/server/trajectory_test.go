package server_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

func TestStatisticsResponse_roundTrip(t *testing.T) {
	t.Parallel()
	since := time.Now().UTC().Truncate(time.Second)
	resp := statistics.Response(since, statistics.StatisticsRollup{
		HooksByKind: map[agentdv1.EventKind]uint64{
			agentdv1.EventKind_EVENT_KIND_TOOL_PRE: 3,
		},
	})
	require.NotNil(t, resp.GetSince())
	assert.Equal(t, since, resp.GetSince().AsTime())
	assert.Equal(t, uint64(3), resp.GetRollup().GetHooksByKind()[int32(agentdv1.EventKind_EVENT_KIND_TOOL_PRE)])
}

func TestStatisticsRPC(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	started := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		yaml      string
		provider  agentdv1.Provider
		observe   bool
		wantCode  codes.Code
		wantHooks uint64
		checkSince bool
	}{
		{
			name: "global",
			yaml: "version: 1\ntrajectory:\n  enabled: true\n  statistics: true\n",
			observe: true,
			wantHooks: 1,
		},
		{
			name: "filter_provider",
			yaml: "version: 1\ntrajectory:\n  enabled: true\n  statistics: true\n",
			provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			observe: true,
			wantHooks: 1,
		},
		{
			name: "unknown_provider",
			yaml: "version: 1\ntrajectory:\n  enabled: true\n  statistics: true\n",
			provider: agentdv1.Provider(999),
		},
		{
			name: "trajectory_disabled",
			yaml: "version: 1\ntrajectory:\n  enabled: false\n  statistics: true\n",
			wantCode: codes.FailedPrecondition,
		},
		{
			name: "statistics_disabled",
			yaml: "version: 1\ntrajectory:\n  enabled: true\n  statistics: false\n",
			wantCode: codes.FailedPrecondition,
		},
		{
			name: "since",
			yaml: "version: 1\ntrajectory:\n  enabled: true\n  statistics: true\n",
			checkSince: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "agentd.yaml")
			require.NoError(t, os.WriteFile(path, []byte(tt.yaml), 0o600))
			store, err := config.Load(ctx, path)
			require.NoError(t, err)

			collector := statistics.NewCollector()
			if tt.observe {
				collector.Observe(trajectoryRecordInput(store.Current()))
			}
			srv := server.New(server.Options{
				Store:     store,
				Collector: collector,
				StartedAt: started,
			})
			conn := dialBuf(t, srv)
			client := agentdv1.NewTrajectoryServiceClient(conn)
			resp, err := client.Statistics(ctx, &agentdv1.StatisticsRequest{Provider: tt.provider})
			if tt.wantCode != codes.OK {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				return
			}
			require.NoError(t, err)
			if tt.checkSince {
				assert.Equal(t, started, resp.GetSince().AsTime())
				return
			}
			if tt.wantHooks > 0 {
				got := resp.GetRollup().GetHooksByKind()[int32(agentdv1.EventKind_EVENT_KIND_TOOL_PRE)]
				assert.Equal(t, tt.wantHooks, got)
			}
		})
	}
}

func trajectoryRecordInput(snap *config.Snapshot) trajectory.RecordInput {
	return trajectory.RecordInput{
		Provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		Result: dispatch.InvokeResult{
			Meta:     dispatch.InvokeMeta{EventKind: "tool.pre"},
			Decision: &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_ALLOW},
		},
		Snap: snap,
	}
}
