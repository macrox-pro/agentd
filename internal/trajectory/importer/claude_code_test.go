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

func TestImportClaudeMapsThinkingAndToolUse(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "claude_session.jsonl")
	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      "m10-claude",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{MaxEventBytes: 262144, RedactSecretRules: true},
	})
	require.NoError(t, err, "ImportClaude")
	require.GreaterOrEqual(t, len(result.Events), 3)

	var thinking, toolMsg bool
	for _, e := range result.Events {
		assert.Equal(t, trajectory.SourceTranscript, e.Source)
		switch e.Type {
		case trajectory.TypeTranscriptThinking:
			thinking = true
		case trajectory.TypeTranscriptMessage:
			if e.Data != nil {
				var d trajectory.TranscriptMessageData
				require.NoError(t, json.Unmarshal(e.Data, &d))
				if d.ToolUseID == "toolu_01" {
					toolMsg = true
				}
			}
		}
	}
	assert.True(t, thinking, "thinking event")
	assert.True(t, toolMsg, "tool_use correlation")
}

func TestImportClaudeDedupCheckpoint(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join("testdata", "claude_session.jsonl")
	sid := "m10-dedup"
	sidecar := trajectory.ImportSidecarPath(root, "claude-code", sid)

	first, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      sid,
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude first")
	require.NotEmpty(t, first.Events)

	require.NoError(t, trajectory.SaveImportCheckpoint(sidecar, trajectory.ImportCheckpoint{
		LastLineIndex: first.LastLineIndex,
		SourcePath:    path,
	}))

	second, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      sid,
		TranscriptPath: path,
		StartIndex:     first.LastLineIndex + 1,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude second")
	assert.Empty(t, second.Events)
}

func TestResolveClaudeTranscriptExplicitPath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	f := filepath.Join(dir, "sess.jsonl")
	require.NoError(t, os.WriteFile(f, []byte("{}"), 0o600))
	got, err := importer.ResolveClaudeTranscriptPath("", f, "")
	require.NoError(t, err, "ResolveClaudeTranscriptPath")
	assert.Equal(t, f, got)
}

func TestResolveClaudeTranscriptExactMatch(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	uuid := "019f9ed8-c891-7dd0-9808-e31c3b38ce48"
	mainDir := filepath.Join(root, "-Users-test-project")
	agentDir := filepath.Join(root, "-Users-other")
	require.NoError(t, os.MkdirAll(mainDir, 0o700))
	require.NoError(t, os.MkdirAll(agentDir, 0o700))

	mainPath := filepath.Join(mainDir, uuid+".jsonl")
	agentPath := filepath.Join(agentDir, "agent-"+uuid+".jsonl")
	require.NoError(t, os.WriteFile(mainPath, []byte("{}"), 0o600))
	require.NoError(t, os.WriteFile(agentPath, []byte("{}"), 0o600))

	got, err := importer.ResolveClaudeTranscriptPath(uuid, "", root)
	require.NoError(t, err, "ResolveClaudeTranscriptPath(%q)", uuid)
	assert.Equal(t, mainPath, got)

	agentID := "agent-" + uuid
	gotAgent, err := importer.ResolveClaudeTranscriptPath(agentID, "", root)
	require.NoError(t, err, "ResolveClaudeTranscriptPath(%q)", agentID)
	assert.Equal(t, agentPath, gotAgent)
}

func TestResolveClaudeTranscriptNewestMatch(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	uuid := "dup-session-id"
	oldDir := filepath.Join(root, "-Users-old")
	newDir := filepath.Join(root, "-Users-new")
	require.NoError(t, os.MkdirAll(oldDir, 0o700))
	require.NoError(t, os.MkdirAll(newDir, 0o700))

	oldPath := filepath.Join(oldDir, uuid+".jsonl")
	newPath := filepath.Join(newDir, uuid+".jsonl")
	require.NoError(t, os.WriteFile(oldPath, []byte("old"), 0o600))
	require.NoError(t, os.WriteFile(newPath, []byte("new"), 0o600))

	oldTime := time.Now().Add(-2 * time.Hour)
	newTime := time.Now().Add(-1 * time.Hour)
	require.NoError(t, os.Chtimes(oldPath, oldTime, oldTime))
	require.NoError(t, os.Chtimes(newPath, newTime, newTime))

	got, err := importer.ResolveClaudeTranscriptPath(uuid, "", root)
	require.NoError(t, err, "ResolveClaudeTranscriptPath")
	assert.Equal(t, newPath, got)
}

func TestImportClaudeSkipsMetadataTypes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "meta.jsonl")
	content := `{"type":"summary","message":{"role":"user","content":"summary text"}}
{"type":"system","message":{"role":"user","content":"system prompt"}}
{"type":"file-history-snapshot","message":{"role":"user","content":"snap"}}
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      "meta-skip",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude")
	assert.Empty(t, result.Events)
}

func TestImportClaudeSkipsToolResultUserLine(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "tool_result.jsonl")
	content := `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_x","content":"output"}]}}
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      "tool-result-skip",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude")
	assert.Empty(t, result.Events)
}

func TestImportClaudeUsesTimestampFromWrapper(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "ts.jsonl")
	wantTS := time.Date(2026, 7, 26, 14, 2, 0, 0, time.UTC)
	content := `{"type":"user","timestamp":"2026-07-26T14:02:00.000Z","message":{"role":"user","content":"timed"}}
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	before := time.Now().UTC()
	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      "ts-test",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude")
	require.Len(t, result.Events, 1)
	assert.Equal(t, wantTS, result.Events[0].TS)
	assert.True(t, result.Events[0].TS.Before(before), "event TS should come from JSON not import wall clock")
}

func TestImportClaudeTimestampFallback(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "claude_session.jsonl")
	before := time.Now().UTC()
	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      "fallback-ts",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude")
	require.NotEmpty(t, result.Events)
	after := time.Now().UTC()
	for _, e := range result.Events {
		assert.False(t, e.TS.Before(before.Add(-time.Second)))
		assert.False(t, e.TS.After(after.Add(time.Second)))
	}
}

func TestImportClaudeOnDiskFixture(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "claude_session_on_disk.jsonl")
	wantTS := time.Date(2026, 7, 26, 14, 2, 0, 0, time.UTC)
	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      "sess-on-disk",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude")

	var thinking, textMsg, toolMsg bool
	var cwdSet bool
	for _, e := range result.Events {
		if e.CWD == "/tmp/agentd-project" {
			cwdSet = true
		}
		switch e.Type {
		case trajectory.TypeTranscriptThinking:
			thinking = true
			assert.Equal(t, wantTS, e.TS)
		case trajectory.TypeTranscriptMessage:
			if e.Data != nil {
				var d trajectory.TranscriptMessageData
				require.NoError(t, json.Unmarshal(e.Data, &d))
				if d.Text == "hi there" {
					textMsg = true
					assert.Equal(t, wantTS, e.TS)
				}
				if d.ToolUseID == "toolu_01" {
					toolMsg = true
				}
			}
		}
	}
	assert.True(t, thinking, "thinking from on-disk fixture")
	assert.True(t, textMsg, "text from on-disk fixture")
	assert.True(t, toolMsg, "tool_use_id from on-disk fixture")
	assert.True(t, cwdSet, "cwd from wrapper")
}

func TestImportClaudeSkipsRedactedThinking(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "claude_session_redacted.jsonl")
	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      "redacted",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportClaude")
	for _, e := range result.Events {
		assert.NotEqual(t, trajectory.TypeTranscriptThinking, e.Type)
	}
	require.Len(t, result.Events, 1)
	assert.Equal(t, trajectory.TypeTranscriptMessage, result.Events[0].Type)
}
