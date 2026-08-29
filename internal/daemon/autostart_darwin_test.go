//go:build darwin

package daemon_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestAutostartDarwin_writes_plist(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	launchAgents := filepath.Join(home, "Library", "LaunchAgents")
	require.NoError(t, os.MkdirAll(launchAgents, 0o700))

	path := filepath.Join(launchAgents, "io.github.macrox-pro.agentd.plist")
	body := daemon.RenderLaunchdPlistForTest(daemon.AutostartSpecForTest("/tmp/agentd", []string{"daemon", "start", "--foreground"}))
	require.NoError(t, daemon.WriteFileAtomicForTest(path, []byte(body), 0o600))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "io.github.macrox-pro.agentd")
}
