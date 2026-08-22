package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDaemonStatusCLI(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		args     []string
		contains string
	}{
		{
			name:     "not running human",
			setup:    func(t *testing.T) string { socket, _ := testSocketDir(t); return socket },
			args:     []string{"daemon", "status"},
			contains: "not running",
		},
		{
			name:     "not running json",
			setup:    func(t *testing.T) string { socket, _ := testSocketDir(t); return socket },
			args:     []string{"daemon", "status", "--json"},
			contains: `"running":false`,
		},
		{
			name:     "running",
			setup:    statusStubServer,
			args:     []string{"daemon", "status"},
			contains: "generation 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			socket := tt.setup(t)
			got := executeRoot(t, execOpts{args: tt.args, socketPath: socket})
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
		})
	}
}
