//go:build integration

package hookedge_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/hookedge"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/transport"
)

// Integration: full daemon↔hookedge round-trip over a real unix socket.
// Run: go test -tags=integration ./internal/hookedge/ -count=1 -run TestRoundTrip
func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")
	cfgPath := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(`version: 1
guards:
  secrets:
    enabled: true
    action: deny
`), 0o600))

	store, err := config.Load(context.Background(), cfgPath)
	require.NoError(t, err, "Load")

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)

	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen")
	t.Cleanup(func() { _ = ln.Close() })

	gs := server.New(server.Options{Store: store, Engine: eng, Version: "integration"})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitForSocket(t, socket)

	payload, err := json.Marshal(map[string]any{
		"session_id":      "int-1",
		"cwd":             dir,
		"hook_event_name": "PreToolUse",
		"tool_name":       "Bash",
		"tool_use_id":     "t1",
		"tool_input":      map[string]any{"command": "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"},
	})
	require.NoError(t, err)

	var stdout, stderr bytes.Buffer
	code := hookedge.Run(context.Background(), hookedge.Options{
		Socket:   socket,
		Provider: "claude-code",
		Stdin:    bytes.NewReader(payload),
		Stdout:   &stdout,
		Stderr:   &stderr,
		Timeout:  5 * time.Second,
	})
	assert.Equal(t, 0, code, "stderr=%s", stderr.String())
	assert.Contains(t, stdout.String(), `"permissionDecision":"deny"`)
	assert.NotContains(t, stdout.String(), "AKIAIOSFODNN7EXAMPLE")
}
