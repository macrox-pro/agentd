package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name: "ok missing file",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "absent.yaml")
			},
			args:     []string{"config", "validate"},
			contains: "ok",
		},
		{
			name: "ok minimal yaml",
			setup: func(t *testing.T) string {
				path := filepath.Join(t.TempDir(), ".agentd.yaml")
				require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600))
				return path
			},
			args:     []string{"config", "validate"},
			contains: "ok",
		},
		{
			name: "invalid yaml",
			setup: func(t *testing.T) string {
				path := filepath.Join(t.TempDir(), ".agentd.yaml")
				require.NoError(t, os.WriteFile(path, []byte(":\n  bad:\n"), 0o600))
				return path
			},
			args:    []string{"config", "validate"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup(t)
			got := executeRoot(t, execOpts{args: tt.args, configPath: cfg})
			if tt.wantErr {
				require.Error(t, got.err)
				return
			}
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
		})
	}
}
