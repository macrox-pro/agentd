package install_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
	"github.com/macrox-pro/agentd/internal/provider"
)

func testDiscoverEnv(t *testing.T, root string) install.DiscoverEnv {
	t.Helper()
	home := filepath.Join(root, "home")
	cwd := filepath.Join(root, "repo")
	require.NoError(t, os.MkdirAll(home, 0o755))
	require.NoError(t, os.MkdirAll(cwd, 0o755))
	return install.DiscoverEnv{
		Cwd:  cwd,
		Home: home,
		Stat: os.Stat,
		LookPath: func(string) (string, error) {
			return "", os.ErrNotExist
		},
	}
}

func TestDiscover(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(t *testing.T, root string)
		wantIDs  []provider.ID
		wantConf install.Confidence
	}{
		{
			name: "empty_discover",
			setup: func(t *testing.T, root string) {
				_ = root
			},
			wantIDs: nil,
		},
		{
			name: "cursor_project",
			setup: func(t *testing.T, root string) {
				env := testDiscoverEnv(t, root)
				require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".cursor"), 0o755))
			},
			wantIDs:  []provider.ID{provider.Cursor},
			wantConf: install.ConfidenceHigh,
		},
		{
			name: "cursor_user",
			setup: func(t *testing.T, root string) {
				env := testDiscoverEnv(t, root)
				require.NoError(t, os.MkdirAll(filepath.Join(env.Home, ".cursor"), 0o755))
			},
			wantIDs:  []provider.ID{provider.Cursor},
			wantConf: install.ConfidenceHigh,
		},
		{
			name: "codex_dot_codex",
			setup: func(t *testing.T, root string) {
				env := testDiscoverEnv(t, root)
				require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".codex"), 0o755))
			},
			wantIDs:  []provider.ID{provider.Codex},
			wantConf: install.ConfidenceHigh,
		},
		{
			name: "codex_home_env",
			setup: func(t *testing.T, root string) {
				codexHome := filepath.Join(root, "codex-home")
				require.NoError(t, os.MkdirAll(codexHome, 0o755))
			},
			wantIDs:  []provider.ID{provider.Codex},
			wantConf: install.ConfidenceHigh,
		},
		{
			name: "kimi_no_project",
			setup: func(t *testing.T, root string) {
				env := testDiscoverEnv(t, root)
				require.NoError(t, os.MkdirAll(filepath.Join(env.Home, ".kimi-code"), 0o755))
			},
			wantIDs:  []provider.ID{provider.KimiCode},
			wantConf: install.ConfidenceHigh,
		},
		{
			name: "opencode_no_user",
			setup: func(t *testing.T, root string) {
				env := testDiscoverEnv(t, root)
				require.NoError(t, os.MkdirAll(filepath.Join(env.Home, ".opencode"), 0o755))
			},
			wantIDs: nil,
		},
		{
			name: "claude_binary_only_medium",
			setup: func(t *testing.T, root string) {
				_ = root
			},
			wantIDs:  []provider.ID{provider.ClaudeCode},
			wantConf: install.ConfidenceMedium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			env := testDiscoverEnv(t, root)
			if tt.name == "codex_home_env" {
				codexHome := filepath.Join(root, "codex-home")
				env.Getenv = func(key string) string {
					if key == "CODEX_HOME" {
						return codexHome
					}
					return os.Getenv(key)
				}
			}
			if tt.name == "claude_binary_only_medium" {
				env.LookPath = func(file string) (string, error) {
					if file == "claude" {
						return "/usr/bin/claude", nil
					}
					return "", os.ErrNotExist
				}
			}
			tt.setup(t, root)
			got, err := install.Discover(context.Background(), env)
			require.NoError(t, err, "Discover(%q)", tt.name)
			if len(tt.wantIDs) == 0 {
				assert.Empty(t, got, "Discover(%q)", tt.name)
				return
			}
			require.Len(t, got, len(tt.wantIDs), "Discover(%q) count", tt.name)
			for i, id := range tt.wantIDs {
				assert.Equal(t, id, got[i].Provider, "Discover(%q) provider", tt.name)
				if tt.wantConf != "" {
					assert.Equal(t, tt.wantConf, got[i].Confidence, "Discover(%q) confidence", tt.name)
				}
			}
		})
	}
}

func TestDiscover_noChdir(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	env := testDiscoverEnv(t, root)
	other := filepath.Join(root, "other")
	require.NoError(t, os.MkdirAll(other, 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(other, ".cursor"), 0o755))
	env.Cwd = other
	got, err := install.Discover(context.Background(), env)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, provider.Cursor, got[0].Provider)
}
