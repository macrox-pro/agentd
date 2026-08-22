package importer_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

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

	require.NoError(t, trajectory.SaveImportCheckpoint(sidecar, trajectory.ImportCheckpoint{LastLineIndex: first.LastLineIndex}))

	second, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      sid,
		TranscriptPath: path,
		StartIndex:     first.LastLineIndex,
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
