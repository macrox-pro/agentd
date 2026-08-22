package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionShowCLI(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T)
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "unknown provider",
			args:     []string{"session", "show", "s1", "--provider", "nope"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name:     "session not found",
			setup:    func(t *testing.T) { tempSessionsEnv(t) },
			args:     []string{"session", "show", "missing", "--provider", "claude-code"},
			wantErr:  true,
			contains: "not found",
		},
		{
			name: "happy",
			setup: func(t *testing.T) {
				root := tempSessionsEnv(t)
				writeSessionLedger(t, root, "claude-code", "s1", 2)
			},
			args:     []string{"session", "show", "s1", "--provider", "claude-code"},
			contains: "hook/invoked",
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
