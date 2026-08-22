package daemon_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/daemon"
)

type reloadFailStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

type reloadOKStub struct {
	agentdv1.UnimplementedDaemonServiceServer
	gen uint64
}

func (s *reloadOKStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (s *reloadOKStub) ReloadConfig(context.Context, *agentdv1.ReloadConfigRequest) (*agentdv1.ReloadConfigResponse, error) {
	s.gen++
	if s.gen == 0 {
		s.gen = 1
	}
	return &agentdv1.ReloadConfigResponse{
		Config: &agentdv1.ConfigGeneration{Generation: s.gen, Fingerprint: "fp"},
	}, nil
}

type reloadEmptyStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

func (reloadEmptyStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (reloadEmptyStub) ReloadConfig(context.Context, *agentdv1.ReloadConfigRequest) (*agentdv1.ReloadConfigResponse, error) {
	return &agentdv1.ReloadConfigResponse{}, nil
}

func (reloadFailStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (reloadFailStub) ReloadConfig(context.Context, *agentdv1.ReloadConfigRequest) (*agentdv1.ReloadConfigResponse, error) {
	return nil, errors.New("reload failed")
}

func TestReloadRPC(t *testing.T) {
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
				_, err := daemon.Reload(ctx, socket)
				require.Error(t, err, "Reload(not running)")
			},
		},
		{
			name: "rapid reload ok",
			run: func(t *testing.T) {
				stub := &reloadOKStub{}
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, stub)
				})

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				first, err := daemon.Reload(ctx, socket)
				require.NoError(t, err, "Reload(first)")
				second, err := daemon.Reload(ctx, socket)
				require.NoError(t, err, "Reload(second)")
				assert.Equal(t, uint64(1), first.Generation, "Generation(first)")
				assert.Equal(t, "fp", first.Fingerprint, "Fingerprint(first)")
				assert.Greater(t, second.Generation, first.Generation, "Generation(second)")
			},
		},
		{
			name: "reload rpc error",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, reloadFailStub{})
				})
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				_, err := daemon.Reload(ctx, socket)
				require.Error(t, err, "Reload(rpc error)")
			},
		},
		{
			name: "reload empty config",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, reloadEmptyStub{})
				})
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				out, err := daemon.Reload(ctx, socket)
				require.NoError(t, err, "Reload(empty config)")
				assert.Equal(t, uint64(0), out.Generation, "Generation")
				assert.Empty(t, out.Fingerprint, "Fingerprint")
			},
		},
		{
			name: "default socket not running",
			run: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				_, err := daemon.Reload(ctx, "")
				require.Error(t, err, "Reload(default socket)")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
