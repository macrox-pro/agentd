package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDisable(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		preEnable  bool
		bootstrap  string
		wantSubstr []string
	}{
		{
			name:       "disable_idempotent_exit_0",
			args:       []string{"config", "disable", "trajectory"},
			wantSubstr: []string{"already disabled"},
			bootstrap:  "version: 1\ntrajectory:\n  enabled: false\n",
		},
		{
			name:       "disable_trajectory_writes_false",
			args:       []string{"config", "disable", "trajectory"},
			preEnable:  true,
			wantSubstr: []string{"trajectory: disabled"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			configPath := filepath.Join(dir, "user.yaml")
			switch {
			case tt.preEnable:
				require.NoError(t, os.WriteFile(configPath, []byte("version: 1\npolicy:\n  fail: fail_closed\ntrajectory:\n  enabled: true\n"), 0o600))
			case tt.bootstrap != "":
				require.NoError(t, os.WriteFile(configPath, []byte(tt.bootstrap), 0o600))
			}
			res := executeRoot(t, execOpts{args: tt.args, configPath: configPath})
			require.NoError(t, res.err, "Execute(%q)", tt.name)
			for _, sub := range tt.wantSubstr {
				assert.Contains(t, res.out, sub, "Execute(%q) output", tt.name)
			}
			if tt.preEnable {
				raw, err := os.ReadFile(configPath)
				require.NoError(t, err, "ReadFile(%q)", configPath)
				assert.Contains(t, string(raw), "enabled: false", "explicit false")
			}
		})
	}
}
