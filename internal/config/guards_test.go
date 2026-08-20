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

func TestGuardsShellMCPPathsDefaultsAndMerge(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("defaults disabled", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "agentd.yaml")
		require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600))
		store, err := config.Load(ctx, path)
		require.NoError(t, err)
		snap := store.Current()
		assert.True(t, snap.Guards.Secrets.Enabled, "secrets default")
		assert.False(t, snap.Guards.Shell.Enabled, "shell default")
		assert.False(t, snap.Guards.MCP.Enabled, "mcp default")
		assert.False(t, snap.Guards.Paths.Enabled, "paths default")
	})

	t.Run("project overlay disables shell", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		userPath := filepath.Join(dir, "user.yaml")
		require.NoError(t, os.WriteFile(userPath, []byte(`version: 1
guards:
  shell:
    enabled: true
    deny_patterns: ["rm -rf /"]
    ask_on: [curl]
  mcp:
    enabled: true
    deny_servers: ["bad-*"]
  paths:
    enabled: true
    deny_read: ["/etc/shadow"]
`), 0o600))

		projDir := filepath.Join(dir, "proj")
		require.NoError(t, os.MkdirAll(projDir, 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(projDir, ".agentd.yaml"), []byte(`version: 1
guards:
  shell:
    enabled: false
`), 0o600))

		store, err := config.Load(ctx, userPath)
		require.NoError(t, err)
		snap := store.SnapshotFor(projDir, "")
		assert.False(t, snap.Guards.Shell.Enabled, "project disables shell")
		assert.Equal(t, []string{"rm -rf /"}, snap.Guards.Shell.DenyPatterns, "patterns retained")
		assert.True(t, snap.Guards.MCP.Enabled, "mcp from user")
		assert.True(t, snap.Guards.Paths.Enabled, "paths from user")
	})

	t.Run("default routes include enabled guards", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "agentd.yaml")
		require.NoError(t, os.WriteFile(path, []byte(`version: 1
guards:
  shell:
    enabled: true
  mcp:
    enabled: true
  paths:
    enabled: true
`), 0o600))
		store, err := config.Load(ctx, path)
		require.NoError(t, err)
		var pre *config.CompiledRoute
		for i := range store.Current().Routes {
			r := &store.Current().Routes[i]
			if r.Name == "default-tool.pre" {
				pre = r
				break
			}
		}
		require.NotNil(t, pre)
		require.Len(t, pre.Sync, 1)
		assert.Equal(t, []string{"secrets", "shell", "mcp", "paths"}, pre.Sync[0].Guards)
	})
}
