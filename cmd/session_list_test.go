package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionListCLI(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, root string)
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "empty",
			args:     []string{"session", "list"},
			contains: "no sessions",
		},
		{
			name: "text rows",
			setup: func(t *testing.T, root string) {
				t.Helper()
				writeSessionLedger(t, root, "claude-code", "s1", 1)
			},
			args:     []string{"session", "list"},
			contains: "claude-code",
		},
		{
			name: "json",
			setup: func(t *testing.T, root string) {
				t.Helper()
				writeSessionLedger(t, root, "cursor", "c1", 1)
			},
			args:     []string{"session", "list", "--json"},
			contains: `"provider"`,
		},
		{
			name: "provider filter",
			setup: func(t *testing.T, root string) {
				t.Helper()
				writeSessionLedger(t, root, "claude-code", "s1", 1)
				writeSessionLedger(t, root, "cursor", "c1", 1)
			},
			args:     []string{"session", "list", "--provider", "cursor"},
			contains: "cursor",
		},
		{
			name:    "unknown provider",
			args:    []string{"session", "list", "--provider", "nope"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tempSessionsEnv(t)
			if tt.setup != nil {
				tt.setup(t, root)
			}
			got := executeRoot(t, execOpts{args: tt.args})
			if tt.wantErr {
				require.Error(t, got.err, "session list(%q)", tt.name)
				return
			}
			require.NoError(t, got.err, "session list(%q)", tt.name)
			assert.Contains(t, got.out, tt.contains, "session list(%q)", tt.name)
		})
	}
}
