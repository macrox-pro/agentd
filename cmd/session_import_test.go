package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionImportCLI(t *testing.T) {
	tempSessionsEnv(t)

	got := executeRoot(t, execOpts{
		args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--path", claudeTranscriptFixture(t),
			"--dry-run", "--json",
		},
	})
	require.NoError(t, got.err, "ImportSession dry-run")
	assert.Contains(t, got.out, `"importer_status"`)
	assert.Contains(t, got.out, `"imported"`)
}
