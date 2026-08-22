package trajectory_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestReplayPolicyFromConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(t *testing.T) (cfgPath, sessionsRoot, provider, sessionID string, seq uint64)
		wantErr    error
		wantErrMsg string
		wantHits   bool
	}{
		{
			name: "bad_config_path",
			setup: func(t *testing.T) (string, string, string, string, uint64) {
				dir := t.TempDir()
				cfg := filepath.Join(dir, "bad.yaml")
				require.NoError(t, os.WriteFile(cfg, []byte(":\n  invalid: ["), 0o600))
				return cfg, t.TempDir(), "claude-code", "s1", 0
			},
			wantErrMsg: "load config",
		},
		{
			name: "missing_session",
			setup: func(t *testing.T) (string, string, string, string, uint64) {
				return filepath.Join(t.TempDir(), "missing.yaml"), t.TempDir(), "claude-code", "missing", 0
			},
			wantErr: trajectory.ErrSessionNotFound,
		},
		{
			name: "seq_not_found",
			setup: func(t *testing.T) (string, string, string, string, uint64) {
				root := t.TempDir()
				raw := agenthookstest.Fixture(t, "claude/pre_tool_use.json")
				sid := "replay-seq"
				writeReplayLedger(t, root, "claude-code", sid, "stdin", raw)
				return filepath.Join(t.TempDir(), "missing.yaml"), root, "claude-code", sid, 99
			},
			wantErr: trajectory.ErrReplaySeqNotFound,
		},
		{
			name: "ok",
			setup: func(t *testing.T) (string, string, string, string, uint64) {
				root := t.TempDir()
				raw := agenthookstest.Fixture(t, "claude/pre_tool_use.json")
				sid := "replay-ok"
				writeReplayLedger(t, root, "claude-code", sid, "stdin", raw)
				return filepath.Join(t.TempDir(), "missing.yaml"), root, "claude-code", sid, 0
			},
			wantHits: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfgPath, root, prov, sid, seq := tt.setup(t)
			result, err := trajectory.ReplayPolicyFromConfig(context.Background(), trajectory.ReplayPolicyConfigOptions{
				ConfigPath:   cfgPath,
				SessionsRoot: root,
				Provider:     prov,
				SessionID:    sid,
				Seq:          seq,
			})
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			if tt.wantErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
				return
			}
			require.NoError(t, err)
			if tt.wantHits {
				assert.NotEmpty(t, result.Hits)
			}
		})
	}
}
