package server_test

import (
	"context"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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

func TestDaemonService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)

	shutdownCh := make(chan struct{}, 1)
	srv := server.New(server.Options{
		Store:     store,
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
			name: "status generation",
			run: func(t *testing.T) {
				t.Helper()
				resp, err := daemon.Status(ctx, &agentdv1.StatusRequest{})
				require.NoError(t, err, "Status()")
				assert.GreaterOrEqual(t, resp.GetConfig().GetGeneration(), uint64(1), "Status()")
			},
		},
		{
			name: "invoke no decision",
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
