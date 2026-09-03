package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestFindProjectConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T) (cwd, root string)
		wantFound  bool
		wantBase   string // expected basename dir of config ("" = any)
		wantInRoot bool
	}{
		{
			name: "missing",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				return t.TempDir(), ""
			},
			wantFound: false,
		},
		{
			name: "in cwd",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(dir, ".agentd.yaml"), []byte("version: 1\n"), 0o600))
				return dir, ""
			},
			wantFound: true,
		},
		{
			name: "in ancestor",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				root := t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(root, ".agentd.yaml"), []byte("version: 1\n"), 0o600))
				nested := filepath.Join(root, "a", "b")
				require.NoError(t, os.MkdirAll(nested, 0o700))
				return nested, ""
			},
			wantFound: true,
		},
		{
			name: "project_root wins over nearer cwd file",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				root := filepath.Join(base, "root")
				near := filepath.Join(base, "near")
				require.NoError(t, os.MkdirAll(root, 0o700))
				require.NoError(t, os.MkdirAll(near, 0o700))
				require.NoError(t, os.WriteFile(filepath.Join(root, ".agentd.yaml"), []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))
				require.NoError(t, os.WriteFile(filepath.Join(near, ".agentd.yaml"), []byte("version: 1\n"), 0o600))
				return near, root
			},
			wantFound:  true,
			wantInRoot: true,
		},
		{
			name: "empty start with project_root",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				root := t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(root, ".agentd.yaml"), []byte("version: 1\n"), 0o600))
				return "", root
			},
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cwd, root := tt.setup(t)
			got, ok := config.FindProjectConfig(cwd, root)
			assert.Equal(t, tt.wantFound, ok, "FindProjectConfig(%q,%q)", cwd, root)
			if !tt.wantFound {
				assert.Empty(t, got)
				return
			}
			require.True(t, filepath.IsAbs(got), "path=%q", got)
			assert.Equal(t, ".agentd.yaml", filepath.Base(got))
			if tt.wantInRoot {
				assert.Equal(t, filepath.Join(root, ".agentd.yaml"), got)
			}
		})
	}
}
