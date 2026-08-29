package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/version"
)

func TestDaemonStatusCLI(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string
		args        []string
		cliVersion  string
		contains    string
		notContains string
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
			contains: `"autostart"`,
		},
		{
			name:     "running human generation",
			setup:    statusStubServer,
			args:     []string{"daemon", "status"},
			contains: "generation 1",
		},
		{
			name:     "running human daemon version",
			setup:    statusStubServer,
			args:     []string{"daemon", "status"},
			contains: "version test",
		},
		{
			name:     "running json daemon version",
			setup:    statusStubServer,
			args:     []string{"daemon", "status", "--json"},
			contains: `"version":"test"`,
		},
		{
			name:        "running uses daemon version not cli",
			setup:       statusStubServer,
			args:        []string{"daemon", "status"},
			cliVersion:  "cli-bin",
			contains:    "version test",
			notContains: "cli-bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cliVersion != "" {
				prev := version.Version
				version.Version = tt.cliVersion
				t.Cleanup(func() { version.Version = prev })
			}
			socket := tt.setup(t)
			got := executeRoot(t, execOpts{args: tt.args, socketPath: socket})
			require.NoError(t, got.err)
			assert.Contains(t, got.out, tt.contains)
			if tt.notContains != "" {
				assert.NotContains(t, got.out, tt.notContains)
			}
		})
	}
}
