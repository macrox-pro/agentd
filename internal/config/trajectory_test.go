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

func TestCompileTrajectoryDefaultsOn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		yaml string
		want config.TrajectoryConfig
	}{
		{
			name: "defaults_on",
			yaml: "",
			want: config.TrajectoryConfig{
				Enabled:           true,
				Statistics:        true,
				IncludeRaw:        true,
				RedactSecretRules: true,
				MaxEventBytes:     262144,
			},
		},
		{
			name: "explicit_enabled_false",
			yaml: `version: 1
trajectory:
  enabled: false`,
			want: config.TrajectoryConfig{
				Enabled:           false,
				Statistics:        true,
				IncludeRaw:        true,
				RedactSecretRules: true,
				MaxEventBytes:     262144,
			},
		},
		{
			name: "statistics_false_only",
			yaml: `version: 1
trajectory:
  statistics: false`,
			want: config.TrajectoryConfig{
				Enabled:           true,
				Statistics:        false,
				IncludeRaw:        true,
				RedactSecretRules: true,
				MaxEventBytes:     262144,
			},
		},
		{
			name: "include_raw_false_only",
			yaml: `version: 1
trajectory:
  include_raw: false`,
			want: config.TrajectoryConfig{
				Enabled:           true,
				Statistics:        true,
				IncludeRaw:        false,
				RedactSecretRules: true,
				MaxEventBytes:     262144,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var snap *config.Snapshot
			if tt.yaml == "" {
				res, err := config.CompileMerged(nil, nil, nil)
				require.NoError(t, err, "CompileMerged(%q)", tt.name)
				snap = &config.Snapshot{Trajectory: res.Trajectory}
			} else {
				dir := t.TempDir()
				path := filepath.Join(dir, "agentd.yaml")
				require.NoError(t, os.WriteFile(path, []byte(tt.yaml), 0o600))
				store, err := config.Load(context.Background(), path)
				require.NoError(t, err, "Load(%q)", tt.name)
				snap = store.Current()
			}
			assert.Equal(t, tt.want.Enabled, snap.Trajectory.Enabled, "Enabled")
			assert.Equal(t, tt.want.Statistics, snap.Trajectory.Statistics, "Statistics")
			assert.Equal(t, tt.want.IncludeRaw, snap.Trajectory.IncludeRaw, "IncludeRaw")
			assert.Equal(t, tt.want.RedactSecretRules, snap.Trajectory.RedactSecretRules, "RedactSecretRules")
			assert.Equal(t, tt.want.MaxEventBytes, snap.Trajectory.MaxEventBytes, "MaxEventBytes")
		})
	}
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
			name: "omitted_uses_compile_default",
			yaml: `version: 1
trajectory:
  enabled: true`,
			want: true,
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
