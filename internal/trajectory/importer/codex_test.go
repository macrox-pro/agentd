package importer_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestImportCodexMapsMessages(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "codex_transcript.jsonl")
	result, err := importer.ImportCodex(importer.ImportOptions{
		SessionID:      "m11-codex",
		TranscriptPath: path,
		Cfg:            config.TrajectoryConfig{},
	})
	require.NoError(t, err, "ImportCodex")
	require.NotEmpty(t, result.Events)
	for _, e := range result.Events {
		assert.Equal(t, "codex", e.Provider)
		assert.Equal(t, trajectory.SourceTranscript, e.Source)
	}
}
