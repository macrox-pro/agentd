package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstallCLI(t *testing.T) {
	tests := []struct {
		name string
		args func(t *testing.T) []string
	}{
		{
			name: "provider required",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install"}
			},
		},
		{
			name: "unknown provider",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"install", "--provider", "nope", "--dir", t.TempDir()}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executeRoot(t, execOpts{args: tt.args(t)})
			require.Error(t, got.err, "install(%q)", tt.name)
		})
	}
}
