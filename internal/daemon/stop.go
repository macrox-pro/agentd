package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/transport"
)

// Stop asks a running daemon to shut down.
func Stop(ctx context.Context, socket string, timeout time.Duration) error {
	if socket == "" {
		socket = transport.DefaultSocketPath()
	}
	paths := NewPaths(socket)
	deadline := time.Now().Add(timeout)

	stopCtx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	if err := requestShutdown(stopCtx, socket); err == nil {
		if waitUntilStopped(paths, deadline) {
			return nil
		}
	}

	pid, err := paths.ReadPID()
	if err != nil {
		return err
	}
	if !processAlive(pid) {
		paths.RemoveStale()
		return nil
	}
	if err := signalTerminate(pid); err != nil {
		return fmt.Errorf("signal: %w", err)
	}
	if waitUntilStopped(paths, deadline) {
		return nil
	}
	return fmt.Errorf("daemon pid %d did not exit within %s", pid, timeout)
}

func requestShutdown(ctx context.Context, socket string) error {
	cli, err := hookclient.Dial(ctx, socket)
	if err != nil {
		return err
	}
	defer cli.Close()
	return cli.Shutdown(ctx)
}

func waitUntilStopped(paths Paths, deadline time.Time) bool {
	for time.Now().Before(deadline) {
		pid, err := paths.ReadPID()
		if err != nil || !processAlive(pid) {
			paths.RemoveStale()
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}
