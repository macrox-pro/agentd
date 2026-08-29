package cmd_test

import (
	"os"
	"path/filepath"
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

func TestSessionImportOut(t *testing.T) {
	transcript := claudeTranscriptFixture(t)

	t.Run("out_stdout", func(t *testing.T) {
		sessionsRoot := tempSessionsEnv(t)
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--session", "out-stdout",
			"--path", transcript,
			"--out", "-",
		}})
		require.NoError(t, got.err)
		assert.Contains(t, got.out, `"type":"transcript/`)
		_, err := os.Stat(filepath.Join(sessionsRoot, "claude-code", "out-stdout.jsonl"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("out_file", func(t *testing.T) {
		sessionsRoot := tempSessionsEnv(t)
		outFile := filepath.Join(t.TempDir(), "events.jsonl")
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--session", "out-file",
			"--path", transcript,
			"--out", outFile,
		}})
		require.NoError(t, got.err)
		b, err := os.ReadFile(outFile)
		require.NoError(t, err)
		assert.Contains(t, string(b), "transcript/")
		_, err = os.Stat(filepath.Join(sessionsRoot, "claude-code", "out-file.jsonl"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("dry_run_regression", func(t *testing.T) {
		tempSessionsEnv(t)
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--path", transcript,
			"--dry-run", "--json",
		}})
		require.NoError(t, got.err)
		assert.Contains(t, got.out, `"imported"`)
		assert.NotContains(t, got.out, `"type":"transcript/`)
	})

	t.Run("dry_run_out_combined", func(t *testing.T) {
		sessionsRoot := tempSessionsEnv(t)
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--session", "dry-out",
			"--path", transcript,
			"--dry-run", "--out", "-", "--json",
		}})
		require.NoError(t, got.err)
		assert.Contains(t, got.out, `"type":"transcript/`)
		assert.Contains(t, got.out, `"dry_run": true`)
		_, err := os.Stat(filepath.Join(sessionsRoot, "claude-code", "dry-out.jsonl"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("empty_transcript", func(t *testing.T) {
		sessionsRoot := tempSessionsEnv(t)
		empty := filepath.Join(t.TempDir(), "empty.jsonl")
		require.NoError(t, os.WriteFile(empty, nil, 0o600))
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--session", "empty-out",
			"--path", empty,
			"--out", "-",
		}})
		require.NoError(t, got.err)
		assert.Contains(t, got.out, "imported=0")
		assert.NotContains(t, got.out, `"type":"transcript/`)
		_, err := os.Stat(filepath.Join(sessionsRoot, "claude-code", "empty-out.jsonl"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("out_empty_invalid", func(t *testing.T) {
		tempSessionsEnv(t)
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--path", transcript,
			"--out=",
		}})
		require.Error(t, got.err)
		assert.Contains(t, got.err.Error(), "import --out requires a path or -")
	})

	t.Run("out_whitespace_invalid", func(t *testing.T) {
		tempSessionsEnv(t)
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--path", transcript,
			"--out", " ",
		}})
		require.Error(t, got.err)
		assert.Contains(t, got.err.Error(), "import --out requires a path or -")
	})

	t.Run("default_import_writes_ledger", func(t *testing.T) {
		sessionsRoot := tempSessionsEnv(t)
		got := executeRoot(t, execOpts{args: []string{
			"session", "import",
			"--provider", "claude-code",
			"--session", "default-ledger",
			"--path", transcript,
		}})
		require.NoError(t, got.err)
		assert.Contains(t, got.out, "imported=")
		_, err := os.Stat(filepath.Join(sessionsRoot, "claude-code", "default-ledger.jsonl"))
		require.NoError(t, err, "ledger should exist")
	})
}
