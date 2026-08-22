package importer_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestResolveCodexBySessionID(t *testing.T) {
	t.Parallel()
	sid := "019f9ed8-c891-7dd0-9808-e31c3b38ce48"

	t.Run("nested rollout newest", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		older := filepath.Join(root, "2026", "07", "26", "rollout-2026-07-26T10-00-00-"+sid+".jsonl")
		newer := filepath.Join(root, "2026", "07", "27", "rollout-2026-07-27T12-00-00-"+sid+".jsonl")
		require.NoError(t, os.MkdirAll(filepath.Dir(older), 0o700))
		require.NoError(t, os.MkdirAll(filepath.Dir(newer), 0o700))
		require.NoError(t, os.WriteFile(older, []byte("{}\n"), 0o600))
		require.NoError(t, os.WriteFile(newer, []byte("{}\n"), 0o600))
		past := time.Now().Add(-time.Hour)
		require.NoError(t, os.Chtimes(older, past, past))

		got, err := importer.ResolveCodexTranscriptPath(sid, "", root)
		require.NoError(t, err, "ResolveCodexTranscriptPath(%q)", sid)
		assert.Equal(t, newer, got)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_, err := importer.ResolveCodexTranscriptPath("missing-sid", "", root)
		require.Error(t, err, "ResolveCodexTranscriptPath(missing)")
	})

	t.Run("explicit path", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		path := filepath.Join(root, "custom.jsonl")
		require.NoError(t, os.WriteFile(path, []byte("{}\n"), 0o600))
		got, err := importer.ResolveCodexTranscriptPath(sid, path, root)
		require.NoError(t, err, "ResolveCodexTranscriptPath explicit")
		assert.Equal(t, path, got)
	})

	t.Run("exact sid.jsonl absent", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		// Old resolver looked for {sid}.jsonl; rollout layout must still resolve.
		path := filepath.Join(root, "2026", "07", "26", "rollout-2026-07-26T17-33-55-"+sid+".jsonl")
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
		require.NoError(t, os.WriteFile(path, []byte("{}\n"), 0o600))
		bare := filepath.Join(root, sid+".jsonl")
		require.NoFileExists(t, bare)

		got, err := importer.ResolveCodexTranscriptPath(sid, "", root)
		require.NoError(t, err, "ResolveCodexTranscriptPath rollout only")
		assert.Equal(t, path, got)
	})
}

func TestImportCodexRolloutMapsMessagesAndTools(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "codex_transcript.jsonl")
	result, err := importer.ImportCodex(importer.ImportOptions{
		SessionID:      "m11-codex",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex")
	require.NotEmpty(t, result.Events)

	var (
		userMsg, agentMsg, thinking bool
		toolCall, toolOut           bool
		customCall, customOut       bool
	)
	for _, e := range result.Events {
		assert.Equal(t, "codex", e.Provider)
		assert.Equal(t, trajectory.SourceTranscript, e.Source)
		switch e.Type {
		case trajectory.TypeTranscriptThinking:
			thinking = true
			var d trajectory.TranscriptThinkingData
			require.NoError(t, json.Unmarshal(e.Data, &d))
			assert.Contains(t, d.Text, "Planning")
		case trajectory.TypeTranscriptMessage:
			var d trajectory.TranscriptMessageData
			require.NoError(t, json.Unmarshal(e.Data, &d))
			switch {
			case d.Role == "user" && d.Text == "run the tests":
				userMsg = true
			case d.Role == "assistant" && d.Text == "running tests now":
				agentMsg = true
			case d.ToolUseID == "call_abc" && d.Text != "ok":
				toolCall = true
			case d.ToolUseID == "call_abc" && d.Text == "ok":
				toolOut = true
			case d.ToolUseID == "call_custom" && d.Text != "done":
				customCall = true
			case d.ToolUseID == "call_custom" && d.Text == "done":
				customOut = true
			}
		}
	}
	assert.True(t, userMsg, "user_message")
	assert.True(t, agentMsg, "agent_message")
	assert.True(t, thinking, "agent_reasoning")
	assert.True(t, toolCall, "function_call")
	assert.True(t, toolOut, "function_call_output")
	assert.True(t, customCall, "custom_tool_call")
	assert.True(t, customOut, "custom_tool_call_output")
}

func TestImportCodexSkipsEncryptedReasoning(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "codex_transcript.jsonl")
	result, err := importer.ImportCodex(importer.ImportOptions{
		SessionID:      "m11-codex",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex")
	for _, e := range result.Events {
		if e.Type != trajectory.TypeTranscriptThinking {
			continue
		}
		var d trajectory.TranscriptThinkingData
		require.NoError(t, json.Unmarshal(e.Data, &d))
		assert.NotContains(t, d.Text, "aabb")
		assert.NotEqual(t, "", d.Text)
	}
}

func TestImportCodexSessionIDFromMeta(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "codex_transcript.jsonl")
	result, err := importer.ImportCodex(importer.ImportOptions{
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex")
	require.NotEmpty(t, result.Events)
	assert.Equal(t, "019f9ed8-c891-7dd0-9808-e31c3b38ce48", result.Events[0].SessionID)
}

func TestImportCodexSessionIDFromFilenameFallback(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	sid := "019fa7c0-b628-7b82-97c2-40d94f5b2efd"
	path := filepath.Join(dir, "rollout-2026-07-28T11-04-12-"+sid+".jsonl")
	// No session_meta — only a user_message line.
	body := `{"type":"event_msg","payload":{"type":"user_message","message":"hello"}}` + "\n"
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))

	result, err := importer.ImportCodex(importer.ImportOptions{
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex")
	require.Len(t, result.Events, 1)
	assert.Equal(t, sid, result.Events[0].SessionID)
}

func TestImportCodexStartIndexResume(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "codex_transcript.jsonl")
	first, err := importer.ImportCodex(importer.ImportOptions{
		SessionID:      "resume-codex",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex first")
	require.NotEmpty(t, first.Events)

	second, err := importer.ImportCodex(importer.ImportOptions{
		SessionID:      "resume-codex",
		TranscriptPath: path,
		StartIndex:     first.LastLineIndex + 1,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex second")
	assert.Empty(t, second.Events)
}

func TestImportCodexEmptyFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	require.NoError(t, os.WriteFile(path, nil, 0o600))
	result, err := importer.ImportCodex(importer.ImportOptions{
		SessionID:      "empty",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex empty")
	assert.Empty(t, result.Events)
	assert.Equal(t, -1, result.LastLineIndex)
}
