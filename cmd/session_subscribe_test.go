package cmd_test

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestSessionSubscribeFilter(t *testing.T) {
	socket, _ := testSocketDir(t)

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains []string
		errExcludes []string
	}{
		{
			name:        "unknown provider filter",
			args:        []string{"session", "subscribe", "--provider", "nope"},
			wantErr:     true,
			errContains: []string{"unknown provider"},
		},
		{
			name:        "provider session source flags accepted",
			args:        []string{"session", "subscribe", "--provider", "claude-code", "--session", "s1", "--source", "hook", "--socket", socket},
			wantErr:     true,
			errExcludes: []string{"unknown provider"},
		},
		{
			name:        "omitted provider filter",
			args:        []string{"session", "subscribe", "--socket", socket},
			wantErr:     true,
			errExcludes: []string{"unknown provider"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executeRoot(t, execOpts{args: tt.args})
			require.Error(t, got.err)
			for _, want := range tt.errContains {
				assert.Contains(t, got.err.Error(), want)
			}
			for _, omit := range tt.errExcludes {
				assert.NotContains(t, got.err.Error(), omit)
			}
		})
	}
}

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
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT))
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("subscribe did not exit after interrupt")
	}
	assert.Contains(t, buf.String(), "hook/invoked")
}
