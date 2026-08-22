package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionForkCLI(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T)
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "unknown provider",
			setup:    func(t *testing.T) { tempSessionsEnv(t) },
			args:     []string{"session", "fork", "--provider", "nope", "--session", "s1", "--new-session", "s2"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name:     "source not found",
			setup:    func(t *testing.T) { tempSessionsEnv(t) },
			args:     []string{"session", "fork", "--provider", "claude-code", "--session", "missing", "--new-session", "forked"},
			wantErr:  true,
			contains: "not found",
		},
		{
			name: "duplicate",
			setup: func(t *testing.T) {
				root := tempSessionsEnv(t)
				writeSessionLedger(t, root, "claude-code", "src", 2)
			},
			args: []string{"session", "fork", "--provider", "claude-code", "--session", "src", "--new-session", "dup"},
		},
		{
			name: "happy",
			setup: func(t *testing.T) {
				root := tempSessionsEnv(t)
				writeSessionLedger(t, root, "claude-code", "src", 2)
			},
			args:     []string{"session", "fork", "--provider", "claude-code", "--session", "src", "--new-session", "forked"},
			contains: "copied=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			} else {
				tempSessionsEnv(t)
			}
			got := executeRoot(t, execOpts{args: tt.args})
			if tt.name == "duplicate" {
				require.NoError(t, got.err)
				got2 := executeRoot(t, execOpts{args: tt.args})
				require.Error(t, got2.err)
				assert.Contains(t, got2.err.Error(), "already exists")
				return
			}
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
