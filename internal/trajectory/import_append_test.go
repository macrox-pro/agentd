package trajectory_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestAppendImportedPreservesHookSeq(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	key := trajectory.SessionKey{Provider: provider.ClaudeCode, SessionID: "s1"}
	path := trajectory.SessionFilePath(root, key)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	hookLine := `{"seq":1,"type":"hook/invoked","source":"hook","session_id":"s1","provider":"claude-code","data":{"tool_use_id":"toolu_01"}}`
	require.NoError(t, os.WriteFile(path, []byte(hookLine+"\n"), 0o600))

	imported := []trajectory.Event{{
		Type:   trajectory.TypeTranscriptMessage,
		Source: trajectory.SourceTranscript,
		TS:     time.Now().UTC(),
		Data: mustJSON(trajectory.TranscriptMessageData{
			Text:      "hi",
			ToolUseID: "toolu_01",
		}),
	}}
	require.NoError(t, trajectory.AppendImported(root, key, imported))

	events, err := trajectory.ReadEvents(path)
	require.NoError(t, err, "ReadEvents")
	require.Len(t, events, 2)
	assert.Equal(t, uint64(1), events[0].Seq)
	assert.Equal(t, trajectory.TypeHookInvoked, events[0].Type)
	assert.Equal(t, uint64(2), events[1].Seq)
	assert.Equal(t, trajectory.TypeTranscriptMessage, events[1].Type)
}

func TestShowExportIncludesTranscriptEvents(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	key := trajectory.SessionKey{Provider: provider.ClaudeCode, SessionID: "s2"}
	imported := []trajectory.Event{{
		Type:   trajectory.TypeTranscriptThinking,
		Source: trajectory.SourceTranscript,
		TS:     time.Now().UTC(),
		Data:   mustJSON(trajectory.TranscriptThinkingData{Text: "reason"}),
	}}
	require.NoError(t, trajectory.AppendImported(root, key, imported))

	path := trajectory.SessionFilePath(root, key)
	events, err := trajectory.ReadEvents(path)
	require.NoError(t, err, "ReadEvents")
	require.Len(t, events, 1)
	assert.Equal(t, trajectory.SourceTranscript, events[0].Source)

	outPath := filepath.Join(root, "out.jsonl")
	require.NoError(t, trajectory.ExportToFile(outPath, root, "claude-code", "s2"))
	b, err := os.ReadFile(outPath)
	require.NoError(t, err, "ReadFile export")
	assert.Contains(t, string(b), "transcript/thinking")
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return b
}
