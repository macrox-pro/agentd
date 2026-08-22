package trajectory_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestProviderImporterStatusTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		provider string
		want     trajectory.ImporterStatus
	}{
		{name: "claude-code", provider: "claude-code", want: trajectory.ImporterSupported},
		{name: "cursor", provider: "cursor", want: trajectory.ImporterPartial},
		{name: "codex", provider: "codex", want: trajectory.ImporterSupported},
		{name: "gemini", provider: "gemini", want: trajectory.ImporterNone},
		{name: "opencode", provider: "opencode", want: trajectory.ImporterNone},
		{name: "kimi-code", provider: "kimi-code", want: trajectory.ImporterNone},
		{name: "kimi alias", provider: "kimicode", want: trajectory.ImporterNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
