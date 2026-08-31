package hookedge_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/hookedge"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/transport"
)

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

func TestRun(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")
	cfg := filepath.Join(dir, "missing.yaml")

	store, err := config.Load(context.Background(), cfg)
	require.NoError(t, err, "Load(%q)", cfg)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil, nil)

	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })

	gs := server.New(server.Options{Store: store, Engine: eng, Version: "test"})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)

	waitForSocket(t, socket)

	tests := []struct {
		name       string
		provider   string
		payload    []byte
		wantCode   int
		wantSubstr string
		wantExact  string
		denySecret bool
	}{
		{
			name:      "claude clean no decision",
			provider:  "claude-code",
			payload:   claudeToolPre(t, "go test ./..."),
			wantExact: `{}`,
			wantCode:  0,
		},
		{
			name:       "claude secret ask",
			provider:   "claude-code",
			payload:    claudeToolPre(t, "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"),
			wantSubstr: `"permissionDecision":"ask"`,
			wantCode:   0,
			denySecret: true,
		},
		{
			name:      "undecodable no-op",
			provider:  "claude-code",
			payload:   []byte(`{}`),
			wantExact: `{}`,
			wantCode:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := hookedge.Run(context.Background(), hookedge.Options{
				Socket:   socket,
				Provider: tt.provider,
				Stdin:    bytes.NewReader(tt.payload),
				Stdout:   &stdout,
				Stderr:   &stderr,
			})
			assert.Equal(t, tt.wantCode, code, "Run(%q)", tt.name)
			out := stdout.String()
			if tt.wantExact != "" {
				assert.Equal(t, tt.wantExact, out, "Run(%q)", tt.name)
			}
			if tt.wantSubstr != "" {
				assert.Contains(t, out, tt.wantSubstr, "Run(%q)", tt.name)
			}
			if tt.denySecret {
				assert.NotContains(t, out, "AKIAIOSFODNN7EXAMPLE", "Run(%q)", tt.name)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		provider      string
		stdin         io.Reader
		wantCode      int
		wantOut       string
		wantErrSubstr string
	}{
		{
			name:          "empty stdin",
			provider:      "claude-code",
			stdin:         bytes.NewReader(nil),
			wantCode:      1,
			wantErrSubstr: "empty stdin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			code := hookedge.Run(context.Background(), hookedge.Options{
				Provider:   tt.provider,
				ConfigPath: filepath.Join(t.TempDir(), "missing.yaml"),
				Stdin:      tt.stdin,
				Stdout:     &stdout,
				Stderr:     &stderr,
			})
			assert.Equal(t, tt.wantCode, code, "Run(%q)", tt.name)
			assert.Equal(t, tt.wantOut, stdout.String(), "Run(%q)", tt.name)
			if tt.wantErrSubstr != "" {
				assert.Contains(t, stderr.String(), tt.wantErrSubstr, "Run(%q)", tt.name)
			}
		})
	}
}

func TestRunOffline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(t *testing.T) (configPath, socket string, stdin io.Reader)
		wantCode int
	}{
		{
			name: "daemon down default fail_open",
			setup: func(t *testing.T) (string, string, io.Reader) {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "missing.yaml"), filepath.Join(dir, "missing.sock"), bytes.NewReader([]byte(`{}`))
			},
			wantCode: 0,
		},
		{
			name: "daemon down fail_closed",
			setup: func(t *testing.T) (string, string, io.Reader) {
				t.Helper()
				dir := t.TempDir()
				cfg := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(cfg, []byte("version: 1\npolicy:\n  offline: fail_closed\n"), 0o600))
				return cfg, filepath.Join(dir, "missing.sock"), bytes.NewReader([]byte(`{}`))
			},
			wantCode: 1,
		},
		{
			name: "daemon down project fail_closed",
			setup: func(t *testing.T) (string, string, io.Reader) {
				t.Helper()
				dir := t.TempDir()
				proj := filepath.Join(dir, "proj")
				require.NoError(t, os.MkdirAll(proj, 0o700))
				require.NoError(t, os.WriteFile(filepath.Join(proj, ".agentd.yaml"), []byte("version: 1\npolicy:\n  offline: fail_closed\n"), 0o600))
				payload, err := json.Marshal(map[string]any{"cwd": proj})
				require.NoError(t, err)
				return filepath.Join(dir, "missing.yaml"), filepath.Join(dir, "missing.sock"), bytes.NewReader(payload)
			},
			wantCode: 1,
		},
		{
			name: "daemon down invalid config",
			setup: func(t *testing.T) (string, string, io.Reader) {
				t.Helper()
				dir := t.TempDir()
				cfg := filepath.Join(dir, "bad.yaml")
				require.NoError(t, os.WriteFile(cfg, []byte(":\n  - invalid\n"), 0o600))
				return cfg, filepath.Join(dir, "missing.sock"), bytes.NewReader([]byte(`{}`))
			},
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, socket, stdin := tt.setup(t)
			var stdout, stderr bytes.Buffer
			code := hookedge.Run(context.Background(), hookedge.Options{
				Socket:     socket,
				ConfigPath: cfg,
				Provider:   "claude-code",
				Stdin:      stdin,
				Stdout:     &stdout,
				Stderr:     &stderr,
			})
			assert.Equal(t, tt.wantCode, code, "Run(%q): %s", tt.name, stderr.String())
			assert.Contains(t, stderr.String(), "daemon not running", "Run(%q)", tt.name)
		})
	}
}

func TestRunDenyDecision(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")

	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })

	gs := grpc.NewServer()
	agentdv1.RegisterHookServiceServer(gs, denyHook{})
	agentdv1.RegisterDaemonServiceServer(gs, okDaemon{})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitForSocket(t, socket)

	var stdout, stderr bytes.Buffer
	code := hookedge.Run(context.Background(), hookedge.Options{
		Socket:   socket,
		Provider: "claude-code",
		Stdin:    bytes.NewReader(claudeToolPre(t, "echo")),
		Stdout:   &stdout,
		Stderr:   &stderr,
	})
	assert.Equal(t, 0, code, "Run(deny)")
	assert.Contains(t, stdout.String(), `"permissionDecision":"deny"`, "Run(deny)")
	assert.True(t, strings.Contains(stdout.String(), "blocked") || strings.Contains(stdout.String(), "deny"), "Run(deny) reason")
}

type denyHook struct {
	agentdv1.UnimplementedHookServiceServer
}

func (denyHook) Invoke(context.Context, *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	return &agentdv1.InvokeResponse{
		Decision: &agentdv1.Decision{
			Kind:   agentdv1.DecisionKind_DECISION_KIND_DENY,
			Reason: "blocked by test",
		},
	}, nil
}

type okDaemon struct {
	agentdv1.UnimplementedDaemonServiceServer
}

func (okDaemon) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
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
