package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionSearchCLI(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T)
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "unknown provider filter",
			setup:    func(t *testing.T) { tempSessionsEnv(t) },
			args:     []string{"session", "search", "--provider", "nope"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name:     "no matches",
			setup:    func(t *testing.T) { tempSessionsEnv(t) },
			args:     []string{"session", "search", "--query", "__nope__"},
			contains: "no matches",
		},
		{
			name: "happy",
			setup: func(t *testing.T) {
				root := tempSessionsEnv(t)
				writeSessionLedger(t, root, "claude-code", "s1", 2)
			},
			args:     []string{"session", "search"},
			contains: "claude-code/s1",
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
