package hookedge_test

import (
	"bytes"
	"context"
	"os"
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
		ConfigPath: filepath.Join(t.TempDir(), "missing.yaml"),
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
		ConfigPath: filepath.Join(dir, "missing.yaml"),
		Provider:   "codex",
		PayloadArg: `{"type":"agent-turn-complete","thread_id":"t1"}`,
		Stderr:     &stderr,
	})
	assert.Equal(t, 0, code, "Notify(): %s", stderr.String())
}

func TestNotifyOffline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(t *testing.T) (configPath, socket string)
		wantCode int
	}{
		{
			name: "daemon down default fail_open",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "missing.yaml"), filepath.Join(dir, "missing.sock")
			},
			wantCode: 0,
		},
		{
			name: "daemon down fail_closed",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				cfg := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(cfg, []byte("version: 1\npolicy:\n  offline: fail_closed\n"), 0o600))
				return cfg, filepath.Join(dir, "missing.sock")
			},
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, socket := tt.setup(t)
			var stderr bytes.Buffer
			code := hookedge.Notify(context.Background(), hookedge.Options{
				Socket:     socket,
				ConfigPath: cfg,
				Provider:   "codex",
				PayloadArg: `{"type":"agent-turn-complete","thread_id":"t1"}`,
				Stderr:     &stderr,
			})
			assert.Equal(t, tt.wantCode, code, "Notify(%q): %s", tt.name, stderr.String())
			assert.Contains(t, stderr.String(), "daemon not running", "Notify(%q)", tt.name)
		})
	}
}
