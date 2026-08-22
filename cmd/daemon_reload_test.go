package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDaemonReloadCLI(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string
		wantErr     bool
		errContains string
		contains    string
	}{
		{
			name:        "not running",
			setup:       func(t *testing.T) string { socket, _ := testSocketDir(t); return socket },
			wantErr:     true,
			errContains: "Unavailable",
		},
		{
			name:     "happy",
			setup:    startReloadStubServer,
			contains: "generation=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			socket := tt.setup(t)
			got := executeRoot(t, execOpts{args: []string{"daemon", "reload"}, socketPath: socket})
			if tt.wantErr {
				require.Error(t, got.err)
				assert.Contains(t, got.err.Error(), tt.errContains)
				return
			}
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
		})
	}
}
