package importer_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestImportSessionFacade(t *testing.T) {
	t.Parallel()
	transcript := filepath.Join("testdata", "claude_session.jsonl")

	t.Run("dry-run", func(t *testing.T) {
		root := t.TempDir()

		out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
			Provider:       provider.ClaudeCode,
			TranscriptPath: transcript,
			DryRun:         true,
			SessionsRoot:   filepath.Join(root, "sessions"),
		})
		require.NoError(t, err, "ImportSession dry-run")
		assert.Equal(t, "claude-code", out.Provider)
		assert.Equal(t, string(importer.ImporterSupported), out.ImporterStatus)
		assert.True(t, out.DryRun)
		assert.Greater(t, out.Imported, 0)
	})

	t.Run("write and checkpoint", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		sessionsRoot := filepath.Join(root, "sessions")
		sid := "facade-write"

		out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
			Provider:       provider.ClaudeCode,
			SessionID:      sid,
			TranscriptPath: transcript,
			SessionsRoot:   sessionsRoot,
		})
		require.NoError(t, err, "ImportSession write")
		assert.False(t, out.DryRun)
		assert.Greater(t, out.Imported, 0)

		ledger := trajectory.SessionFilePath(sessionsRoot, trajectory.SessionKey{
			Provider:  provider.ClaudeCode,
			SessionID: sid,
		})
		_, err = os.Stat(ledger)
		require.NoError(t, err, "ledger file")

		cp, err := trajectory.LoadImportCheckpoint(trajectory.ImportSidecarPath(sessionsRoot, "claude-code", sid))
		require.NoError(t, err, "LoadImportCheckpoint")
		assert.Equal(t, transcript, cp.SourcePath)
	})

	t.Run("hub publish", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		sessionsRoot := filepath.Join(root, "sessions")
		hub := trajectory.NewHub(nil)
		ch, unregister := hub.Register(trajectory.SubscribeFilter{Provider: "claude-code"})
		defer unregister()
		defer hub.Close()

		_, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
			Provider:       provider.ClaudeCode,
			SessionID:      "facade-hub",
			TranscriptPath: transcript,
			SessionsRoot:   sessionsRoot,
			Hub:            hub,
		})
		require.NoError(t, err, "ImportSession hub")

		select {
		case ev := <-ch:
			assert.Equal(t, "claude-code", ev.Provider)
			assert.Equal(t, "facade-hub", ev.SessionID)
		case <-time.After(time.Second):
			t.Fatal("expected hub event")
		}
	})

	t.Run("empty events skip", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		sessionsRoot := filepath.Join(root, "sessions")
		emptyPath := filepath.Join(root, "empty.jsonl")
		require.NoError(t, os.WriteFile(emptyPath, nil, 0o600))

		hub := trajectory.NewHub(nil)
		ch, unregister := hub.Register(trajectory.SubscribeFilter{})
		defer unregister()
		defer hub.Close()

		out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
			Provider:       provider.ClaudeCode,
			SessionID:      "empty-skip",
			TranscriptPath: emptyPath,
			SessionsRoot:   sessionsRoot,
			Hub:            hub,
		})
		require.NoError(t, err, "ImportSession empty")
		assert.Equal(t, 0, out.Imported)

		_, err = os.Stat(trajectory.ImportSidecarPath(sessionsRoot, "claude-code", "empty-skip"))
		assert.True(t, os.IsNotExist(err))

		select {
		case <-ch:
			t.Fatal("unexpected hub event")
		case <-time.After(50 * time.Millisecond):
		}
	})

	t.Run("snapshot cfg", func(t *testing.T) {
		t.Parallel()
		cfg := config.TrajectoryConfig{MaxEventBytes: 262144, RedactSecretRules: true}
		out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
			Provider:       provider.ClaudeCode,
			TranscriptPath: transcript,
			DryRun:         true,
			SnapshotCfg:    &cfg,
		})
		require.NoError(t, err, "ImportSession snapshot cfg")
		assert.Greater(t, out.Imported, 0)
	})
}
