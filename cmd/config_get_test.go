package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigGet(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	configPath := filepath.Join(dir, "user.yaml")

	off := executeRoot(t, execOpts{args: []string{"config", "get", "trajectory"}, configPath: configPath})
	require.NoError(t, off.err, "config get default")
	assert.Contains(t, off.out, "trajectory: off (default)", "get_default_off")

	on := executeRoot(t, execOpts{args: []string{"config", "enable", "trajectory"}, configPath: configPath})
	require.NoError(t, on.err, "config enable")

	got := executeRoot(t, execOpts{args: []string{"config", "get", "trajectory"}, configPath: configPath})
	require.NoError(t, got.err, "config get after enable")
	assert.True(t, strings.Contains(got.out, "trajectory: on (user)"), "get_after_enable: %q", got.out)
}
