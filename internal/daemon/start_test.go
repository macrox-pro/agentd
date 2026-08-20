package daemon_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestStartStopForeground(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")
	cfg := filepath.Join(dir, "missing.yaml")

	ctx := t.Context()

	errCh := make(chan error, 1)
	go func() {
		errCh <- daemon.Start(ctx, daemon.StartOptions{
			Socket:     socket,
			ConfigPath: cfg,
			Foreground: true,
			Version:    "test",
		})
	}()

	require.NoError(t, waitReady(t, socket, errCh, 3*time.Second), "Start(%q)", socket)

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	require.NoError(t, daemon.Stop(stopCtx, socket, 5*time.Second), "Stop(%q)", socket)

	select {
	case err := <-errCh:
		require.NoError(t, err, "Start(%q)", socket)
	case <-time.After(5 * time.Second):
		t.Fatal("Start did not return after Stop")
	}
}

func waitReady(t *testing.T, socket string, errCh <-chan error, timeout time.Duration) error {
	t.Helper()
	deadline := time.After(timeout)
	paths := daemon.NewPaths(socket)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("daemon exited before ready: %w", err)
			}
			return fmt.Errorf("daemon exited before writing pid")
		case <-deadline:
			return fmt.Errorf("timeout waiting for pid file")
		case <-ticker.C:
			if pid, err := paths.ReadPID(); err == nil && pid > 0 {
				return nil
			}
		}
	}
}
