package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionReplayPolicy(t *testing.T) {
	cfgMissing := filepath.Join(t.TempDir(), "missing.yaml")

	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "policy flag required",
			setup:    func(t *testing.T) string { tempSessionsEnv(t); return cfgMissing },
			args:     []string{"session", "replay", "--provider", "claude-code", "--session", "s1"},
			wantErr:  true,
			contains: "--policy is required",
		},
		{
			name:     "missing session",
			setup:    func(t *testing.T) string { tempSessionsEnv(t); return cfgMissing },
			args:     []string{"session", "replay", "--policy", "--provider", "claude-code"},
			wantErr:  true,
			contains: "required flag",
		},
		{
			name:     "unknown provider",
			setup:    func(t *testing.T) string { tempSessionsEnv(t); return cfgMissing },
			args:     []string{"session", "replay", "--policy", "--provider", "nope", "--session", "s1"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name:     "session not found",
			setup:    func(t *testing.T) string { tempSessionsEnv(t); return cfgMissing },
			args:     []string{"session", "replay", "--policy", "--provider", "claude-code", "--session", "missing"},
			wantErr:  true,
			contains: "not found",
		},
		{
			name: "no raw at record time",
			setup: func(t *testing.T) string {
				root := tempSessionsEnv(t)
				writeSessionLedgerNoRaw(t, root, "claude-code", "no-raw")
				return cfgMissing
			},
			args:     []string{"session", "replay", "--policy", "--provider", "claude-code", "--session", "no-raw"},
			wantErr:  true,
			contains: "include_raw",
		},
		{
			name: "happy hits",
			setup: func(t *testing.T) string {
				root := tempSessionsEnv(t)
				raw := agenthookstest.Fixture(t, "claude/pre_tool_use.json")
				writeReplayLedger(t, root, "claude-code", "replay-ok", "stdin", raw)
				return cfgMissing
			},
			args:     []string{"session", "replay", "--policy", "--provider", "claude-code", "--session", "replay-ok"},
			contains: "hits=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup(t)
			got := executeRoot(t, execOpts{args: tt.args, configPath: cfg})
			if tt.wantErr {
				require.Error(t, got.err)
				assert.Contains(t, got.err.Error(), tt.contains)
				return
			}
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
		})
	}
}
