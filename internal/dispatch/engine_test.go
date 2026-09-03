package dispatch_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net"
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
	"github.com/macrox-pro/agentd/internal/transport"
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
	res, err := config.CompileMerged(nil, nil, nil)
	require.NoError(t, err)
	return &config.Snapshot{
		Generation: 1,
		Async:      res.Async,
		Guards:     res.Guards,
		Routes:     res.Routes,
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
	eng := dispatch.NewEngine(q, nil, nil)
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
	eng := dispatch.NewEngine(q, nil, nil)
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

	eng := dispatch.NewEngine(q, nil, nil)
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
	eng := dispatch.NewEngine(q, nil, nil)
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
	eng := dispatch.NewEngine(q, nil, nil)

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

func paritySnap(t *testing.T, sync []config.CompiledTarget, async []config.CompiledTarget, mode config.DispatchMode) *config.Snapshot {
	t.Helper()
	base := testSnap(t)
	base.Routes = []config.CompiledRoute{{
		Name:  "parity",
		Match: config.RouteMatch{Kinds: []string{"tool.pre"}},
		Mode:  mode,
		Sync:  sync,
		Async: async,
	}}
	return base
}

func TestEngine_RunSyncParity(t *testing.T) {
	t.Parallel()
	missingPeer := "unix://" + filepath.Join(t.TempDir(), "missing.sock")

	tests := []struct {
		name     string
		sync     []config.CompiledTarget
		command  string
		wantKind agentdv1.DecisionKind
	}{
		{
			name: "builtin_deny_first",
			sync: []config.CompiledTarget{{
				Kind:   config.TargetBuiltin,
				Guards: []string{"secrets"},
			}},
			command:  "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
			wantKind: agentdv1.DecisionKind_DECISION_KIND_ASK,
		},
		{
			name: "grpc_fail_open",
			sync: []config.CompiledTarget{{
				Kind:     config.TargetGRPC,
				Endpoint: missingPeer,
				Timeout:  200 * time.Millisecond,
				OnError:  config.FailOpen,
			}},
			command:  "echo ok",
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
		},
		{
			name: "grpc_fail_closed",
			sync: []config.CompiledTarget{{
				Kind:     config.TargetGRPC,
				Endpoint: missingPeer,
				Timeout:  200 * time.Millisecond,
				OnError:  config.FailClosed,
			}},
			command:  "echo ok",
			wantKind: agentdv1.DecisionKind_DECISION_KIND_DENY,
		},
		{
			name: "skip_non_sync_kind",
			sync: []config.CompiledTarget{
				{Kind: config.TargetLog},
				{Kind: config.TargetBuiltin, Guards: []string{"secrets"}},
			},
			command:  "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
			wantKind: agentdv1.DecisionKind_DECISION_KIND_ASK,
		},
		{
			name: "first_conclusive_stops",
			sync: []config.CompiledTarget{
				{Kind: config.TargetBuiltin, Guards: []string{"secrets"}},
				{
					Kind:     config.TargetGRPC,
					Endpoint: missingPeer,
					Timeout:  200 * time.Millisecond,
					OnError:  config.FailClosed,
				},
			},
			command:  "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
			wantKind: agentdv1.DecisionKind_DECISION_KIND_ASK,
		},
		{
			name:     "empty_sync_list_neutral",
			sync:     nil,
			command:  "echo ok",
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			eng := dispatch.NewEngine(nil, nil, nil)
			snap := paritySnap(t, tt.sync, nil, config.ModeSyncOnly)
			res, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
				Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
				RawPayload: claudeToolPre(t, tt.command),
				Snap:       snap,
			})
			require.NoError(t, err)
			require.NotNil(t, res.Decision)
			assert.Equal(t, tt.wantKind, res.Decision.Kind)
			assert.Equal(t, uint32(0), res.AsyncDispatchedCount)
		})
	}
}

func TestEngineHybrid(t *testing.T) {
	t.Parallel()
	audit := filepath.Join(t.TempDir(), "audit.jsonl")
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 8,
		WorkerLimit:   2,
		TargetTimeout: time.Second,
	}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil, nil)
	snap := paritySnap(t,
		[]config.CompiledTarget{{Kind: config.TargetBuiltin, Guards: []string{"secrets"}}},
		[]config.CompiledTarget{{Kind: config.TargetFile, Path: audit}},
		config.ModeParallel,
	)

	res, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
		Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload: claudeToolPre(t, "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"),
		Snap:       snap,
	})
	require.NoError(t, err)
	require.NotNil(t, res.Decision)
	assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_ASK, res.Decision.Kind)
	assert.GreaterOrEqual(t, res.AsyncDispatchedCount, uint32(1))

	deadline := time.Now().Add(2 * time.Second)
	for {
		if b, err := os.ReadFile(audit); err == nil && len(b) > 0 {
			assert.Contains(t, string(b), "tool.pre")
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("async audit file not written; sync decision must not wait but async should still run")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func syncFailureSnap(t *testing.T, fail config.FailMode, mode config.DispatchMode, kinds []string, async []config.CompiledTarget) *config.Snapshot {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "agentd-hang-")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	sock := filepath.Join(dir, "s.sock")
	ln, err := transport.Listen(sock)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			c, acceptErr := ln.Accept()
			if acceptErr != nil {
				return
			}
			go func(conn net.Conn) {
				time.Sleep(time.Minute)
				_ = conn.Close()
			}(c)
		}
	}()
	snap := testSnap(t)
	snap.Policy.Fail = fail
	snap.Routes = []config.CompiledRoute{{
		Name:  "sync-failure",
		Match: config.RouteMatch{Kinds: kinds},
		Mode:  mode,
		Sync: []config.CompiledTarget{{
			Kind:     config.TargetGRPC,
			Endpoint: "unix://" + sock,
			Timeout:  500 * time.Millisecond,
			OnError:  config.FailOpen,
		}},
		Async: async,
	}}
	return snap
}

func syncFailureDeadline() time.Time {
	return time.Now().Add(50 * time.Millisecond)
}

func claudePromptSubmitted(t *testing.T) []byte {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"session_id":      "s",
		"cwd":             "/w",
		"hook_event_name": "UserPromptSubmit",
		"prompt":          "hello",
	})
	require.NoError(t, err)
	return b
}

func TestEngineSyncFailurePolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		fail           config.FailMode
		mode           config.DispatchMode
		provider       agentdv1.Provider
		invocationMode agentdv1.InvocationMode
		raw            func(t *testing.T) []byte
		kinds          []string
		async          []config.CompiledTarget
		wantKind       agentdv1.DecisionKind
		wantAsync      uint32
	}{
		{
			name:           "sync error fail open",
			fail:           config.FailOpen,
			mode:           config.ModeSyncOnly,
			provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			invocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			raw:            func(t *testing.T) []byte { return claudeToolPre(t, "echo") },
			kinds:          []string{"tool.pre"},
			wantKind:       agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
		},
		{
			name:           "sync error fail closed tool pre",
			fail:           config.FailClosed,
			mode:           config.ModeSyncOnly,
			provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			invocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			raw:            func(t *testing.T) []byte { return claudeToolPre(t, "echo") },
			kinds:          []string{"tool.pre"},
			wantKind:       agentdv1.DecisionKind_DECISION_KIND_DENY,
		},
		{
			name:           "sync error fail closed prompt submitted",
			fail:           config.FailClosed,
			mode:           config.ModeSyncOnly,
			provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			invocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			raw:            func(t *testing.T) []byte { return claudePromptSubmitted(t) },
			kinds:          []string{"prompt.submitted"},
			wantKind:       agentdv1.DecisionKind_DECISION_KIND_BLOCK_PROMPT,
		},
		{
			name:           "sync error fail closed nonblocking event",
			fail:           config.FailClosed,
			mode:           config.ModeSyncOnly,
			provider:       agentdv1.Provider_PROVIDER_CODEX,
			invocationMode: agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY,
			raw: func(t *testing.T) []byte {
				return []byte(`{"type":"agent-turn-complete","thread_id":"t1"}`)
			},
			kinds:    []string{"notification"},
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
		},
		{
			name:           "parallel sync error keeps async dispatched",
			fail:           config.FailClosed,
			mode:           config.ModeParallel,
			provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			invocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			raw:            func(t *testing.T) []byte { return claudeToolPre(t, "echo") },
			kinds:          []string{"tool.pre"},
			async:          []config.CompiledTarget{{Kind: config.TargetBuiltin, Observe: true}},
			wantKind:       agentdv1.DecisionKind_DECISION_KIND_DENY,
			wantAsync:      1,
		},
		{
			name:           "after sync error dispatches converted outcome",
			fail:           config.FailClosed,
			mode:           config.ModeAfterSync,
			provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			invocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			raw:            func(t *testing.T) []byte { return claudeToolPre(t, "echo") },
			kinds:          []string{"tool.pre"},
			async:          []config.CompiledTarget{{Kind: config.TargetBuiltin, Observe: true}},
			wantKind:       agentdv1.DecisionKind_DECISION_KIND_DENY,
			wantAsync:      1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := dispatch.NewQueue(config.AsyncConfig{
				QueueCapacity: 8,
				WorkerLimit:   2,
				TargetTimeout: time.Second,
			}, nil)
			t.Cleanup(func() { q.Close(2 * time.Second) })
			eng := dispatch.NewEngine(q, nil, nil)
			snap := syncFailureSnap(t, tt.fail, tt.mode, tt.kinds, tt.async)

			res, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
				Provider:       tt.provider,
				RawPayload:     tt.raw(t),
				Snap:           snap,
				InvocationMode: tt.invocationMode,
				Deadline:       syncFailureDeadline(),
			})
			require.NoError(t, err)
			require.NotNil(t, res.Decision)
			assert.Equal(t, tt.wantKind, res.Decision.Kind)
			assert.Equal(t, tt.wantAsync, res.AsyncDispatchedCount)
		})
	}
}

func TestEngineRunSyncSkipWarn(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
	eng := dispatch.NewEngine(nil, log, nil)
	snap := paritySnap(t, []config.CompiledTarget{{Kind: config.TargetLog}}, nil, config.ModeSyncOnly)

	_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
		Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload: claudeToolPre(t, "echo"),
		Snap:       snap,
	})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "skip sync target")
	assert.Contains(t, buf.String(), "log")
}

func claudeKindPayload(t *testing.T, sessionID, hookEvent string) []byte {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"session_id":      sessionID,
		"cwd":             "/w",
		"hook_event_name": hookEvent,
	})
	require.NoError(t, err)
	return b
}

func TestEngineAsyncOnlySkipsSessionLock(t *testing.T) {
	t.Parallel()
	const sessionID = "sess-lock"

	tests := []struct {
		name     string
		payload  []byte
		wantFast bool
	}{
		{
			name:     "session.start async_only",
			payload:  claudeKindPayload(t, sessionID, "SessionStart"),
			wantFast: true,
		},
		{
			name:     "tool.pre waits on lock",
			payload:  claudeToolPreSession(t, sessionID, "echo hi"),
			wantFast: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := dispatch.NewQueue(config.AsyncConfig{
				QueueCapacity: 8,
				WorkerLimit:   2,
				TargetTimeout: time.Second,
			}, nil)
			t.Cleanup(func() { q.Close(2 * time.Second) })
			eng := dispatch.NewEngine(q, nil, nil)
			snap := testSnap(t)
			unlock := eng.Sessions().Lock(sessionID)
			done := make(chan error, 1)
			go func() {
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: tt.payload,
					Snap:       snap,
				})
				done <- err
			}()
			select {
			case err := <-done:
				unlock()
				require.NoError(t, err)
				assert.True(t, tt.wantFast, "Invoke(%q) returned while session lock held", tt.name)
			case <-time.After(80 * time.Millisecond):
				assert.False(t, tt.wantFast, "Invoke(%q) blocked on session lock", tt.name)
				unlock()
				require.NoError(t, <-done)
			}
		})
	}
}

func claudeToolPreSession(t *testing.T, sessionID, command string) []byte {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"session_id":      sessionID,
		"cwd":             "/w",
		"hook_event_name": "PreToolUse",
		"tool_name":       "Bash",
		"tool_use_id":     "t1",
		"tool_input":      map[string]any{"command": command},
	})
	require.NoError(t, err)
	return b
}

func TestEngineInvokeObserveKinds(t *testing.T) {
	t.Parallel()
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 8,
		WorkerLimit:   2,
		TargetTimeout: time.Second,
	}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil, nil)
	snap := testSnap(t)

	tests := []struct {
		name    string
		payload []byte
		kind    string
	}{
		{name: "session_start", payload: claudeKindPayload(t, "s", "SessionStart"), kind: "session.start"},
		{name: "subagent_start", payload: claudeKindPayload(t, "s", "SubagentStart"), kind: "subagent.start"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
				Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
				RawPayload: tt.payload,
				Snap:       snap,
			})
			require.NoError(t, err, "Invoke(%q)", tt.name)
			require.NotNil(t, res.Decision)
			assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_NO_DECISION, res.Decision.Kind)
			assert.True(t, res.Meta.HasRoute, "HasRoute %q", tt.name)
			assert.Equal(t, tt.kind, res.Meta.EventKind, "EventKind %q", tt.name)
			assert.GreaterOrEqual(t, res.AsyncDispatchedCount, uint32(1), "AsyncDispatched %q", tt.name)
		})
	}
}
