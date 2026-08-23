package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigPatchCLI(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		setup   func(t *testing.T) []string
		contain string
	}{
		{
			name: "missing file flag",
			args: []string{"config", "patch"},
		},
		{
			name: "daemon not running",
			setup: func(t *testing.T) []string {
				t.Helper()
				path := filepath.Join(t.TempDir(), "patch.yaml")
				require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600))
				sock := filepath.Join(t.TempDir(), "missing.sock")
				return []string{"config", "patch", "--file", path, "--socket", sock}
			},
			contain: "Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args
			if tt.setup != nil {
				args = tt.setup(t)
			}
			got := executeRoot(t, execOpts{args: args})
			require.Error(t, got.err, "config patch(%q)", tt.name)
			if tt.contain != "" {
				assert.Contains(t, got.err.Error(), tt.contain, "config patch(%q)", tt.name)
			}
		})
	}
}
