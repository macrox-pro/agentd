package trajectory_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestProviderImporterStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		provider string
		want     trajectory.ImporterStatus
	}{
		{provider: "claude-code", want: trajectory.ImporterSupported},
		{provider: "cursor", want: trajectory.ImporterPartial},
		{provider: "codex", want: trajectory.ImporterSupported},
		{provider: "gemini", want: trajectory.ImporterNone},
		{provider: "opencode", want: trajectory.ImporterNone},
		{provider: "kimi-code", want: trajectory.ImporterNone},
	}
	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, trajectory.ProviderImporterStatus(tt.provider))
		})
	}
}

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
