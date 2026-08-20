package hookedge_test

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookedge"
	"github.com/macrox-pro/agentd/internal/transport"
)

func TestServeOpenCodeInitialize(t *testing.T) {
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

	stdin := strings.NewReader(`{"seq":1,"hook":"initialize","input":{"serverUrl":"http://127.0.0.1:1","directory":"/work","worktree":""}}` + "\n")
	var stdout, stderr bytes.Buffer
	code := hookedge.Serve(context.Background(), hookedge.Options{
		Socket:   socket,
		Provider: "opencode",
		Stdin:    stdin,
		Stdout:   &stdout,
		Stderr:   &stderr,
		Timeout:  2 * time.Second,
	})
	assert.Equal(t, 0, code, "Serve(): %s", stderr.String())
	assert.Contains(t, stdout.String(), `"seq":1`)
}

func TestServeRejectsNonOpenCode(t *testing.T) {
	t.Parallel()
	var stderr bytes.Buffer
	code := hookedge.Serve(context.Background(), hookedge.Options{
		Provider: "claude-code",
		Stderr:   &stderr,
	})
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "opencode")
}
