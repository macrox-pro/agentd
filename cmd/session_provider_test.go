package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionProviderCLI(t *testing.T) {
	tempSessionsEnv(t)

	tests := []struct {
		name     string
		setup    func(t *testing.T)
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name: "list all providers",
			args: []string{"session", "list"},
		},
		{
			name:     "list unknown provider",
			args:     []string{"session", "list", "--provider", "nope"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name: "list json importer_status",
			setup: func(t *testing.T) {
				root := tempSessionsEnv(t)
				writeSessionLedger(t, root, "claude-code", "s1", 2)
			},
			args:     []string{"session", "list", "--json"},
			contains: `"importer_status": "supported"`,
		},
		{
			name:     "import cursor without path",
			args:     []string{"session", "import", "--provider", "cursor", "--session", "x"},
			wantErr:  true,
			contains: "cursor import requires --path",
		},
		{
			name:     "import claude without session or path",
			args:     []string{"session", "import", "--provider", "claude-code"},
			wantErr:  true,
			contains: "import requires --session or --path",
		},
		{
			name:     "import unsupported provider",
			args:     []string{"session", "import", "--provider", "gemini", "--session", "x"},
			wantErr:  true,
			contains: "not supported",
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
				require.Error(t, got.err, "args=%v", tt.args)
				if tt.contains != "" {
					assert.Contains(t, got.err.Error(), tt.contains)
				}
				return
			}
			require.NoError(t, got.err, "args=%v", tt.args)
			if tt.name == "list all providers" {
				assert.Contains(t, got.out, "no sessions")
			}
			if tt.contains != "" {
				assert.Contains(t, got.out, tt.contains)
			}
		})
	}
}
