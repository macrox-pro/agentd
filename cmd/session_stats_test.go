package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStatsCLI(t *testing.T) {
	enabledYAML := "version: 1\ntrajectory:\n  enabled: true\n  statistics: true\n"
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "provider_required",
			setup:    func(t *testing.T) string { return writeStatsConfig(t, enabledYAML) },
			args:     []string{"session", "stats", "s1"},
			wantErr:  true,
			contains: "required",
		},
		{
			name:     "invalid_provider",
			setup:    func(t *testing.T) string { return writeStatsConfig(t, enabledYAML) },
			args:     []string{"session", "stats", "s1", "--provider", "nope"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name: "trajectory_disabled",
			setup: func(t *testing.T) string {
				return writeStatsConfig(t, "version: 1\ntrajectory:\n  enabled: false\n  statistics: true\n")
			},
			args:     []string{"session", "stats", "s1", "--provider", "claude-code"},
			wantErr:  true,
			contains: "trajectory is disabled",
		},
		{
			name: "statistics_disabled",
			setup: func(t *testing.T) string {
				return writeStatsConfig(t, "version: 1\ntrajectory:\n  enabled: true\n  statistics: false\n")
			},
			args:     []string{"session", "stats", "s1", "--provider", "claude-code"},
			wantErr:  true,
			contains: "trajectory statistics is disabled",
		},
		{
			name: "not_found",
			setup: func(t *testing.T) string {
				tempSessionsEnv(t)
				return writeStatsConfig(t, enabledYAML)
			},
			args:     []string{"session", "stats", "missing", "--provider", "claude-code"},
			wantErr:  true,
			contains: "not found",
		},
		{
			name: "human",
			setup: func(t *testing.T) string {
				root := tempSessionsEnv(t)
				writeSessionLedger(t, root, "claude-code", "s1", 2)
				return writeStatsConfig(t, enabledYAML)
			},
			args:     []string{"session", "stats", "s1", "--provider", "claude-code"},
			contains: "provider=claude-code",
		},
		{
			name: "json",
			setup: func(t *testing.T) string {
				root := tempSessionsEnv(t)
				writeSessionLedger(t, root, "claude-code", "s1", 1)
				return writeStatsConfig(t, enabledYAML)
			},
			args:     []string{"session", "stats", "s1", "--provider", "claude-code", "--json"},
			contains: `"session_id"`,
		},
		{
			name: "corrupt_jsonl",
			setup: func(t *testing.T) string {
				root := tempSessionsEnv(t)
				dir := filepath.Join(root, "claude-code")
				require.NoError(t, os.MkdirAll(dir, 0o700))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "bad.jsonl"), []byte("{not json\n"), 0o600))
				return writeStatsConfig(t, enabledYAML)
			},
			args:    []string{"session", "stats", "bad", "--provider", "claude-code"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup(t)
			got := executeRoot(t, execOpts{args: tt.args, configPath: cfg})
			if tt.wantErr {
				require.Error(t, got.err)
				if tt.contains != "" {
					assert.Contains(t, got.err.Error(), tt.contains)
				}
				return
			}
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
		})
	}
}
