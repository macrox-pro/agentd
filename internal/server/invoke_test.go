package server_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
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
