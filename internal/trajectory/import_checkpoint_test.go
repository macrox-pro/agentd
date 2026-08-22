package trajectory_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestImportCheckpointRoundTrip(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "claude-code", "s1.import.json")
	cp := trajectory.ImportCheckpoint{LastLineIndex: 3, SourcePath: "/tmp/t.jsonl"}
	require.NoError(t, trajectory.SaveImportCheckpoint(path, cp))
	got, err := trajectory.LoadImportCheckpoint(path)
	require.NoError(t, err, "LoadImportCheckpoint")
	assert.Equal(t, 3, got.LastLineIndex)
	assert.Equal(t, "/tmp/t.jsonl", got.SourcePath)
}
