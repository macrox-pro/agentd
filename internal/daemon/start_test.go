package daemon_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/daemon"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/transport"
)

type healthStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

func (healthStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func testSocket(t *testing.T) (socket, cfg string) {
	t.Helper()
	dir, err := os.MkdirTemp("", "agentd-daemon-")
	require.NoError(t, err, "MkdirTemp")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return filepath.Join(dir, "s.sock"), filepath.Join(dir, "missing.yaml")
}

func launchForegroundDaemon(t *testing.T, socket, cfg string) <-chan error {
	t.Helper()
	errCh := make(chan error, 1)
	go func() {
		errCh <- daemon.Start(t.Context(), daemon.StartOptions{
			Socket:     socket,
			ConfigPath: cfg,
			Foreground: true,
			Version:    "test",
		})
	}()
	require.NoError(t, waitReady(t, socket, errCh, 5*time.Second), "Start(%q)", socket)
	return errCh
}

func startForegroundDaemon(t *testing.T, socket, cfg string) <-chan error {
	t.Helper()
	errCh := launchForegroundDaemon(t, socket, cfg)
	t.Cleanup(func() { _ = stopForegroundDaemon(t, socket, errCh) })
	return errCh
}

func stopForegroundDaemon(t *testing.T, socket string, errCh <-chan error) error {
	t.Helper()
	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, daemon.Stop(stopCtx, socket, 10*time.Second), "Stop(%q)", socket)
	select {
	case err := <-errCh:
		return err
	case <-time.After(10 * time.Second):
		t.Fatal("Start did not return after Stop")
		return nil
	}
}

func TestStartStopTable(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "foreground start stop",
			run: func(t *testing.T) {
				socket, cfg := testSocket(t)
				errCh := launchForegroundDaemon(t, socket, cfg)
				require.NoError(t, stopForegroundDaemon(t, socket, errCh), "Start(%q)", socket)
			},
		},
		{
			name: "already running health",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterDaemonServiceServer(gs, healthStub{})
				})
				err := daemon.Start(context.Background(), daemon.StartOptions{
					Socket:     socket,
					ConfigPath: filepath.Join(filepath.Dir(socket), "missing.yaml"),
					Foreground: true,
					Version:    "test2",
				})
				require.ErrorIs(t, err, daemon.ErrAlreadyRunning, "Start(second)")
			},
		},
		{
			name: "live pid blocks start",
			run: func(t *testing.T) {
				socket, cfg := testSocket(t)
				paths := daemon.NewPaths(socket)
				require.NoError(t, paths.WritePID(os.Getpid()), "WritePID")
				t.Cleanup(paths.ClearPID)

				err := daemon.Start(context.Background(), daemon.StartOptions{
					Socket:     socket,
					ConfigPath: cfg,
					Foreground: true,
					Version:    "test",
				})
				require.ErrorIs(t, err, daemon.ErrAlreadyRunning, "Start(live pid)")
			},
		},
		{
			name: "stale pid then start",
			run: func(t *testing.T) {
				socket, cfg := testSocket(t)
				paths := daemon.NewPaths(socket)
				require.NoError(t, paths.WritePID(999999999), "WritePID(stale)")

				errCh := launchForegroundDaemon(t, socket, cfg)
				require.NoError(t, stopForegroundDaemon(t, socket, errCh), "Start(%q)", socket)
			},
		},
		{
			name: "lock held",
			run: func(t *testing.T) {
				socket, cfg := testSocket(t)
				paths := daemon.NewPaths(socket)
				lock, err := paths.AcquireLock()
				require.NoError(t, err, "AcquireLock()")
				t.Cleanup(func() { daemon.ReleaseLock(lock) })

				err = daemon.Start(context.Background(), daemon.StartOptions{
					Socket:     socket,
					ConfigPath: cfg,
					Foreground: true,
					Version:    "test",
				})
				require.Error(t, err, "Start(lock held)")
				require.True(t, errors.Is(err, daemon.ErrAlreadyRunning), "Start(lock held)")
			},
		},
		{
			name: "lock reacquired after stop",
			run: func(t *testing.T) {
				socket, cfg := testSocket(t)
				errCh := launchForegroundDaemon(t, socket, cfg)
				require.NoError(t, stopForegroundDaemon(t, socket, errCh), "Start(%q)", socket)

				paths := daemon.NewPaths(socket)
				lock, err := paths.AcquireLock()
				require.NoError(t, err, "AcquireLock() after stop")
				daemon.ReleaseLock(lock)
			},
		},
		{
			name: "listen error",
			run: func(t *testing.T) {
				dir, err := os.MkdirTemp("", "agentd-daemon-")
				require.NoError(t, err, "MkdirTemp")
				t.Cleanup(func() { _ = os.RemoveAll(dir) })
				block := filepath.Join(dir, "block")
				require.NoError(t, os.WriteFile(block, []byte("x"), 0o600), "WriteFile(%q)", block)
				socket := filepath.Join(block, "s.sock")
				cfg := filepath.Join(dir, "missing.yaml")

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				err = daemon.Start(ctx, daemon.StartOptions{
					Socket:     socket,
					ConfigPath: cfg,
					Foreground: true,
					Version:    "test",
				})
				require.Error(t, err, "Start(listen error)")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}

func startStubServer(t *testing.T, register func(*grpc.Server)) string {
	t.Helper()
	socket, _ := testSocket(t)
	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })

	gs := grpc.NewServer()
	register(gs)
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	require.NoError(t, waitReady(t, socket, nil, 2*time.Second), "stub ready")
	return socket
}

func waitReady(t *testing.T, socket string, errCh <-chan error, timeout time.Duration) error {
	t.Helper()
	deadline := time.After(timeout)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("daemon exited before ready: %w", err)
			}
			return fmt.Errorf("daemon exited before ready")
		case <-deadline:
			return fmt.Errorf("timeout waiting for health")
		case <-ticker.C:
			cli, err := hookclient.Dial(context.Background(), socket)
			if err != nil {
				continue
			}
			_, err = cli.Health(context.Background())
			_ = cli.Close()
			if err == nil {
				return nil
			}
		}
	}
}
