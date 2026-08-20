package hookclient_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookclient"
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

func TestDialHealthInvokeClose(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")

	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })

	gs := grpc.NewServer()
	svc := healthOnly{}
	agentdv1.RegisterDaemonServiceServer(gs, svc)
	agentdv1.RegisterHookServiceServer(gs, svc)
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitForSocket(t, socket)

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
}

func TestHealthMissingSocket(t *testing.T) {
	t.Parallel()
	cli, err := hookclient.Dial(context.Background(), filepath.Join(t.TempDir(), "no-such.sock"))
	require.NoError(t, err, "Dial is lazy")
	t.Cleanup(func() { _ = cli.Close() })

	_, err = cli.Health(context.Background())
	require.Error(t, err, "Health(missing)")
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
