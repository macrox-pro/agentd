package server_test

import (
	"context"
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

func dialHook(t *testing.T, hook agentdv1.HookServiceServer) *grpc.ClientConn {
	t.Helper()
	s := grpc.NewServer()
	agentdv1.RegisterHookServiceServer(s, hook)
	return dialBuf(t, s)
}

func TestDaemonService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil, nil)

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
				assert.Equal(t, uint64(0), resp.GetAsyncDroppedCount(), "async_dropped_count")
				assert.Equal(t, uint32(0), resp.GetAsyncQueueDepth(), "async_queue_depth")
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

func TestDaemonServiceAsyncDropped(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)

	async := store.Current().Async
	async.QueueCapacity = 1
	async.WorkerLimit = 1
	q := dispatch.NewQueue(async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })

	block := make(chan struct{})
	started := make(chan struct{})
	require.True(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {
		close(started)
		<-block
	}}), "first job")
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not start first job")
	}
	// Fill the single buffer slot while the worker is blocked.
	require.True(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {}}), "buffered job")
	require.False(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {}}), "overflow drop")

	srv := server.New(server.Options{
		Store:     store,
		Engine:    dispatch.NewEngine(q, nil, nil),
		StartedAt: time.Now().UTC(),
		Version:   "test",
	})
	conn := dialBuf(t, srv)
	daemon := agentdv1.NewDaemonServiceClient(conn)

	resp, err := daemon.Status(ctx, &agentdv1.StatusRequest{})
	require.NoError(t, err, "Status()")
	assert.Equal(t, uint64(1), resp.GetAsyncDroppedCount(), "async_dropped_count")
	assert.Equal(t, uint32(1), resp.GetAsyncQueueDepth(), "async_queue_depth")
	close(block)
}
