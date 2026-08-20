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
