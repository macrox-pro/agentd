package hookedge_test

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/hookedge"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/transport"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")
	cfg := filepath.Join(dir, "missing.yaml")

	store, err := config.Load(context.Background(), cfg)
	require.NoError(t, err, "Load(%q)", cfg)

	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })

	gs := server.New(server.Options{Store: store, Version: "test"})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)

	waitForSocket(t, socket)

	tests := []struct {
		name     string
		socket   string
		provider string
		payload  string
		stdin    io.Reader
		wantOut  string
		wantCode int
		wantErr  string
	}{
		{
			name:     "claude no decision",
			socket:   socket,
			provider: "claude-code",
			payload:  `{}`,
			wantOut:  `{}`,
			wantCode: 0,
		},
		{
			name:     "codex no decision",
			socket:   socket,
			provider: "codex",
			payload:  `{}`,
			wantOut:  ``,
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			stdin := tt.stdin
			if stdin == nil {
				stdin = bytes.NewReader([]byte(tt.payload))
			}
			code := hookedge.Run(context.Background(), hookedge.Options{
				Socket:   tt.socket,
				Provider: tt.provider,
				Stdin:    stdin,
				Stdout:   &stdout,
				Stderr:   &stderr,
			})
			assert.Equal(t, tt.wantCode, code, "Run(%q)", tt.name)
			assert.Equal(t, tt.wantOut, stdout.String(), "Run(%q)", tt.name)
			if tt.wantErr != "" {
				assert.Contains(t, stderr.String(), tt.wantErr, "Run(%q)", tt.name)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		missingSocket bool
		provider      string
		stdin         io.Reader
		wantCode      int
		wantOut       string
		wantErrSubstr string
	}{
		{
			name:          "daemon down",
			missingSocket: true,
			provider:      "claude-code",
			stdin:         bytes.NewReader([]byte(`{}`)),
			wantCode:      1,
			wantErrSubstr: "daemon not running",
		},
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
			socket := ""
			if tt.missingSocket {
				socket = filepath.Join(t.TempDir(), "missing.sock")
			}

			var stdout, stderr bytes.Buffer
			code := hookedge.Run(context.Background(), hookedge.Options{
				Socket:   socket,
				Provider: tt.provider,
				Stdin:    tt.stdin,
				Stdout:   &stdout,
				Stderr:   &stderr,
			})
			assert.Equal(t, tt.wantCode, code, "Run(%q)", tt.name)
			assert.Equal(t, tt.wantOut, stdout.String(), "Run(%q)", tt.name)
			if tt.wantErrSubstr != "" {
				assert.Contains(t, stderr.String(), tt.wantErrSubstr, "Run(%q)", tt.name)
			}
		})
	}
}

func TestRunUnsupportedDecision(t *testing.T) {
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
		Stdin:    bytes.NewReader([]byte(`{}`)),
		Stdout:   &stdout,
		Stderr:   &stderr,
	})
	assert.Equal(t, 1, code, "Run(unsupported)")
	assert.Empty(t, stdout.String(), "Run(unsupported)")
	assert.Contains(t, stderr.String(), "unsupported decision", "Run(unsupported)")
}

type denyHook struct {
	agentdv1.UnimplementedHookServiceServer
}

func (denyHook) Invoke(context.Context, *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	return &agentdv1.InvokeResponse{
		Decision: &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_DENY},
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
