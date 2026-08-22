package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestDispatchRoutes(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (cfgPath string, wantErr bool)
		args     []string
		contains string
		jsonOK   bool
	}{
		{
			name: "human",
			setup: func(t *testing.T) (string, bool) {
				return filepath.Join(t.TempDir(), ".agentd.yaml"), false
			},
			args:     []string{"dispatch", "routes"},
			contains: "mode=",
		},
		{
			name: "json",
			setup: func(t *testing.T) (string, bool) {
				return filepath.Join(t.TempDir(), ".agentd.yaml"), false
			},
			args:   []string{"dispatch", "routes", "--json"},
			jsonOK: true,
		},
		{
			name: "invalid config",
			setup: func(t *testing.T) (string, bool) {
				path := filepath.Join(t.TempDir(), ".agentd.yaml")
				require.NoError(t, os.WriteFile(path, []byte(":\n  bad:\n"), 0o600))
				return path, true
			},
			args: []string{"dispatch", "routes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, wantErr := tt.setup(t)
			got := executeRoot(t, execOpts{args: tt.args, configPath: cfg})
			if wantErr {
				require.Error(t, got.err)
				return
			}
			require.NoError(t, got.err, "dispatch routes")
			if tt.jsonOK {
				var routes []config.CompiledRoute
				require.NoError(t, json.Unmarshal([]byte(got.out), &routes), "dispatch routes json")
				require.NotEmpty(t, routes)
				assert.NotEmpty(t, routes[0].Name)
				return
			}
			assert.Contains(t, got.out, tt.contains)
		})
	}
}
