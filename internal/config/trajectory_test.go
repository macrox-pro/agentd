package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestCompileTrajectoryDefaultsOff(t *testing.T) {
	t.Parallel()
	res, err := config.CompileMerged(nil, nil, nil)
	require.NoError(t, err)
	assert.False(t, res.Trajectory.Enabled, "default off")
	assert.False(t, res.Trajectory.Statistics, "default statistics off")
	assert.False(t, res.Trajectory.IncludeRaw)
	assert.True(t, res.Trajectory.RedactSecretRules)
	assert.Equal(t, 262144, res.Trajectory.MaxEventBytes)
}

func TestCompileTrajectoryFromYAML(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  enabled: true
  include_raw: true
  max_event_bytes: 4096
`), 0o600))
	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	snap := store.Current()
	assert.True(t, snap.Trajectory.Enabled)
	assert.True(t, snap.Trajectory.IncludeRaw)
	assert.Equal(t, 4096, snap.Trajectory.MaxEventBytes)
}

func TestCompileTrajectoryInvalidMaxBytes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  max_event_bytes: -1
`), 0o600))
	_, err := config.Load(context.Background(), path)
	require.Error(t, err)
}

func TestFingerprintIncludesTrajectory(t *testing.T) {
	t.Parallel()
	a, err := config.CompileMerged(nil, nil, nil)
	require.NoError(t, err)
	fpA, err := config.Fingerprint(a.Merged)
	require.NoError(t, err)

	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  enabled: true
`), 0o600))
	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	fpB := store.Current().Fingerprint
	assert.NotEqual(t, fpA, fpB)
}

func TestParseTrajectory_statistics(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		yaml string
		want bool
	}{
		{
			name: "default_false",
			yaml: `version: 1
trajectory:
  enabled: true`,
			want: false,
		},
		{
			name: "enabled_and_statistics",
			yaml: `version: 1
trajectory:
  enabled: true
  statistics: true`,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "agentd.yaml")
			require.NoError(t, os.WriteFile(path, []byte(tt.yaml), 0o600))
			store, err := config.Load(context.Background(), path)
			require.NoError(t, err)
			assert.Equal(t, tt.want, store.Current().Trajectory.Statistics)
		})
	}
}

func TestFingerprintIncludesTrajectoryStatistics(t *testing.T) {
	t.Parallel()
	a, err := config.CompileMerged(nil, nil, nil)
	require.NoError(t, err)
	fpA, err := config.Fingerprint(a.Merged)
	require.NoError(t, err)

	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  enabled: true
  statistics: true
`), 0o600))
	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	fpB := store.Current().Fingerprint
	assert.NotEqual(t, fpA, fpB)
}

func TestFingerprintIncludesTrajectoryImport(t *testing.T) {
	t.Parallel()
	a, err := config.CompileMerged(nil, nil, nil)
	require.NoError(t, err)
	fpA, err := config.Fingerprint(a.Merged)
	require.NoError(t, err)

	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  import:
    claude-code:
      enabled: true
      path: /tmp/claude
`), 0o600))
	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	fpB := store.Current().Fingerprint
	assert.NotEqual(t, fpA, fpB)
	assert.True(t, store.Current().Trajectory.ClaudeImport().Enabled)
}

func TestCompileTrajectoryImportCursorCodex(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  import:
    cursor:
      enabled: true
      path: /tmp/cursor
    codex:
      enabled: true
      path: /tmp/codex
`), 0o600))
	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	assert.True(t, store.Current().Trajectory.CursorImport().Enabled)
	assert.Equal(t, "/tmp/cursor", store.Current().Trajectory.CursorImport().Path)
	assert.True(t, store.Current().Trajectory.CodexImport().Enabled)
	assert.Equal(t, "/tmp/codex", store.Current().Trajectory.CodexImport().Path)
}

func TestCompileTrajectoryImportProviderAlias(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  import:
    Claude-Code:
      enabled: true
      path: /tmp/claude-alias
`), 0o600))
	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	assert.True(t, store.Current().Trajectory.ClaudeImport().Enabled)
	assert.Equal(t, "/tmp/claude-alias", store.Current().Trajectory.ClaudeImport().Path)
}
