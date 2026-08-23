package install_test

import (
	"testing"

	ahinstall "github.com/speakeasy-api/agenthooks/install"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
	"github.com/macrox-pro/agentd/internal/provider"
)

func TestResolveDir(t *testing.T) {
	t.Parallel()

	const (
		cwd  = "/work/repo"
		home = "/home/me"
	)

	tests := []struct {
		name       string
		id         provider.ID
		scope      ahinstall.Scope
		dirFlag    string
		dirFlagSet bool
		env        map[string]string
		home       string
		want       string
		wantErr    error
	}{
		{
			name:       "user cursor defaults to home/.cursor",
			id:         provider.Cursor,
			scope:      ahinstall.ScopeUser,
			home:       home,
			want:       "/home/me/.cursor",
		},
		{
			name:       "user respects explicit --dir",
			id:         provider.Cursor,
			scope:      ahinstall.ScopeUser,
			dirFlag:    "/custom/cursor",
			dirFlagSet: true,
			want:       "/custom/cursor",
		},
		{
			name:       "project cursor uses cwd",
			id:         provider.Cursor,
			scope:      ahinstall.ScopeProject,
			want:       cwd,
		},
		{
			name:       "project codex defaults to cwd/.codex",
			id:         provider.Codex,
			scope:      ahinstall.ScopeProject,
			want:       "/work/repo/.codex",
		},
		{
			name:       "user codex uses CODEX_HOME",
			id:         provider.Codex,
			scope:      ahinstall.ScopeUser,
			env:        map[string]string{"CODEX_HOME": "/codex-home"},
			home:       home,
			want:       "/codex-home",
		},
		{
			name:       "user kimi uses KIMI_CODE_HOME",
			id:         provider.KimiCode,
			scope:      ahinstall.ScopeUser,
			env:        map[string]string{"KIMI_CODE_HOME": "/kimi-home"},
			home:       home,
			want:       "/kimi-home",
		},
		{
			name:       "plugin without dir errors",
			id:         provider.Cursor,
			scope:      ahinstall.ScopePlugin,
			wantErr:    install.ErrDirRequired,
		},
		{
			name:       "opencode user without dir errors",
			id:         provider.OpenCode,
			scope:      ahinstall.ScopeUser,
			home:       home,
			wantErr:    install.ErrDirRequired,
		},
		{
			name:       "HOME unset user dir error",
			id:         provider.Cursor,
			scope:      ahinstall.ScopeUser,
			home:       "",
			wantErr:    install.ErrHomeRequired,
		},
		{
			name:       "dir flag Changed skips defaults",
			id:         provider.ClaudeCode,
			scope:      ahinstall.ScopeUser,
			dirFlag:    "/explicit/.claude",
			dirFlagSet: true,
			want:       "/explicit/.claude",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			getenv := func(key string) string {
				if tt.env == nil {
					return ""
				}
				return tt.env[key]
			}
			home := tt.home
			if home == "" && tt.wantErr != install.ErrHomeRequired {
				home = "/home/me"
			}
			got, err := install.ResolveDir(tt.id, tt.scope, tt.dirFlag, tt.dirFlagSet, cwd, home, getenv)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
