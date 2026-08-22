package daemon_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/daemon"
)

type shutdownStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

func (shutdownStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (shutdownStub) Shutdown(context.Context, *agentdv1.ShutdownRequest) (*agentdv1.ShutdownResponse, error) {
	return &agentdv1.ShutdownResponse{}, nil
}

type shutdownFailStub struct {
	shutdownStub
}

func (shutdownFailStub) Shutdown(context.Context, *agentdv1.ShutdownRequest) (*agentdv1.ShutdownResponse, error) {
	return nil, errors.New("shutdown failed")
}

func TestStopTable(t *testing.T) {
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
				err := daemon.Stop(ctx, socket, time.Second)
				require.NoError(t, err, "Stop(not running)")
			},
		},
		{
			name: "clean shutdown",
			run: func(t *testing.T) {
				socket, cfg := testSocket(t)
				errCh := launchForegroundDaemon(t, socket, cfg)
				require.NoError(t, stopForegroundDaemon(t, socket, errCh), "Start(%q)", socket)

				paths := daemon.NewPaths(socket)
				_, err := paths.ReadPID()
				require.Error(t, err, "ReadPID() after stop")
				require.True(t, errors.Is(err, daemon.ErrNotRunning), "ReadPID() after stop")
			},
		},
		{
			name: "missing socket rpc",
			run: func(t *testing.T) {
				socket, _ := testSocket(t)
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				err := daemon.Stop(ctx, socket, time.Second)
				require.NoError(t, err, "Stop(missing socket)")
			},
		},
		{
			name: "stale dead pid",
			run: func(t *testing.T) {
				socket, _ := testSocket(t)
				paths := daemon.NewPaths(socket)
				require.NoError(t, paths.WritePID(999999999), "WritePID(stale)")

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				err := daemon.Stop(ctx, socket, time.Second)
				require.NoError(t, err, "Stop(stale dead pid)")
			},
		},
		{
			name: "own pid timeout",
			run: func(t *testing.T) {
				socket, _ := testSocket(t)
				paths := daemon.NewPaths(socket)
				require.NoError(t, paths.WritePID(os.Getpid()), "WritePID(self)")
				t.Cleanup(paths.ClearPID)

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				err := daemon.Stop(ctx, socket, 200*time.Millisecond)
				require.Error(t, err, "Stop(own pid timeout)")
				require.Contains(t, err.Error(), "foreground daemon did not exit", "Stop(own pid timeout)")
			},
		},
		{
			name: "shutdown rpc ok",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, shutdownStub{})
				})
				paths := daemon.NewPaths(socket)
				require.NoError(t, paths.WritePID(999999999), "WritePID(stale)")

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				err := daemon.Stop(ctx, socket, time.Second)
				require.NoError(t, err, "Stop(shutdown rpc)")
			},
		},
		{
			name: "shutdown rpc error ignored",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, shutdownFailStub{})
				})
				paths := daemon.NewPaths(socket)
				require.NoError(t, paths.WritePID(999999999), "WritePID(stale)")

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				err := daemon.Stop(ctx, socket, time.Second)
				require.NoError(t, err, "Stop(shutdown rpc error)")
			},
		},
		{
			name: "default socket not running",
			run: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				err := daemon.Stop(ctx, "", time.Second)
				require.NoError(t, err, "Stop(default socket)")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
