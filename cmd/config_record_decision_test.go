package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigRecordDecisionCLI(t *testing.T) {
	sock := filepath.Join(t.TempDir(), "missing.sock")

	tests := []struct {
		name    string
		args    []string
		contain string
	}{
		{
			name: "missing fingerprint",
			args: []string{"config", "record-decision"},
		},
		{
			name:    "invalid scope",
			args:    []string{"config", "record-decision", "--fingerprint", "sha256:x", "--scope", "bogus", "--socket", sock},
			contain: "scope must be project or session",
		},
		{
			name:    "daemon unavailable project",
			args:    []string{"config", "record-decision", "--fingerprint", "sha256:x", "--scope", "project", "--project-root", "/tmp/p", "--socket", sock},
			contain: "Unavailable",
		},
		{
			name:    "daemon unavailable session",
			args:    []string{"config", "record-decision", "--fingerprint", "sha256:x", "--scope", "session", "--session-id", "s1", "--socket", sock},
			contain: "Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executeRoot(t, execOpts{args: tt.args})
			require.Error(t, got.err, "record-decision(%q)", tt.name)
			if tt.contain != "" {
				assert.Contains(t, got.err.Error(), tt.contain, "record-decision(%q)", tt.name)
			}
		})
	}
}
