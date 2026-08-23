package cmd_test

import (
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
