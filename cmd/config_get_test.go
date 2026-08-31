package cmd_test

import (
	"os"
	"path/filepath"
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
	assert.Contains(t, off.out, "trajectory: on (default)", "get_default_on")

	on := executeRoot(t, execOpts{args: []string{"config", "enable", "trajectory"}, configPath: configPath})
	require.NoError(t, on.err, "config enable")
	assert.Contains(t, on.out, "already enabled", "enable_idempotent_default_on")

	got := executeRoot(t, execOpts{args: []string{"config", "get", "trajectory"}, configPath: configPath})
	require.NoError(t, got.err, "config get after enable")
	assert.Contains(t, got.out, "trajectory: on (default)", "get_after_enable: %q", got.out)
}
