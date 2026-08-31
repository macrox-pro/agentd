package daemon_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/daemon"
)

type statusFailStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

type statusOKStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

type statusMinimalStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

func (statusMinimalStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (statusMinimalStub) Status(context.Context, *agentdv1.StatusRequest) (*agentdv1.StatusResponse, error) {
	return &agentdv1.StatusResponse{Version: "bare"}, nil
}

func (statusOKStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (statusOKStub) Status(context.Context, *agentdv1.StatusRequest) (*agentdv1.StatusResponse, error) {
	return &agentdv1.StatusResponse{
		Version:                "test",
		StartedAt:              timestamppb.Now(),
		Config:                 &agentdv1.ConfigGeneration{Generation: 2, Fingerprint: "fp"},
		AsyncQueueDepth:        3,
		AsyncDroppedCount:      4,
		TrajectoryDroppedCount: 5,
		CompiledRouteCount:     6,
		MetricsListen:          "127.0.0.1:2112",
	}, nil
}

func (statusFailStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (statusFailStub) Status(context.Context, *agentdv1.StatusRequest) (*agentdv1.StatusResponse, error) {
	return nil, errors.New("status failed")
}

func TestStatusTable(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "not running",
			run: func(t *testing.T) {
				socket, _ := testSocket(t)
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				rep, err := daemon.Status(ctx, socket)
				require.NoError(t, err, "Status(%q)", socket)
				assert.False(t, rep.Running, "Running")
				assert.Equal(t, socket, rep.Socket, "Socket")
			},
		},
		{
			name: "running",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, statusOKStub{})
				})

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				rep, err := daemon.Status(ctx, socket)
				require.NoError(t, err, "Status(%q)", socket)
				assert.True(t, rep.Running, "Running")
				assert.Equal(t, "test", rep.Version, "Version")
				assert.Equal(t, uint64(2), rep.Generation, "Generation")
				assert.Equal(t, "fp", rep.Fingerprint, "Fingerprint")
				assert.Equal(t, uint32(3), rep.AsyncQueueDepth, "AsyncQueueDepth")
				assert.Equal(t, uint64(4), rep.AsyncDroppedCount, "AsyncDroppedCount")
				assert.Equal(t, uint64(5), rep.TrajectoryDroppedCount, "TrajectoryDroppedCount")
				assert.Equal(t, uint32(6), rep.CompiledRouteCount, "CompiledRouteCount")
				assert.Equal(t, "127.0.0.1:2112", rep.MetricsListen, "MetricsListen")
			},
		},
		{
			name: "running minimal fields",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, statusMinimalStub{})
				})
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				rep, err := daemon.Status(ctx, socket)
				require.NoError(t, err, "Status(%q)", socket)
				assert.True(t, rep.Running, "Running")
				assert.Equal(t, "bare", rep.Version, "Version")
				assert.True(t, rep.StartedAt.IsZero(), "StartedAt")
				assert.Equal(t, uint64(0), rep.Generation, "Generation")
			},
		},
		{
			name: "status rpc error",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, statusFailStub{})
				})
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				rep, err := daemon.Status(ctx, socket)
				require.NoError(t, err, "Status(%q)", socket)
				assert.False(t, rep.Running, "Running")
			},
		},
		{
			name: "default socket not running",
			run: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				rep, err := daemon.Status(ctx, "")
				require.NoError(t, err, "Status(default socket)")
				assert.False(t, rep.Running, "Running")
				assert.NotEmpty(t, rep.Socket, "Socket")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
