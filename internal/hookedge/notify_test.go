package hookedge_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookedge"
	"github.com/macrox-pro/agentd/internal/transport"
)

func TestNotifyRejectsNonCodex(t *testing.T) {
	t.Parallel()
	var stderr bytes.Buffer
	code := hookedge.Notify(context.Background(), hookedge.Options{
		Provider:   "claude-code",
		PayloadArg: `{"type":"agent-turn-complete"}`,
		Stderr:     &stderr,
	})
	assert.Equal(t, 1, code, "Notify(non-codex): %s", stderr.String())
	assert.Contains(t, stderr.String(), "codex")
}

func TestNotify(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")

	ln, err := transport.Listen(socket)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	gs := grpc.NewServer()
	agentdv1.RegisterHookServiceServer(gs, denyHook{})
	agentdv1.RegisterDaemonServiceServer(gs, okDaemon{})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitForSocket(t, socket)

	var stderr bytes.Buffer
	code := hookedge.Notify(context.Background(), hookedge.Options{
		Socket:     socket,
		Provider:   "codex",
		PayloadArg: `{"type":"agent-turn-complete","thread_id":"t1"}`,
		Stderr:     &stderr,
	})
	assert.Equal(t, 0, code, "Notify(): %s", stderr.String())
}
