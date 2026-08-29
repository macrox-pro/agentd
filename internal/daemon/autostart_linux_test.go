//go:build linux

package daemon_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestAutostartLinux_writes_unit_file(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	unitDir := filepath.Join(home, ".config", "systemd", "user")
	require.NoError(t, os.MkdirAll(unitDir, 0o700))
	path := filepath.Join(unitDir, "agentd.service")
	body := daemon.RenderSystemdUnitForTest(daemon.AutostartSpecForTest("/usr/bin/agentd", []string{"daemon", "start", "--foreground"}))
	require.NoError(t, daemon.WriteFileAtomicForTest(path, []byte(body), 0o600))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "ExecStart=")
}

func TestAutostartLinux_preflight_no_runtime_dir(t *testing.T) {
	t.Setenv("XDG_RUNTIME_DIR", "")
	home := t.TempDir()
	t.Setenv("HOME", home)

	err := daemon.Enable(t.Context(), daemon.AutostartOptions{})
	require.Error(t, err)
	assert.ErrorIs(t, err, daemon.ErrAutostartNotAvailable)
}
