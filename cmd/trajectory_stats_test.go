package cmd_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
	"github.com/macrox-pro/agentd/internal/transport"
)

func writeStatsConfig(t *testing.T, yaml string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0o600))
	return path
}

func startStatsServer(t *testing.T, cfgPath string, gateFail bool) string {
	t.Helper()
	store, err := config.Load(context.Background(), cfgPath)
	require.NoError(t, err)
	if gateFail {
		snap := store.Current()
		snap.Trajectory.Enabled = false
	}
	socket, _ := testSocketDir(t)
	collector := statistics.NewCollector()
	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	gs := server.New(server.Options{
		Store:     store,
		Engine:    dispatch.NewEngine(q, nil, nil),
		Collector: collector,
		StartedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Version:   "test",
	})
	ln, err := transport.Listen(socket)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitDaemonReady(t, socket)
	return socket
}

type statsGateStub struct {
	agentdv1.UnimplementedDaemonServiceServer
	agentdv1.UnimplementedTrajectoryServiceServer
	fail statisticsGateFail
}

type statisticsGateFail int

const (
	gateFailNone statisticsGateFail = iota
	gateFailTrajectory
	gateFailStatistics
)

func (s statsGateStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (s statsGateStub) Statistics(context.Context, *agentdv1.StatisticsRequest) (*agentdv1.StatisticsResponse, error) {
	switch s.fail {
	case gateFailTrajectory:
		return nil, status.Errorf(codes.FailedPrecondition, "trajectory disabled")
	case gateFailStatistics:
		return nil, status.Errorf(codes.FailedPrecondition, "trajectory statistics disabled")
	default:
		return statistics.Response(time.Now().UTC(), statistics.StatisticsRollup{}), nil
	}
}

func startStatsGateStub(t *testing.T, fail statisticsGateFail) string {
	t.Helper()
	socket, _ := testSocketDir(t)
	ln, err := transport.Listen(socket)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	gs := grpc.NewServer()
	agentdv1.RegisterDaemonServiceServer(gs, &statsGateStub{fail: fail})
	agentdv1.RegisterTrajectoryServiceServer(gs, &statsGateStub{fail: fail})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitDaemonReady(t, socket)
	return socket
}

func TestTrajectoryStatsCLI(t *testing.T) {
	enabledYAML := "version: 1\ntrajectory:\n  enabled: true\n  statistics: true\n"
	tests := []struct {
		name     string
		setup    func(t *testing.T) (cfg, sock string)
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name: "invalid_provider",
			setup: func(t *testing.T) (string, string) {
				cfg := writeStatsConfig(t, enabledYAML)
				return cfg, startStatsServer(t, cfg, false)
			},
			args:     []string{"trajectory", "stats", "--provider", "nope"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name: "trajectory_disabled",
			setup: func(t *testing.T) (string, string) {
				cfg := writeStatsConfig(t, "version: 1\ntrajectory:\n  enabled: false\n  statistics: true\n")
				return cfg, ""
			},
			args:     []string{"trajectory", "stats"},
			wantErr:  true,
			contains: "trajectory is disabled",
		},
		{
			name: "statistics_disabled",
			setup: func(t *testing.T) (string, string) {
				cfg := writeStatsConfig(t, "version: 1\ntrajectory:\n  enabled: true\n  statistics: false\n")
				return cfg, ""
			},
			args:     []string{"trajectory", "stats"},
			wantErr:  true,
			contains: "trajectory statistics is disabled",
		},
		{
			name: "daemon_down",
			setup: func(t *testing.T) (string, string) {
				cfg := writeStatsConfig(t, enabledYAML)
				return cfg, filepath.Join(t.TempDir(), "missing.sock")
			},
			args:     []string{"trajectory", "stats"},
			wantErr:  true,
			contains: "daemon not running",
		},
		{
			name: "json",
			setup: func(t *testing.T) (string, string) {
				cfg := writeStatsConfig(t, enabledYAML)
				return cfg, startStatsServer(t, cfg, false)
			},
			args:     []string{"trajectory", "stats", "--json"},
			contains: `"since"`,
		},
		{
			name: "rpc_trajectory_disabled",
			setup: func(t *testing.T) (string, string) {
				cfg := writeStatsConfig(t, enabledYAML)
				return cfg, startStatsGateStub(t, gateFailTrajectory)
			},
			args:     []string{"trajectory", "stats"},
			wantErr:  true,
			contains: "trajectory is disabled",
		},
		{
			name: "rpc_statistics_disabled",
			setup: func(t *testing.T) (string, string) {
				cfg := writeStatsConfig(t, enabledYAML)
				return cfg, startStatsGateStub(t, gateFailStatistics)
			},
			args:     []string{"trajectory", "stats"},
			wantErr:  true,
			contains: "trajectory statistics is disabled",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, sock := tt.setup(t)
			got := executeRoot(t, execOpts{args: tt.args, configPath: cfg, socketPath: sock})
			if tt.wantErr {
				require.Error(t, got.err)
				assert.Contains(t, got.err.Error(), tt.contains)
				return
			}
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
		})
	}
}
