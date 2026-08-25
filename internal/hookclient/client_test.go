package hookclient_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/transport"
)

type healthOnly struct {
	agentdv1.UnimplementedDaemonServiceServer
	agentdv1.UnimplementedHookServiceServer
}

func (healthOnly) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (healthOnly) Invoke(_ context.Context, _ *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	return &agentdv1.InvokeResponse{
		Decision: &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION},
	}, nil
}

type daemonStub struct {
	agentdv1.UnimplementedDaemonServiceServer
	generation uint64
}

func (d *daemonStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (d *daemonStub) Status(context.Context, *agentdv1.StatusRequest) (*agentdv1.StatusResponse, error) {
	gen := d.generation
	if gen == 0 {
		gen = 1
	}
	return &agentdv1.StatusResponse{
		Version:   "test",
		StartedAt: timestamppb.Now(),
		Config:    &agentdv1.ConfigGeneration{Generation: gen, Fingerprint: "fp"},
	}, nil
}

func (d *daemonStub) ReloadConfig(context.Context, *agentdv1.ReloadConfigRequest) (*agentdv1.ReloadConfigResponse, error) {
	d.generation++
	if d.generation == 0 {
		d.generation = 2
	}
	return &agentdv1.ReloadConfigResponse{
		Config: &agentdv1.ConfigGeneration{Generation: d.generation, Fingerprint: "fp2"},
	}, nil
}

func (d *daemonStub) Shutdown(context.Context, *agentdv1.ShutdownRequest) (*agentdv1.ShutdownResponse, error) {
	return &agentdv1.ShutdownResponse{}, nil
}

type sessionSendStub struct {
	agentdv1.UnimplementedSessionServiceServer
	event *agentdv1.SessionEvent
}

func (s sessionSendStub) Subscribe(_ *agentdv1.SubscribeRequest, stream agentdv1.SessionService_SubscribeServer) error {
	go func() {
		time.Sleep(20 * time.Millisecond)
		_ = stream.Send(&agentdv1.SubscribeResponse{Event: s.event})
	}()
	<-stream.Context().Done()
	return stream.Context().Err()
}

type sessionBlockStub struct {
	agentdv1.UnimplementedSessionServiceServer
}

func (sessionBlockStub) Subscribe(_ *agentdv1.SubscribeRequest, stream agentdv1.SessionService_SubscribeServer) error {
	<-stream.Context().Done()
	return stream.Context().Err()
}

func startStubServer(t *testing.T, register func(*grpc.Server)) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "agentd-hookclient-")
	require.NoError(t, err, "MkdirTemp")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	socket := filepath.Join(dir, "s.sock")

	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })

	gs := grpc.NewServer()
	register(gs)
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitForSocket(t, socket)
	return socket
}

func TestHookclientDial(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "lazy dial missing socket",
			run: func(t *testing.T) {
				t.Parallel()
				cli, err := hookclient.Dial(context.Background(), filepath.Join(t.TempDir(), "no-such.sock"))
				require.NoError(t, err, "Dial is lazy")
				t.Cleanup(func() { _ = cli.Close() })
				_, err = cli.Health(context.Background())
				require.Error(t, err, "Health(missing)")
			},
		},
		{
			name: "dial ready missing socket",
			run: func(t *testing.T) {
				t.Parallel()
				cli, err := hookclient.DialReady(context.Background(), filepath.Join(t.TempDir(), "no-such.sock"))
				require.Error(t, err, "DialReady(missing)")
				assert.Nil(t, cli, "DialReady(missing)")
			},
		},
		{
			name: "dial ready healthy",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					svc := healthOnly{}
					agentdv1.RegisterDaemonServiceServer(gs, svc)
					agentdv1.RegisterHookServiceServer(gs, svc)
				})
				cli, err := hookclient.DialReady(context.Background(), socket)
				require.NoError(t, err, "DialReady(%q)", socket)
				t.Cleanup(func() { _ = cli.Close() })
			},
		},
		{
			name: "canceled ctx health",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					svc := healthOnly{}
					agentdv1.RegisterDaemonServiceServer(gs, svc)
					agentdv1.RegisterHookServiceServer(gs, svc)
				})
				cli, err := hookclient.Dial(context.Background(), socket)
				require.NoError(t, err, "Dial(%q)", socket)
				t.Cleanup(func() { _ = cli.Close() })

				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err = cli.Health(ctx)
				require.Error(t, err, "Health(canceled)")
			},
		},
		{
			name: "health invoke close",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					svc := healthOnly{}
					agentdv1.RegisterDaemonServiceServer(gs, svc)
					agentdv1.RegisterHookServiceServer(gs, svc)
				})
				cli, err := hookclient.Dial(context.Background(), socket)
				require.NoError(t, err, "Dial(%q)", socket)
				t.Cleanup(func() { _ = cli.Close() })

				health, err := cli.Health(context.Background())
				require.NoError(t, err, "Health()")
				assert.Equal(t, "ok", health.GetStatus(), "Health()")

				resp, err := cli.Invoke(context.Background(), &agentdv1.InvokeRequest{
					Provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
				})
				require.NoError(t, err, "Invoke()")
				assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_NO_DECISION, resp.GetDecision().GetKind(), "Invoke()")

				require.NoError(t, cli.Close(), "Close()")
			},
		},
		{
			name: "empty socket uses default path",
			run: func(t *testing.T) {
				cli, err := hookclient.Dial(context.Background(), "")
				require.NoError(t, err, "Dial(empty)")
				require.NotNil(t, cli, "Dial(empty)")
				require.NoError(t, cli.Close(), "Close()")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}

func TestHookclientDaemonRPC(t *testing.T) {
	stub := &daemonStub{}
	socket := startStubServer(t, func(gs *grpc.Server) {
		agentdv1.RegisterDaemonServiceServer(gs, stub)
	})

	tests := []struct {
		name string
		run  func(t *testing.T, cli *hookclient.Client)
	}{
		{
			name: "status ok",
			run: func(t *testing.T, cli *hookclient.Client) {
				resp, err := cli.Status(context.Background())
				require.NoError(t, err, "Status()")
				assert.GreaterOrEqual(t, resp.GetConfig().GetGeneration(), uint64(1), "Generation")
			},
		},
		{
			name: "reload ok",
			run: func(t *testing.T, cli *hookclient.Client) {
				resp, err := cli.Reload(context.Background())
				require.NoError(t, err, "Reload()")
				assert.NotNil(t, resp.GetConfig(), "Config")
				assert.GreaterOrEqual(t, resp.GetConfig().GetGeneration(), uint64(1), "Generation")
			},
		},
		{
			name: "shutdown ok",
			run: func(t *testing.T, cli *hookclient.Client) {
				require.NoError(t, cli.Shutdown(context.Background()), "Shutdown()")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := hookclient.Dial(context.Background(), socket)
			require.NoError(t, err, "Dial(%q)", socket)
			t.Cleanup(func() { _ = cli.Close() })
			tt.run(t, cli)
		})
	}
}

func TestSubscribeClient(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "recv one event",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterSessionServiceServer(gs, sessionSendStub{
						event: &agentdv1.SessionEvent{
							Type:     trajectory.TypeHookInvoked,
							Provider: "claude-code",
						},
					})
				})
				cli, err := hookclient.Dial(context.Background(), socket)
				require.NoError(t, err, "Dial(%q)", socket)
				t.Cleanup(func() { _ = cli.Close() })

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				stream, err := cli.Subscribe(ctx, &agentdv1.SubscribeRequest{})
				require.NoError(t, err, "Subscribe()")

				msg, err := stream.Recv()
				require.NoError(t, err, "Recv()")
				assert.Equal(t, trajectory.TypeHookInvoked, msg.GetEvent().GetType(), "Type")
				assert.Equal(t, "claude-code", msg.GetEvent().GetProvider(), "Provider")
			},
		},
		{
			name: "cancel propagation",
			run: func(t *testing.T) {
				socket := startStubServer(t, func(gs *grpc.Server) {
					agentdv1.RegisterSessionServiceServer(gs, sessionBlockStub{})
				})
				cli, err := hookclient.Dial(context.Background(), socket)
				require.NoError(t, err, "Dial(%q)", socket)
				t.Cleanup(func() { _ = cli.Close() })

				ctx, cancel := context.WithCancel(context.Background())
				stream, err := cli.Subscribe(ctx, &agentdv1.SubscribeRequest{})
				require.NoError(t, err, "Subscribe()")
				cancel()

				_, err = stream.Recv()
				require.Error(t, err, "Recv(canceled)")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}

func TestCloseNil(t *testing.T) {
	t.Parallel()
	var cli *hookclient.Client
	require.NoError(t, cli.Close(), "Close(nil)")
}

func waitForSocket(t *testing.T, socket string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		if time.Now().After(deadline) {
			t.Fatal("socket not ready")
		}
		c, err := transport.Dial(context.Background(), socket)
		if err == nil {
			_ = c.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
