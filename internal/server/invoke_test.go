package server_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func claudeToolPre(t *testing.T, command string) []byte {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"session_id":      "s",
		"cwd":             "/w",
		"hook_event_name": "PreToolUse",
		"tool_name":       "Bash",
		"tool_use_id":     "t1",
		"tool_input":      map[string]any{"command": command},
	})
	require.NoError(t, err)
	return b
}

func TestHookServiceInvoke(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)

	srv := server.New(server.Options{
		Store:     store,
		Engine:    eng,
		StartedAt: time.Now().UTC(),
		Version:   "test",
	})
	conn := dialBuf(t, srv)
	hook := agentdv1.NewHookServiceClient(conn)

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "invoke clean no decision",
			run: func(t *testing.T) {
				t.Helper()
				resp, err := hook.Invoke(ctx, &agentdv1.InvokeRequest{
					Provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload:     claudeToolPre(t, "go test"),
					InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
				})
				require.NoError(t, err, "Invoke()")
				assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_NO_DECISION, resp.GetDecision().GetKind(), "Invoke()")
				assert.GreaterOrEqual(t, resp.GetAsyncDispatchedCount(), uint32(1), "async")
			},
		},
		{
			name: "invoke secret ask",
			run: func(t *testing.T) {
				t.Helper()
				resp, err := hook.Invoke(ctx, &agentdv1.InvokeRequest{
					Provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload:     claudeToolPre(t, "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"),
					InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
				})
				require.NoError(t, err, "Invoke()")
				assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_ASK, resp.GetDecision().GetKind(), "Invoke()")
				assert.NotContains(t, resp.GetDecision().GetReason(), "AKIAIOSFODNN7EXAMPLE")
			},
		},
		{
			name: "invoke undecodable is neutral",
			run: func(t *testing.T) {
				t.Helper()
				resp, err := hook.Invoke(ctx, &agentdv1.InvokeRequest{
					Provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload:     []byte(`{}`),
					InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
				})
				require.NoError(t, err, "Invoke()")
				assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_NO_DECISION, resp.GetDecision().GetKind(), "Invoke()")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}

func TestHookServiceInvokeTrajectory(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateDir)

	path := filepath.Join(t.TempDir(), "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  enabled: true
`), 0o600))
	store, err := config.Load(ctx, path)
	require.NoError(t, err)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	recorder := trajectory.NewRecorder(trajectory.DefaultSessionsDir(), store.Current().Trajectory.QueueCapacity, nil)
	t.Cleanup(func() { recorder.Close(2 * time.Second) })

	srv := server.New(server.Options{
		Store:     store,
		Engine:    dispatch.NewEngine(q, nil),
		Recorder:  recorder,
		StartedAt: time.Now().UTC(),
		Version:   "test",
	})
	conn := dialBuf(t, srv)
	hook := agentdv1.NewHookServiceClient(conn)

	_, err = hook.Invoke(ctx, &agentdv1.InvokeRequest{
		Provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload:     claudeToolPre(t, "go test"),
		InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
		Cwd:            "/w",
	})
	require.NoError(t, err)

	time.Sleep(150 * time.Millisecond)
	summaries, err := trajectory.ListSessions(trajectory.DefaultSessionsDir(), "claude-code")
	require.NoError(t, err)
	require.NotEmpty(t, summaries)
	events, err := trajectory.ReadEvents(summaries[0].Path)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(events), 3)
	types := map[string]bool{}
	for _, e := range events {
		types[e.Type] = true
	}
	assert.True(t, types[trajectory.TypeSessionOpen])
	assert.True(t, types[trajectory.TypeHookInvoked])
	assert.True(t, types[trajectory.TypeHookDecided])
	for i, e := range events {
		assert.Equal(t, uint64(i+1), e.Seq)
	}
}
