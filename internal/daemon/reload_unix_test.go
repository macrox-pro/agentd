//go:build unix

package daemon_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestReloadSIGHUP(t *testing.T) {
	socket, cfg := testSocket(t)
	startForegroundDaemon(t, socket, cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	before, err := daemon.Status(ctx, socket)
	require.NoError(t, err, "Status(before)")
	require.True(t, before.Running, "Running")

	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGHUP), "Kill(SIGHUP)")

	deadline := time.Now().Add(3 * time.Second)
	var after daemon.StatusReport
	for time.Now().Before(deadline) {
		after, err = daemon.Status(ctx, socket)
		require.NoError(t, err, "Status(after)")
		if after.Generation > before.Generation {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	assert.Greater(t, after.Generation, before.Generation, "Generation after SIGHUP")
}
