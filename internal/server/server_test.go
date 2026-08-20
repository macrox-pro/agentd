package server_test

import (
	"context"
	"encoding/json"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
)

func dialBuf(t *testing.T, srv *grpc.Server) *grpc.ClientConn {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(func() {
		srv.Stop()
		_ = lis.Close()
	})
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "dial bufconn")
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

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

func TestDaemonService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)

	shutdownCh := make(chan struct{}, 1)
	srv := server.New(server.Options{
		Store:     store,
		Engine:    eng,
		StartedAt: time.Now().UTC(),
		Version:   "test",
		OnShutdown: func() {
			select {
			case shutdownCh <- struct{}{}:
			default:
			}
		},
	})
	conn := dialBuf(t, srv)
	daemon := agentdv1.NewDaemonServiceClient(conn)
	hook := agentdv1.NewHookServiceClient(conn)

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "health ok",
			run: func(t *testing.T) {
				t.Helper()
				resp, err := daemon.Health(ctx, &agentdv1.HealthRequest{})
				require.NoError(t, err, "Health()")
				assert.Equal(t, "ok", resp.GetStatus(), "Health()")
			},
		},
		{
			name: "status generation and routes",
			run: func(t *testing.T) {
				t.Helper()
				resp, err := daemon.Status(ctx, &agentdv1.StatusRequest{})
				require.NoError(t, err, "Status()")
				assert.GreaterOrEqual(t, resp.GetConfig().GetGeneration(), uint64(1), "Status()")
				assert.Greater(t, resp.GetCompiledRouteCount(), uint32(0), "compiled_route_count")
			},
		},
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
		{
			name: "reload bumps generation",
			run: func(t *testing.T) {
				t.Helper()
				before := store.Current().Generation
				resp, err := daemon.ReloadConfig(ctx, &agentdv1.ReloadConfigRequest{})
				require.NoError(t, err, "ReloadConfig()")
				assert.Greater(t, resp.GetConfig().GetGeneration(), before, "ReloadConfig()")
			},
		},
		{
			name: "shutdown callback",
			run: func(t *testing.T) {
				t.Helper()
				_, err := daemon.Shutdown(ctx, &agentdv1.ShutdownRequest{})
				require.NoError(t, err, "Shutdown()")
				select {
				case <-shutdownCh:
				case <-time.After(2 * time.Second):
					t.Fatal("OnShutdown not called")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
