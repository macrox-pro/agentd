package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestOfflineFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T) (config.LoadOptions, string)
		want    config.FailMode
		wantErr bool
	}{
		{
			name: "default fail_open",
			setup: func(t *testing.T) (config.LoadOptions, string) {
				t.Helper()
				dir := t.TempDir()
				return config.LoadOptions{
					UserPath:    filepath.Join(dir, "missing.yaml"),
					RuntimePath: filepath.Join(dir, "missing-runtime.yaml"),
				}, ""
			},
			want: config.FailOpen,
		},
		{
			name: "user fail_closed",
			setup: func(t *testing.T) (config.LoadOptions, string) {
				t.Helper()
				dir := t.TempDir()
				user := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(user, []byte("version: 1\npolicy:\n  offline: fail_closed\n"), 0o600))
				return config.LoadOptions{
					UserPath:    user,
					RuntimePath: filepath.Join(dir, "missing-runtime.yaml"),
				}, ""
			},
			want: config.FailClosed,
		},
		{
			name: "project overrides",
			setup: func(t *testing.T) (config.LoadOptions, string) {
				t.Helper()
				dir := t.TempDir()
				proj := filepath.Join(dir, "proj")
				require.NoError(t, os.MkdirAll(proj, 0o700))
				require.NoError(t, os.WriteFile(filepath.Join(proj, ".agentd.yaml"), []byte("version: 1\npolicy:\n  offline: fail_closed\n"), 0o600))
				return config.LoadOptions{
					UserPath:    filepath.Join(dir, "missing.yaml"),
					RuntimePath: filepath.Join(dir, "missing-runtime.yaml"),
				}, proj
			},
			want: config.FailClosed,
		},
		{
			name: "runtime overrides",
			setup: func(t *testing.T) (config.LoadOptions, string) {
				t.Helper()
				dir := t.TempDir()
				user := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(user, []byte("version: 1\npolicy:\n  offline: fail_open\n"), 0o600))
				rt := filepath.Join(dir, "runtime.yaml")
				require.NoError(t, os.WriteFile(rt, []byte("version: 1\npolicy:\n  offline: fail_closed\n"), 0o600))
				return config.LoadOptions{UserPath: user, RuntimePath: rt}, ""
			},
			want: config.FailClosed,
		},
		{
			name: "missing user file",
			setup: func(t *testing.T) (config.LoadOptions, string) {
				t.Helper()
				dir := t.TempDir()
				return config.LoadOptions{
					UserPath:    filepath.Join(dir, "nope.yaml"),
					RuntimePath: filepath.Join(dir, "nope-rt.yaml"),
				}, ""
			},
			want: config.FailOpen,
		},
		{
			name: "invalid yaml",
			setup: func(t *testing.T) (config.LoadOptions, string) {
				t.Helper()
				dir := t.TempDir()
				user := filepath.Join(dir, "bad.yaml")
				require.NoError(t, os.WriteFile(user, []byte(":\n  - invalid\n"), 0o600))
				return config.LoadOptions{
					UserPath:    user,
					RuntimePath: filepath.Join(dir, "missing-runtime.yaml"),
				}, ""
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts, cwd := tt.setup(t)
			got, err := config.OfflineFor(opts, cwd)
			if tt.wantErr {
				require.Error(t, err, "OfflineFor(%q)", tt.name)
				return
			}
			require.NoError(t, err, "OfflineFor(%q)", tt.name)
			assert.Equal(t, tt.want, got, "OfflineFor(%q)", tt.name)
		})
	}
}
