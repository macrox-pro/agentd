package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigEnable(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		configPath string
		workDir    string
		wantErr    bool
		wantSubstr []string
	}{
		{
			name:    "unknown_feature_exit",
			args:    []string{"config", "enable", "not-real"},
			wantErr: true,
			wantSubstr: []string{"unknown feature"},
		},
		{
			name:    "unknown_feature_case_sensitive",
			args:    []string{"config", "enable", "Trajectory"},
			wantErr: true,
			wantSubstr: []string{"unknown feature"},
		},
		{
			name:    "missing_feature_arg",
			args:    []string{"config", "enable"},
			wantErr: true,
		},
		{
			name:    "extra_feature_arg",
			args:    []string{"config", "enable", "trajectory", "extra"},
			wantErr: true,
		},
		{
			name:    "invalid_scope_flag",
			args:    []string{"config", "enable", "trajectory", "--scope", "global"},
			wantErr: true,
			wantSubstr: []string{"unknown scope"},
		},
		{
			name: "enable_trajectory_happy_path",
			args: []string{"config", "enable", "trajectory"},
			wantSubstr: []string{"trajectory: already enabled"},
		},
		{
			name:    "enable_guard_shell_project_cwd",
			args:    []string{"config", "enable", "guard-shell", "--scope", "project"},
			wantSubstr: []string{"guard-shell: enabled", "project"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			configPath := tt.configPath
			if configPath == "" {
				configPath = filepath.Join(dir, "user.yaml")
			}
			workDir := tt.workDir
			if workDir == "" {
				workDir = dir
			}
			if tt.name == "enable_guard_shell_project_cwd" {
				workDir = dir
			}
			cwd, err := os.Getwd()
			require.NoError(t, err, "Getwd()")
			require.NoError(t, os.Chdir(workDir), "Chdir(%q)", workDir)
			t.Cleanup(func() { _ = os.Chdir(cwd) })

			res := executeRoot(t, execOpts{args: tt.args, configPath: configPath})
			if tt.wantErr {
				require.Error(t, res.err, "Execute(%q)", tt.name)
			} else {
				require.NoError(t, res.err, "Execute(%q)", tt.name)
				valRes := executeRoot(t, execOpts{args: []string{"config", "validate"}, configPath: configPath})
				require.NoError(t, valRes.err, "config validate after enable")
			}
			for _, sub := range tt.wantSubstr {
				assert.Contains(t, res.out, sub, "Execute(%q) output", tt.name)
			}
			if tt.name == "enable_guard_shell_project_cwd" {
				_, statErr := os.Stat(filepath.Join(workDir, ".agentd.yaml"))
				require.NoError(t, statErr, "Stat project config")
			}
			if tt.name == "enable_trajectory_happy_path" {
				_, statErr := os.Stat(configPath)
				require.Error(t, statErr, "user config not written when already at default")
			}
		})
	}
}
