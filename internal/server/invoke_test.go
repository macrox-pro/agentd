package server_test

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/trajectory"
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

func TestHookServiceInvoke(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil, nil)

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
			name: "cursor argv decodes",
			run: func(t *testing.T) {
				t.Helper()
				resp, err := hook.Invoke(ctx, &agentdv1.InvokeRequest{
					Provider:       agentdv1.Provider_PROVIDER_CURSOR,
					RawPayload:     agenthookstest.Fixture(t, "cursor/pre_tool_use.json"),
					InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
					Cwd:            "/work/repo",
				})
				require.NoError(t, err, "Invoke(cursor)")
				assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_NO_DECISION, resp.GetDecision().GetKind(), "Invoke(cursor)")
				assert.GreaterOrEqual(t, resp.GetAsyncDispatchedCount(), uint32(1), "cursor decoded and routed")
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

func TestInvokeTrajectoryTable(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		configYAML string
		provider   agentdv1.Provider
		raw        func(t *testing.T) []byte
		mode       agentdv1.InvocationMode
		cwd        string
		checkSeq   bool
		providerID string
	}{
		{
			name: "claude",
			configYAML: `version: 1
trajectory:
  enabled: true
`,
			provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			raw:        func(t *testing.T) []byte { return claudeToolPre(t, "go test") },
			mode:       agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			cwd:        "/w",
			checkSeq:   true,
			providerID: "claude-code",
		},
		{
			name: "cursor argv",
			configYAML: `version: 1
trajectory:
  enabled: true
dispatch:
  - name: observe-all
    match:
      kind: ["*"]
      provider: ["*"]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: builtin
        observe: true
`,
			provider: agentdv1.Provider_PROVIDER_CURSOR,
			raw: func(t *testing.T) []byte {
				return agenthookstest.Fixture(t, "cursor/pre_tool_use.json")
			},
			mode:       agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
			cwd:        "/work/repo",
			providerID: "cursor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			t.Setenv("XDG_STATE_HOME", stateDir)

			path := filepath.Join(t.TempDir(), "agentd.yaml")
			require.NoError(t, os.WriteFile(path, []byte(tt.configYAML), 0o600), "WriteFile(%q)", tt.name)
			store, err := config.Load(ctx, path)
			require.NoError(t, err, "Load(%q)", tt.name)

			q := dispatch.NewQueue(store.Current().Async, nil)
			t.Cleanup(func() { q.Close(2 * time.Second) })
			recorder := trajectory.NewRecorder(trajectory.DefaultSessionsDir(), store.Current().Trajectory.QueueCapacity, nil)
			t.Cleanup(func() { recorder.Close(2 * time.Second) })

			srv := server.New(server.Options{
				Store:     store,
				Engine:    dispatch.NewEngine(q, nil, nil),
				Recorder:  recorder,
				StartedAt: time.Now().UTC(),
				Version:   "test",
			})
			conn := dialBuf(t, srv)
			hook := agentdv1.NewHookServiceClient(conn)

			_, err = hook.Invoke(ctx, &agentdv1.InvokeRequest{
				Provider:       tt.provider,
				RawPayload:     tt.raw(t),
				InvocationMode: tt.mode,
				Cwd:            tt.cwd,
			})
			require.NoError(t, err, "Invoke(%q)", tt.name)

			time.Sleep(150 * time.Millisecond)
			summaries, err := trajectory.ListSessions(trajectory.DefaultSessionsDir(), tt.providerID)
			require.NoError(t, err, "ListSessions(%q)", tt.name)
			require.NotEmpty(t, summaries, "ListSessions(%q)", tt.name)

			if !tt.checkSeq {
				return
			}
			events, err := trajectory.ReadEvents(summaries[0].Path)
			require.NoError(t, err, "ReadEvents(%q)", tt.name)
			require.GreaterOrEqual(t, len(events), 3, "ReadEvents(%q)", tt.name)
			types := map[string]bool{}
			for _, e := range events {
				types[e.Type] = true
			}
			assert.True(t, types[trajectory.TypeSessionOpen], "ReadEvents(%q)", tt.name)
			assert.True(t, types[trajectory.TypeHookInvoked], "ReadEvents(%q)", tt.name)
			assert.True(t, types[trajectory.TypeHookDecided], "ReadEvents(%q)", tt.name)
			for i, e := range events {
				assert.Equal(t, uint64(i+1), e.Seq, "ReadEvents(%q)", tt.name)
			}
		})
	}
}

func TestInvokeTrajectorySyncFailure(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateDir)

	hangDir, err := os.MkdirTemp("/tmp", "agentd-hang-")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(hangDir) })
	sock := filepath.Join(hangDir, "s.sock")
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

	path := filepath.Join(t.TempDir(), "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  enabled: true
policy:
  fail: fail_closed
dispatch:
  - name: grpc-sync
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: grpc
        endpoint: unix://`+sock+`
        timeout: 500ms
        on_error: fail_open
`), 0o600))

	store, err := config.Load(ctx, path)
	require.NoError(t, err)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	recorder := trajectory.NewRecorder(trajectory.DefaultSessionsDir(), store.Current().Trajectory.QueueCapacity, nil)
	t.Cleanup(func() { recorder.Close(2 * time.Second) })

	srv := server.New(server.Options{
		Store:     store,
		Engine:    dispatch.NewEngine(q, nil, nil),
		Recorder:  recorder,
		StartedAt: time.Now().UTC(),
		Version:   "test",
	})
	conn := dialBuf(t, srv)
	hook := agentdv1.NewHookServiceClient(conn)

	invokeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = hook.Invoke(invokeCtx, &agentdv1.InvokeRequest{
		Provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload:     claudeToolPre(t, "echo"),
		InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
		Cwd:            "/w",
		Deadline:       timestamppb.New(time.Now().Add(50 * time.Millisecond)),
	})
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)
	summaries, err := trajectory.ListSessions(trajectory.DefaultSessionsDir(), "claude-code")
	require.NoError(t, err)
	require.NotEmpty(t, summaries)

	events, err := trajectory.ReadEvents(summaries[0].Path)
	require.NoError(t, err)
	var decided bool
	for _, e := range events {
		if e.Type == trajectory.TypeHookDecided {
			decided = true
			break
		}
	}
	assert.True(t, decided, "trajectory should record converted sync failure decision")
}
