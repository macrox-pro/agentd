package daemon

import (
	"context"
	"fmt"
	"os"
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

	_ = requestShutdown(stopCtx, socket)
	if waitUntilStopped(paths, deadline) {
		return nil
	}

	pid, err := paths.ReadPID()
	if err != nil {
		return err
	}
	if !processAlive(pid) {
		paths.RemoveStale()
		return nil
	}
	// Foreground mode records this process in the PID file; never SIGTERM ourselves
	// (that kills CLI/tests once signal.Notify has been stopped during shutdown).
	if pid == os.Getpid() {
		return fmt.Errorf("foreground daemon did not exit within %s", timeout)
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
		if err != nil {
			paths.RemoveStale()
			return true
		}
		// Own PID means foreground in this process; only the PID file going away
		// signals shutdown (processAlive would stay true until the test/CLI exits).
		if pid != os.Getpid() && !processAlive(pid) {
			paths.RemoveStale()
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}
