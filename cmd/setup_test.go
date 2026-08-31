package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install/tui"
)

func TestSetupCLI(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "setup_non_tty_error",
			env:  map[string]string{"AGENTD_NO_TUI": "1"},
		},
		{
			name: "setup_ci_error",
			env:  map[string]string{"CI": "true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AGENTD_NO_TUI", "")
			t.Setenv("CI", "")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got := executeRoot(t, execOpts{args: []string{"setup"}})
			require.ErrorIs(t, got.err, tui.ErrNonInteractive, "setup(%q)", tt.name)
			combined := got.out
			if got.err != nil {
				combined += got.err.Error()
			}
			assert.Contains(t, combined, "--provider", "setup(%q)", tt.name)
			assert.Contains(t, combined, "--all-detected", "setup(%q)", tt.name)
		})
	}
}
