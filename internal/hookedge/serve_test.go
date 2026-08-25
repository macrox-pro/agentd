package hookedge_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		Socket:     socket,
		ConfigPath: filepath.Join(dir, "missing.yaml"),
		Provider:   "opencode",
		Stdin:      stdin,
		Stdout:     &stdout,
		Stderr:     &stderr,
		Timeout:    2 * time.Second,
	})
	assert.Equal(t, 0, code, "Serve(): %s", stderr.String())
	assert.Contains(t, stdout.String(), `"seq":1`)
}

func TestServeRejectsNonOpenCode(t *testing.T) {
	t.Parallel()
	var stderr bytes.Buffer
	code := hookedge.Serve(context.Background(), hookedge.Options{
		Provider:   "claude-code",
		ConfigPath: filepath.Join(t.TempDir(), "missing.yaml"),
		Stderr:     &stderr,
	})
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "opencode")
}

func TestServeOffline(t *testing.T) {
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

	initLine := `{"seq":1,"hook":"initialize","input":{"serverUrl":"http://127.0.0.1:1","directory":"/work","worktree":""}}` + "\n"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, socket := tt.setup(t)
			var stdout, stderr bytes.Buffer
			code := hookedge.Serve(context.Background(), hookedge.Options{
				Socket:     socket,
				ConfigPath: cfg,
				Provider:   "opencode",
				Stdin:      strings.NewReader(initLine),
				Stdout:     &stdout,
				Stderr:     &stderr,
				Timeout:    2 * time.Second,
			})
			assert.Equal(t, tt.wantCode, code, "Serve(%q): %s", tt.name, stderr.String())
			assert.Contains(t, stderr.String(), "daemon not running", "Serve(%q)", tt.name)
		})
	}
}

func TestServeOfflineInvokeCache(t *testing.T) {
	// Corner case name (PROGRESS / intent): invoke fail caches offline once
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")

	ln, err := transport.Listen(socket)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	gs := grpc.NewServer()
	agentdv1.RegisterHookServiceServer(gs, failInvokeHook{})
	agentdv1.RegisterDaemonServiceServer(gs, okDaemon{})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitForSocket(t, socket)

	// tool.execute.before frames carry Raw into middleware; initialize often skips Invoke.
	stdin := strings.NewReader(
		`{"seq":1,"hook":"tool.execute.before","input":{"sessionID":"s","callID":"c1","tool":"bash"},"output":{"args":{}}}` + "\n" +
			`{"seq":2,"hook":"tool.execute.before","input":{"sessionID":"s","callID":"c2","tool":"bash"},"output":{"args":{}}}` + "\n",
	)
	var stdout, stderr bytes.Buffer
	code := hookedge.Serve(context.Background(), hookedge.Options{
		Socket:     socket,
		ConfigPath: filepath.Join(dir, "missing.yaml"),
		Provider:   "opencode",
		Stdin:      stdin,
		Stdout:     &stdout,
		Stderr:     &stderr,
		Timeout:    2 * time.Second,
	})
	assert.Equal(t, 0, code, "Serve(invoke cache): %s", stderr.String())
	assert.Equal(t, 1, strings.Count(stderr.String(), "daemon not running"),
		"stderr should print daemon not running once, got %q", stderr.String())
}

type failInvokeHook struct {
	agentdv1.UnimplementedHookServiceServer
}

func (failInvokeHook) Invoke(context.Context, *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	return nil, status.Error(codes.Unavailable, "daemon invoke unavailable")
}
