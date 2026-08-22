package cmd_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/cmd"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

type execOpts struct {
	args       []string
	configPath string
	socketPath string
}

type execResult struct {
	out string
	err error
}

func resetFlag(f *pflag.Flag) {
	switch f.Value.Type() {
	case "stringSlice", "stringArray":
		// DefValue for nil slice defaults is "[]", which Set parses as one element "[]".
		_ = f.Value.Set("")
	default:
		_ = f.Value.Set(f.DefValue)
	}
	f.Changed = false
}

func resetCommandFlags(c *cobra.Command) {
	if c == nil {
		return
	}
	c.PersistentFlags().VisitAll(resetFlag)
	c.Flags().VisitAll(resetFlag)
	for _, sub := range c.Commands() {
		resetCommandFlags(sub)
	}
}

func tempSessionsEnv(t *testing.T) string {
	t.Helper()
	state := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(state, "agentd", "sessions"), 0o700))
	t.Setenv("XDG_STATE_HOME", state)
	return filepath.Join(state, "agentd", "sessions")
}

func executeRoot(t *testing.T, opts execOpts) execResult {
	t.Helper()
	root := cmd.RootCommand()
	resetCommandFlags(root)
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)

	args := make([]string, 0, len(opts.args)+4)
	if opts.configPath != "" {
		args = append(args, "--config", opts.configPath)
	}
	if opts.socketPath != "" {
		args = append(args, "--socket", opts.socketPath)
	}
	args = append(args, opts.args...)
	root.SetArgs(args)
	err := root.Execute()
	return execResult{out: buf.String(), err: err}
}

func writeSessionLedger(t *testing.T, root, provider, sessionID string, n int) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	path := filepath.Join(dir, sessionID+".jsonl")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()
	enc := json.NewEncoder(f)
	now := time.Now().UTC()
	for i := 1; i <= n; i++ {
		require.NoError(t, enc.Encode(trajectory.Event{
			Seq:       uint64(i),
			Type:      trajectory.TypeHookInvoked,
			Source:    trajectory.SourceHook,
			TS:        now,
			Provider:  provider,
			SessionID: sessionID,
			Data:      json.RawMessage(`{"kind":"tool.pre"}`),
		}))
	}
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

func writeSessionLedgerNoRaw(t *testing.T, root, provider, sessionID string) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	line := `{"seq":1,"type":"hook/invoked","source":"hook","provider":"` + provider +
		`","session_id":"` + sessionID + `","data":{"kind":"tool.pre"}}` + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, sessionID+".jsonl"), []byte(line), 0o600))
}

func claudeTranscriptFixture(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "internal", "trajectory", "importer", "testdata", "claude_session.jsonl")
}
