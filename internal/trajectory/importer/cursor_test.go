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

func TestImportCursorMapsMessagesSkipsRedacted(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "cursor_transcript.jsonl")
	result, err := importer.ImportCursor(importer.ImportOptions{
		SessionID:      "m11-cursor",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{MaxEventBytes: 262144},
	})
	require.NoError(t, err, "ImportCursor")
	require.NotEmpty(t, result.Events)

	var messages int
	for _, e := range result.Events {
		assert.Equal(t, trajectory.SourceTranscript, e.Source)
		assert.Equal(t, "cursor", e.Provider)
		assert.NotEqual(t, trajectory.TypeTranscriptThinking, e.Type, "cursor must not invent thinking")
		if e.Type == trajectory.TypeTranscriptMessage {
			messages++
			var d trajectory.TranscriptMessageData
			require.NoError(t, json.Unmarshal(e.Data, &d))
			assert.NotEqual(t, "[REDACTED]", d.Text)
		}
	}
	assert.GreaterOrEqual(t, messages, 2)
}

func TestImportCursorRequiresPath(t *testing.T) {
	t.Parallel()
	_, err := importer.ImportCursor(importer.ImportOptions{SessionID: "s1"})
	require.Error(t, err)
	assert.ErrorIs(t, err, importer.ErrTranscriptRootRequired)
}

func TestImportCursorEmptyFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	require.NoError(t, os.WriteFile(path, nil, 0o600))
	result, err := importer.ImportCursor(importer.ImportOptions{
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err)
	assert.Empty(t, result.Events)
}
