package daemon_test

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestSetupLogWritesFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "agentd.log")

	log, cleanup, err := daemon.SetupLog(daemon.SetupLogOptions{
		Logging: config.LoggingConfig{
			Level: config.LogLevelInfo,
			File:  logPath,
		},
	})
	require.NoError(t, err, "SetupLog")
	require.NotNil(t, cleanup)
	t.Cleanup(cleanup)

	log.Info("daemon ready", "socket", "/tmp/s.sock")
	cleanup()

	b, err := os.ReadFile(logPath)
	require.NoError(t, err, "ReadFile")
	assert.Contains(t, string(b), "daemon ready", "log line")
	assert.Contains(t, string(b), "/tmp/s.sock", "attr")
}

func TestSetupLogLevelFilter(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "agentd.log")

	log, cleanup, err := daemon.SetupLog(daemon.SetupLogOptions{
		Logging: config.LoggingConfig{
			Level: config.LogLevelWarn,
			File:  logPath,
		},
	})
	require.NoError(t, err, "SetupLog")
	t.Cleanup(cleanup)

	log.Info("hidden")
	log.Warn("visible")
	cleanup()

	b, err := os.ReadFile(logPath)
	require.NoError(t, err, "ReadFile")
	body := string(b)
	assert.NotContains(t, body, "hidden", "info filtered")
	assert.Contains(t, body, "visible", "warn kept")
}

func TestSetupLogForegroundMirrorsStderr(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "agentd.log")
	old := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err, "Pipe")
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = old
		_ = w.Close()
		_ = r.Close()
	})

	log, cleanup, err := daemon.SetupLog(daemon.SetupLogOptions{
		Logging: config.LoggingConfig{
			Level: config.LogLevelInfo,
			File:  logPath,
		},
		Foreground: true,
	})
	require.NoError(t, err, "SetupLog")
	t.Cleanup(cleanup)

	log.Info("mirror me")
	_ = w.Close()

	var stderr bytes.Buffer
	_, _ = io.Copy(&stderr, r)
	cleanup()

	assert.Contains(t, stderr.String(), "mirror me", "stderr")
	b, err := os.ReadFile(logPath)
	require.NoError(t, err, "ReadFile")
	assert.Contains(t, string(b), "mirror me", "file")
}

func TestSetupLogCLILevelOverride(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "agentd.log")

	log, cleanup, err := daemon.SetupLog(daemon.SetupLogOptions{
		Logging:  config.LoggingConfig{Level: config.LogLevelWarn, File: logPath},
		LogLevel: "debug",
	})
	require.NoError(t, err, "SetupLog")
	t.Cleanup(cleanup)

	log.Debug("debug line")
	cleanup()

	b, err := os.ReadFile(logPath)
	require.NoError(t, err, "ReadFile")
	assert.Contains(t, string(b), "debug line", "debug visible with override")
}

func TestSetupLogInvalidLevel(t *testing.T) {
	dir := t.TempDir()
	_, _, err := daemon.SetupLog(daemon.SetupLogOptions{
		Logging:  config.LoggingConfig{Level: config.LogLevelInfo, File: filepath.Join(dir, "x.log")},
		LogLevel: "nope",
	})
	require.Error(t, err, "SetupLog invalid level")
}

func TestSetupLogRestoresDefault(t *testing.T) {
	dir := t.TempDir()
	prev := slog.Default()

	log, cleanup, err := daemon.SetupLog(daemon.SetupLogOptions{
		Logging: config.LoggingConfig{
			Level: config.LogLevelInfo,
			File:  filepath.Join(dir, "agentd.log"),
		},
	})
	require.NoError(t, err, "SetupLog")
	assert.NotSame(t, prev, slog.Default(), "SetDefault")
	assert.Same(t, log, slog.Default(), "returned logger is default")

	cleanup()
	assert.Same(t, prev, slog.Default(), "restored default")
}
