package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallCLI(t *testing.T) {
	tests := []struct {
		name    string
		args    func(t *testing.T) []string
		wantErr bool
		contain []string
	}{
		{
			name: "provider required",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install"}
			},
			wantErr: true,
			contain: []string{"--provider", "--all-detected"},
		},
		{
			name: "unknown provider",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install", "--provider", "nope", "--dir", t.TempDir()}
			},
			wantErr: true,
		},
		{
			name: "user default under HOME",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				return []string{"install", "--provider", "cursor", "--scope", "user"}
			},
			contain: []string{
				"provider=cursor",
				"scope=user",
				"create",
				"hooks.json",
			},
		},
		{
			name: "global same as scope user",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				return []string{"install", "--provider", "cursor", "--global"}
			},
			contain: []string{
				"provider=cursor",
				"scope=user",
				"create",
				"hooks.json",
			},
		},
		{
			name: "global conflicts with scope project",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install", "--provider", "cursor", "--global", "--scope", "project"}
			},
			wantErr: true,
			contain: []string{"--global", "--scope"},
		},
		{
			name: "plugin without dir errors",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install", "--provider", "cursor", "--scope", "plugin"}
			},
			wantErr: true,
			contain: []string{"--dir"},
		},
		{
			name: "all_detected_plan_only",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				require.NoError(t, os.MkdirAll(filepath.Join(home, ".cursor"), 0o755))
				return []string{"install", "--all-detected"}
			},
			contain: []string{"provider=cursor", "hooks=missing"},
		},
		{
			name: "all_detected_yes_writes",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				require.NoError(t, os.MkdirAll(filepath.Join(home, ".cursor"), 0o755))
				return []string{"install", "--all-detected", "--yes"}
			},
			contain: []string{"provider=cursor", "create", "hooks.json"},
		},
		{
			name: "all_detected_empty",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				return []string{"install", "--all-detected"}
			},
			contain: []string{"doctor"},
		},
		{
			name: "provider_with_dry_run",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				return []string{"install", "--provider", "cursor", "--scope", "user", "--dry-run"}
			},
			contain: []string{"provider=cursor", "create"},
		},
		{
			name: "provider_with_all_detected_error",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install", "--provider", "cursor", "--all-detected"}
			},
			wantErr: true,
			contain: []string{"mutually exclusive"},
		},
		{
			name: "yes_without_all_detected_error",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install", "--yes"}
			},
			wantErr: true,
			contain: []string{"--yes requires --all-detected"},
		},
		{
			name: "dry_run_with_yes_error",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install", "--all-detected", "--yes", "--dry-run"}
			},
			wantErr: true,
			contain: []string{"mutually exclusive"},
		},
		{
			name: "install_bare_non_tty_error",
			args: func(t *testing.T) []string {
				t.Helper()
				t.Setenv("AGENTD_NO_TUI", "1")
				return []string{"install"}
			},
			wantErr: true,
			contain: []string{"--provider", "--all-detected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executeRoot(t, execOpts{args: tt.args(t)})
			if tt.wantErr {
				require.Error(t, got.err)
				for _, s := range tt.contain {
					assert.Contains(t, got.out+got.err.Error(), s)
				}
				return
			}
			require.NoError(t, got.err)
			for _, s := range tt.contain {
				assert.Contains(t, got.out, s)
			}
		})
	}
}
