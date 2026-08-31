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

func TestTargetsFromHighConfidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(t *testing.T, env *install.DiscoverEnv)
		wantCount int
		wantScope []string
	}{
		{
			name: "scope_both_claude",
			setup: func(t *testing.T, env *install.DiscoverEnv) {
				require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".claude"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(env.Home, ".claude"), 0o755))
			},
			wantCount: 2,
			wantScope: []string{"project", "user"},
		},
		{
			name: "kimi_forces_user",
			setup: func(t *testing.T, env *install.DiscoverEnv) {
				require.NoError(t, os.MkdirAll(filepath.Join(env.Home, ".kimi-code"), 0o755))
			},
			wantCount: 1,
			wantScope: []string{"user"},
		},
		{
			name: "opencode_project_only",
			setup: func(t *testing.T, env *install.DiscoverEnv) {
				require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".opencode"), 0o755))
			},
			wantCount: 1,
			wantScope: []string{"project"},
		},
		{
			name: "plugin_never_auto",
			setup: func(t *testing.T, env *install.DiscoverEnv) {
				env.LookPath = func(file string) (string, error) {
					if file == "cursor-agent" || file == "cursor" {
						return "/bin/" + file, nil
					}
					return "", os.ErrNotExist
				}
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			env := testDiscoverEnv(t, root)
			tt.setup(t, &env)
			findings, err := install.Discover(context.Background(), env)
			require.NoError(t, err, "Discover(%q)", tt.name)
			got, err := install.TargetsFromHighConfidence(findings, env)
			require.NoError(t, err, "TargetsFromHighConfidence(%q)", tt.name)
			if tt.wantCount == 0 {
				assert.Empty(t, got)
				return
			}
			require.Len(t, got, tt.wantCount)
			scopes := make([]string, len(got))
			for i, tg := range got {
				scopes[i] = tg.Scope
				assert.NotEmpty(t, tg.Dir)
			}
			assert.ElementsMatch(t, tt.wantScope, scopes)
			if tt.name == "kimi_forces_user" {
				assert.Equal(t, provider.KimiCode, got[0].Provider)
			}
		})
	}
}
