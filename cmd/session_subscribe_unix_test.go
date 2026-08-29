//go:build unix

package cmd_test

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestSessionSubscribeStreamingHappy(t *testing.T) {
	socket, hub := startSubscribeServer(t)
	buf, done := executeRootAsync(t, execOpts{
		socketPath: socket,
		args:       []string{"session", "subscribe"},
	})
	publishHubEvent(hub, trajectory.Event{
		SchemaVersion: trajectory.SchemaVersion,
		Type:          trajectory.TypeHookInvoked,
		Source:        trajectory.SourceHook,
		Provider:      "claude-code",
		SessionID:     "s1",
		Seq:           1,
		TS:            time.Now().UTC(),
	})

	time.Sleep(200 * time.Millisecond)
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT), "Kill(SIGINT)")
	select {
	case err := <-done:
		require.NoError(t, err, "subscribe exit")
	case <-time.After(3 * time.Second):
		t.Fatal("subscribe did not exit after interrupt")
	}
	assert.Contains(t, buf.String(), "hook/invoked")
}
