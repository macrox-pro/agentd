package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHookCLIErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		contain string
	}{
		{
			name:    "notify missing payload",
			args:    []string{"hook", "notify", "--provider", "codex"},
			contain: "notify payload required",
		},
		{
			name:    "run argv missing payload",
			args:    []string{"hook", "run", "--provider", "cursor", "--argv-payload"},
			contain: "argv payload required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executeRoot(t, execOpts{args: tt.args})
			require.Error(t, got.err, "hook(%q)", tt.name)
			assert.Contains(t, got.err.Error(), tt.contain, "hook(%q)", tt.name)
		})
	}
}
