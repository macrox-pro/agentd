package trajectory_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestReplayPolicy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		prov    string
		mode    string
		fixture string
		proto   agentdv1.Provider
	}{
		{name: "claude stdin", prov: "claude-code", mode: "stdin", fixture: "claude/pre_tool_use.json", proto: agentdv1.Provider_PROVIDER_CLAUDE_CODE},
		{name: "cursor argv", prov: "cursor", mode: "argv", fixture: "cursor/pre_tool_use.json", proto: agentdv1.Provider_PROVIDER_CURSOR},
		{name: "codex stdin", prov: "codex", mode: "stdin", fixture: "codex/pre_tool_use.json", proto: agentdv1.Provider_PROVIDER_CODEX},
		{name: "gemini stdin", prov: "gemini", mode: "stdin", fixture: "gemini/before_tool.json", proto: agentdv1.Provider_PROVIDER_GEMINI},
		{name: "kimi stdin", prov: "kimi-code", mode: "stdin", fixture: "kimi/pre_tool_use.json", proto: agentdv1.Provider_PROVIDER_KIMI_CODE},
		{name: "opencode stdin", prov: "opencode", mode: "stdin", fixture: "opencode/tool_execute_before.json", proto: agentdv1.Provider_PROVIDER_OPENCODE},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			raw := agenthookstest.Fixture(t, tt.fixture)
			sid := "replay-" + tt.prov
			writeReplayLedger(t, root, tt.prov, sid, tt.mode, raw)

			res, err := config.CompileMerged(nil, nil, nil)
			require.NoError(t, err)
			snap := &config.Snapshot{
				Generation: 1,
				Async:      res.Async,
				Guards:     res.Guards,
				Routes:     res.Routes,
				Policy:     config.Policy{Fail: config.FailClosed, AskFallback: config.AskFallbackDeny},
			}
			q := dispatch.NewQueue(res.Async, nil)
			t.Cleanup(func() { q.Close(2 * time.Second) })
			eng := dispatch.NewEngine(q, nil, nil)

			result, err := trajectory.ReplayPolicy(context.Background(), trajectory.ReplayOptions{
				SessionsRoot: root,
				Provider:     tt.prov,
				SessionID:    sid,
				Snap:         snap,
				Engine:       eng,
			})
			require.NoError(t, err, "ReplayPolicy(%s)", tt.prov)
			require.NotEmpty(t, result.Hits)
			assert.Empty(t, result.Hits[0].Error)
			assert.NotEmpty(t, result.Hits[0].ReplayDecision)
		})
	}
}

func TestReplayMissingRaw(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	prov, sid := "claude-code", "no-raw"
	dir := filepath.Join(root, prov)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	path := filepath.Join(dir, sid+".jsonl")
	line := `{"seq":1,"type":"hook/invoked","source":"hook","provider":"claude-code","session_id":"no-raw","data":{"kind":"tool.pre"}}` + "\n"
	require.NoError(t, os.WriteFile(path, []byte(line), 0o600))

	res, err := config.CompileMerged(nil, nil, nil)
	require.NoError(t, err)
	snap := &config.Snapshot{Generation: 1, Async: res.Async, Guards: res.Guards, Routes: res.Routes}
	q := dispatch.NewQueue(res.Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil, nil)

	_, err = trajectory.ReplayPolicy(context.Background(), trajectory.ReplayOptions{
		SessionsRoot: root,
		Provider:     prov,
		SessionID:    sid,
		Snap:         snap,
		Engine:       eng,
	})
	require.ErrorIs(t, err, trajectory.ErrReplayNoRaw)
}

func writeReplayLedger(t *testing.T, root, provider, sessionID, mode string, raw []byte) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	invData, err := json.Marshal(trajectory.HookInvokedData{Kind: "tool.pre", HasRoute: true})
	require.NoError(t, err)
	decData, err := json.Marshal(trajectory.HookDecidedData{
		Kind:     "tool.pre",
		Decision: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION.String(),
	})
	require.NoError(t, err)
	// Embed raw as JSON object/array from fixture bytes.
	events := []map[string]any{
		{
			"seq":             1,
			"type":            trajectory.TypeHookInvoked,
			"source":          trajectory.SourceHook,
			"provider":        provider,
			"session_id":      sessionID,
			"invocation_mode": mode,
			"data":            json.RawMessage(invData),
			"raw":             json.RawMessage(raw),
		},
		{
			"seq":        2,
			"type":       trajectory.TypeHookDecided,
			"source":     trajectory.SourceDecision,
			"provider":   provider,
			"session_id": sessionID,
			"data":       json.RawMessage(decData),
		},
	}
	f, err := os.Create(filepath.Join(dir, sessionID+".jsonl"))
	require.NoError(t, err)
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range events {
		require.NoError(t, enc.Encode(e))
	}
}
