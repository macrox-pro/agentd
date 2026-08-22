package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDaemonStopCLI(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid timeout",
			args:        []string{"daemon", "stop", "--timeout", "nope"},
			wantErr:     true,
			errContains: "invalid --timeout",
		},
		{
			name:    "not running",
			args:    []string{"daemon", "stop"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			socket, _ := testSocketDir(t)
			got := executeRoot(t, execOpts{args: tt.args, socketPath: socket})
			if tt.wantErr {
				require.Error(t, got.err)
				assert.Contains(t, got.err.Error(), tt.errContains)
				return
			}
			require.NoError(t, got.err)
		})
	}
}
