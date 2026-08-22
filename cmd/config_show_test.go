package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeConfigYAML(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestConfigShow(t *testing.T) {
	tests := []struct {
		name     string
		cfg      string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "merged and layer conflict",
			cfg:      "version: 1\n",
			args:     []string{"config", "show", "--merged", "--layer", "user"},
			wantErr:  true,
			contains: "not both",
		},
		{
			name:     "unknown layer",
			cfg:      "version: 1\n",
			args:     []string{"config", "show", "--layer", "bogus"},
			wantErr:  true,
			contains: "unknown layer",
		},
		{
			name:     "layer user happy",
			cfg:      "version: 1\n",
			args:     []string{"config", "show", "--layer", "user"},
			contains: "version",
		},
		{
			name:     "merged default happy",
			cfg:      "version: 1\n",
			args:     []string{"config", "show", "--merged"},
			contains: "version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := writeConfigYAML(t, tt.cfg)
			got := executeRoot(t, execOpts{args: tt.args, configPath: cfg})
			if tt.wantErr {
				require.Error(t, got.err)
				assert.Contains(t, got.err.Error(), tt.contains)
				return
			}
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
			assert.NotEmpty(t, got.out)
		})
	}
}
