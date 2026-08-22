package config_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestCompileLoggingDefaults(t *testing.T) {
	t.Parallel()

	store, err := config.Load(context.Background(), filepath.Join(t.TempDir(), "missing.yaml"))
	require.NoError(t, err, "Load")
	snap := store.Current()
	assert.Equal(t, config.LogLevelInfo, snap.Logging.Level, "level")
	assert.Empty(t, snap.Logging.File, "file")
}

func TestCompileLoggingUserOverride(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "agentd.yaml")
	content := `version: 1
logging:
  level: warn
  file: /tmp/custom.log
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	store, err := config.Load(context.Background(), path)
	require.NoError(t, err, "Load")
	snap := store.Current()
	assert.Equal(t, config.LogLevelWarn, snap.Logging.Level, "level")
	assert.Equal(t, "/tmp/custom.log", snap.Logging.File, "file")
}

func TestCompileLoggingInvalidLevel(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "agentd.yaml")
	content := `version: 1
logging:
  level: trace
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	_, err := config.Load(context.Background(), path)
	require.Error(t, err, "Load invalid level")
}

func TestLoggingEffectiveLevel(t *testing.T) {
	t.Parallel()

	cfg := config.LoggingConfig{Level: config.LogLevelInfo}
	lvl, err := cfg.EffectiveLevel("debug")
	require.NoError(t, err, "EffectiveLevel")
	assert.Equal(t, slog.LevelDebug, lvl, "override")

	_, err = cfg.EffectiveLevel("nope")
	require.Error(t, err, "EffectiveLevel invalid")
}

func TestLoggingEffectiveFile(t *testing.T) {
	t.Parallel()

	cfg := config.LoggingConfig{File: "/var/log/agentd.log"}
	assert.Equal(t, "/tmp/override.log", cfg.EffectiveFile("/tmp/override.log"), "cli override")
	assert.Equal(t, "/var/log/agentd.log", cfg.EffectiveFile(""), "config file")
}

func TestDefaultLogPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	got := config.DefaultLogPath()
	require.NotEmpty(t, got, "DefaultLogPath")
	assert.Contains(t, got, "agentd.log", "filename")
}
