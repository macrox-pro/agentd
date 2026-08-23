package targets_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch/targets"
)

func TestFileInvokeAsync(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "audit", "out.jsonl")
	f := &targets.File{}
	err := f.InvokeAsync(context.Background(), targets.AsyncRequest{
		Provider:  "claude-code",
		EventKind: "tool.pre",
		Raw:       []byte(`{}`),
		Target:    config.CompiledTarget{Kind: config.TargetFile, Path: path},
		SyncOutcome: &targets.SyncOutcome{
			Kind:   agentdv1.DecisionKind_DECISION_KIND_ASK,
			Reason: "secret",
		},
	})
	require.NoError(t, err)
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"kind":"tool.pre"`)
	assert.Contains(t, string(b), "ASK")
}

func TestLogInvokeAsync(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	l := &targets.Log{Logger: log}
	require.NoError(t, l.InvokeAsync(context.Background(), targets.AsyncRequest{
		Provider:  "claude-code",
		EventKind: "tool.pre",
		Target:    config.CompiledTarget{Kind: config.TargetLog, Level: "info"},
	}))
	assert.Contains(t, buf.String(), "dispatch async")
}

func TestHTTPInvokeAsync(t *testing.T) {
	t.Parallel()
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)

	h := &targets.HTTP{Client: srv.Client()}
	require.NoError(t, h.InvokeAsync(context.Background(), targets.AsyncRequest{
		Provider:  "claude-code",
		EventKind: "tool.pre",
		Raw:       []byte(`{"ok":true}`),
		Target:    config.CompiledTarget{Kind: config.TargetHTTP, URL: srv.URL},
	}))
	var env map[string]any
	require.NoError(t, json.Unmarshal(gotBody, &env))
	assert.Equal(t, "claude-code", env["provider"])
	assert.Equal(t, "tool.pre", env["kind"])
}

func TestExecInvokeAsync(t *testing.T) {
	t.Parallel()
	out := filepath.Join(t.TempDir(), "stdin.txt")
	e := &targets.Exec{}
	err := e.InvokeAsync(context.Background(), targets.AsyncRequest{
		Raw: []byte("hello"),
		Target: config.CompiledTarget{
			Kind:    config.TargetExec,
			Command: []string{"sh", "-c", "cat > " + out},
			Stdin:   "raw",
		},
	})
	require.NoError(t, err)
	b, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(b))
}

func TestGRPCInvokeSync(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		peerKind agentdv1.DecisionKind
		peerErr  bool
		wantKind agenthooks.DecisionKind
	}{
		{
			name:     "peer deny",
			peerKind: agentdv1.DecisionKind_DECISION_KIND_DENY,
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:     "peer allow",
			peerKind: agentdv1.DecisionKind_DECISION_KIND_ALLOW,
			wantKind: agenthooks.DecisionAllow,
		},
		{
			name:    "peer error",
			peerErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := &targets.GRPC{
				InvokePeer: func(ctx context.Context, endpoint string, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
					assert.Equal(t, "unix:///peer.sock", endpoint)
					if tt.peerErr {
						return nil, assert.AnError
					}
					return &agentdv1.InvokeResponse{
						Decision: &agentdv1.Decision{Kind: tt.peerKind, Reason: "from-peer"},
					}, nil
				},
			}
			d, err := g.InvokeSync(context.Background(), targets.SyncRequest{
				Provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
				Raw:      []byte(`{}`),
				Target: config.CompiledTarget{
					Kind:     config.TargetGRPC,
					Endpoint: "unix:///peer.sock",
					Timeout:  time.Second,
				},
			})
			if tt.peerErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, d)
			assert.Equal(t, tt.wantKind, d.Kind())
		})
	}
}

func TestGRPCInvokeAsync(t *testing.T) {
	t.Parallel()
	called := false
	g := &targets.GRPC{
		InvokePeer: func(ctx context.Context, endpoint string, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
			called = true
			assert.Equal(t, agentdv1.Provider_PROVIDER_CLAUDE_CODE, req.GetProvider())
			return &agentdv1.InvokeResponse{
				Decision: &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_ALLOW},
			}, nil
		},
	}
	require.NoError(t, g.InvokeAsync(context.Background(), targets.AsyncRequest{
		Provider: "claude-code",
		Raw:      []byte(`{}`),
		Target:   config.CompiledTarget{Kind: config.TargetGRPC, Endpoint: "/tmp/x.sock"},
	}))
	assert.True(t, called)
}

func TestNewSyncInvoker(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		kind    config.TargetKind
		wantErr bool
	}{
		{name: "builtin", kind: config.TargetBuiltin},
		{name: "grpc", kind: config.TargetGRPC},
		{name: "log_not_sync", kind: config.TargetLog, wantErr: true},
		{name: "unknown_kind", kind: config.TargetKind("nope"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inv, err := targets.NewSyncInvoker(config.CompiledTarget{Kind: tt.kind}, &targets.Builtin{}, nil)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, inv)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, inv)
		})
	}
}

func TestNewAsyncInvoker(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		kind    config.TargetKind
		wantErr bool
	}{
		{name: "builtin", kind: config.TargetBuiltin},
		{name: "log", kind: config.TargetLog},
		{name: "file", kind: config.TargetFile},
		{name: "http", kind: config.TargetHTTP},
		{name: "exec", kind: config.TargetExec},
		{name: "grpc", kind: config.TargetGRPC},
		{name: "unknown", kind: config.TargetKind("nope"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inv, err := targets.NewAsyncInvoker(config.CompiledTarget{Kind: tt.kind}, &targets.Builtin{}, nil)
			if tt.wantErr {
				require.Error(t, err, "NewAsyncInvoker(%q)", tt.name)
				assert.Nil(t, inv)
				return
			}
			require.NoError(t, err, "NewAsyncInvoker(%q)", tt.name)
			require.NotNil(t, inv)
		})
	}
}

func TestEventKindOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		typed any
		want  string
	}{
		{name: "nil", typed: nil, want: string(agenthooks.KindOther)},
		{name: "tool pre", typed: &agenthooks.ToolPreEvent{Event: agenthooks.Event{Kind: agenthooks.KindToolPre}}, want: string(agenthooks.KindToolPre)},
		{name: "empty kind", typed: &agenthooks.ToolPreEvent{}, want: string(agenthooks.KindOther)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, targets.EventKindOf(tt.typed), "EventKindOf(%q)", tt.name)
		})
	}
}

func TestProjectRootOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		snap *config.Snapshot
		want string
	}{
		{name: "nil", snap: nil, want: ""},
		{name: "empty path", snap: &config.Snapshot{}, want: ""},
		{name: "with project file", snap: &config.Snapshot{ProjectPath: "/repo/.agentd.yaml"}, want: "/repo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, targets.ProjectRootOf(tt.snap), "ProjectRootOf(%q)", tt.name)
		})
	}
}

func TestLogInvokeAsyncLevels(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		level string
	}{
		{name: "debug", level: "debug"},
		{name: "warn", level: "warn"},
		{name: "error", level: "error"},
		{name: "default info", level: "nope"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
			l := &targets.Log{Logger: log}
			require.NoError(t, l.InvokeAsync(context.Background(), targets.AsyncRequest{
				Provider:  "claude-code",
				EventKind: "tool.pre",
				Target:    config.CompiledTarget{Kind: config.TargetLog, Level: tt.level},
			}), "Log.InvokeAsync(%q)", tt.name)
			assert.Contains(t, buf.String(), "dispatch async", "Log.InvokeAsync(%q)", tt.name)
		})
	}
}

func TestGRPCSyncFailMode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		onError  config.FailMode
		wantKind agenthooks.DecisionKind
	}{
		{name: "fail_open", onError: config.FailOpen, wantKind: agenthooks.DecisionNoDecision},
		{name: "fail_closed", onError: config.FailClosed, wantKind: agenthooks.DecisionDeny},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &targets.GRPCSync{
				Inner: &targets.GRPC{
					InvokePeer: func(ctx context.Context, endpoint string, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
						return nil, assert.AnError
					},
				},
			}
			d, err := w.InvokeSync(context.Background(), targets.SyncRequest{
				Provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
				Raw:      []byte(`{}`),
				Target: config.CompiledTarget{
					Kind:     config.TargetGRPC,
					Endpoint: "unix:///peer.sock",
					OnError:  tt.onError,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, d)
			assert.Equal(t, tt.wantKind, d.Kind())
		})
	}
}
