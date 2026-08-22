package dispatch_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
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

func testSnap(t *testing.T) *config.Snapshot {
	t.Helper()
	_, async, guards, routes, err := config.Compile(nil)
	require.NoError(t, err)
	return &config.Snapshot{
		Generation: 1,
		Async:      async,
		Guards:     guards,
		Routes:     routes,
		Policy: config.Policy{
			Fail:        config.FailClosed,
			AskFallback: config.AskFallbackDeny,
		},
	}
}

func TestEngineInvokeClean(t *testing.T) {
	t.Parallel()
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 8,
		WorkerLimit:   2,
		TargetTimeout: time.Second,
	}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)
	snap := testSnap(t)

	res, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
		Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload: claudeToolPre(t, "go test ./..."),
		Snap:       snap,
	})
	require.NoError(t, err)
	require.NotNil(t, res.Decision)
	assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_NO_DECISION, res.Decision.Kind)
	assert.GreaterOrEqual(t, res.AsyncDispatchedCount, uint32(1), "parallel async observe")
}

func TestEngineInvokeSecretAsk(t *testing.T) {
	t.Parallel()
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 8,
		WorkerLimit:   2,
		TargetTimeout: time.Second,
	}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)
	snap := testSnap(t)

	res, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
		Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload: claudeToolPre(t, "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"),
		Snap:       snap,
	})
	require.NoError(t, err)
	require.NotNil(t, res.Decision)
	assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_ASK, res.Decision.Kind)
	assert.NotContains(t, res.Decision.Reason, "AKIAIOSFODNN7EXAMPLE")
}

func TestEngineParallelAsyncDoesNotBlock(t *testing.T) {
	t.Parallel()
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 8,
		WorkerLimit:   1,
		TargetTimeout: 5 * time.Second,
	}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })

	// Replace observe by filling queue with a slow job first — Invoke should still return quickly.
	block := make(chan struct{})
	require.True(t, q.Enqueue(dispatch.Job{Run: func(context.Context) { <-block }}))

	eng := dispatch.NewEngine(q, nil)
	snap := testSnap(t)
	start := time.Now()
	_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
		Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload: claudeToolPre(t, "echo hi"),
		Snap:       snap,
	})
	require.NoError(t, err)
	assert.Less(t, time.Since(start), time.Second, "Invoke must not wait on async workers")
	close(block)
}

func TestEngineInvokeCursorArgv(t *testing.T) {
	t.Parallel()
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 8,
		WorkerLimit:   2,
		TargetTimeout: time.Second,
	}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)
	snap := testSnap(t)

	_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
		Provider:       agentdv1.Provider_PROVIDER_CURSOR,
		RawPayload:     agenthookstest.Fixture(t, "cursor/pre_tool_use.json"),
		InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
		Snap:           snap,
		CWD:            "/work/repo",
	})
	require.NoError(t, err, "Invoke(cursor, argv)")
}

func TestMatchRouteDefaultKind(t *testing.T) {
	t.Parallel()
	snap := testSnap(t)
	typed, err := dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CLAUDE_CODE, agentdv1.InvocationMode_INVOCATION_MODE_STDIN, claudeToolPre(t, "x"))
	require.NoError(t, err)
	r := dispatch.MatchRoute(snap.Routes, typed)
	require.NotNil(t, r)
	assert.Equal(t, "tool.pre", r.Kind)
	assert.Equal(t, config.ModeParallel, r.Mode)
}

func TestEngineFileAsync(t *testing.T) {
	t.Parallel()
	audit := filepath.Join(t.TempDir(), "audit.jsonl")
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	content := "version: 1\ndispatch:\n  - name: audit\n    match:\n      kind: [tool.pre]\n    mode: parallel\n    sync:\n      - target: builtin\n        guards: [secrets]\n    async:\n      - target: file\n        path: " + audit + "\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	snap := store.Current()

	q := dispatch.NewQueue(config.AsyncConfig{QueueCapacity: 8, WorkerLimit: 2, TargetTimeout: time.Second}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)

	_, err = eng.Invoke(context.Background(), dispatch.InvokeInput{
		Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload: claudeToolPre(t, "go test"),
		Snap:       snap,
	})
	require.NoError(t, err)

	deadline := time.Now().Add(2 * time.Second)
	for {
		if b, err := os.ReadFile(audit); err == nil && len(b) > 0 {
			assert.Contains(t, string(b), "tool.pre")
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("audit file not written")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestDecodeTyped(t *testing.T) {
	t.Parallel()
	var n atomic.Int32
	_, err := dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CLAUDE_CODE, agentdv1.InvocationMode_INVOCATION_MODE_STDIN, claudeToolPre(t, "x"))
	require.NoError(t, err)
	n.Add(1)
	assert.Equal(t, int32(1), n.Load())
}
