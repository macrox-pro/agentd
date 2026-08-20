package hookedge_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/hookedge"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/transport"
)

func TestConformanceFixtures(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "agentd.sock")
	cfg := filepath.Join(dir, "missing.yaml")

	store, err := config.Load(context.Background(), cfg)
	require.NoError(t, err, "Load(%q)", cfg)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)

	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })

	gs := server.New(server.Options{Store: store, Engine: eng, Version: "test"})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitForSocket(t, socket)

	tests := []struct {
		name        string
		cliProvider string
		ahProvider  agenthooks.Provider
		fixture     string
		serve       bool // OpenCode uses NDJSON serve bridge
	}{
		{name: "claude pre_tool_use", cliProvider: "claude-code", ahProvider: agenthooks.ProviderClaudeCode, fixture: "claude/pre_tool_use.json"},
		{name: "codex pre_tool_use", cliProvider: "codex", ahProvider: agenthooks.ProviderCodex, fixture: "codex/pre_tool_use.json"},
		{name: "cursor pre_tool_use", cliProvider: "cursor", ahProvider: agenthooks.ProviderCursor, fixture: "cursor/pre_tool_use.json"},
		{name: "gemini before_tool", cliProvider: "gemini", ahProvider: agenthooks.ProviderGemini, fixture: "gemini/before_tool.json"},
		{name: "kimi pre_tool_use", cliProvider: "kimi-code", ahProvider: agenthooks.ProviderKimi, fixture: "kimi/pre_tool_use.json"},
		{name: "opencode tool_execute_before", cliProvider: "opencode", ahProvider: agenthooks.ProviderOpenCode, fixture: "opencode/tool_execute_before.json", serve: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := agenthookstest.Fixture(t, tt.fixture)
			var stdout, stderr bytes.Buffer
			opts := hookedge.Options{
				Socket:   socket,
				Provider: tt.cliProvider,
				Stdin:    bytes.NewReader(payload),
				Stdout:   &stdout,
				Stderr:   &stderr,
				Timeout:  5 * time.Second,
			}
			var code int
			if tt.serve {
				opts.Stdin = bytes.NewReader(append(append([]byte(nil), payload...), '\n'))
				code = hookedge.Serve(context.Background(), opts)
			} else {
				code = hookedge.Run(context.Background(), opts)
			}
			assert.Equal(t, 0, code, "exit code; stderr=%s", stderr.String())
			agenthookstest.AssertNoOp(t, tt.ahProvider, agenthookstest.Result{
				Stdout:   stdout.Bytes(),
				Stderr:   stderr.String(),
				ExitCode: code,
			})
		})
	}
}
