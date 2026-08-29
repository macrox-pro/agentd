package importer_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestImportSessionEmitOnly(t *testing.T) {
	t.Parallel()
	transcript := filepath.Join("testdata", "claude_session.jsonl")

	tests := []struct {
		name string
		run  func(t *testing.T) (importer.ImportSessionResult, string)
	}{
		{
			name: "emit_only_returns_events_no_ledger",
			run: func(t *testing.T) (importer.ImportSessionResult, string) {
				root := t.TempDir()
				sessionsRoot := filepath.Join(root, "sessions")
				sid := "emit-ledger"
				out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
					Provider:       provider.ClaudeCode,
					SessionID:      sid,
					TranscriptPath: transcript,
					SessionsRoot:   sessionsRoot,
					EmitOnly:       true,
				})
				require.NoError(t, err, "ImportSession emit_only_returns_events_no_ledger")
				ledger := trajectory.SessionFilePath(sessionsRoot, trajectory.SessionKey{
					Provider: provider.ClaudeCode, SessionID: sid,
				})
				_, statErr := os.Stat(ledger)
				assert.True(t, os.IsNotExist(statErr))
				return out, sessionsRoot
			},
		},
		{
			name: "emit_only_no_checkpoint",
			run: func(t *testing.T) (importer.ImportSessionResult, string) {
				root := t.TempDir()
				sessionsRoot := filepath.Join(root, "sessions")
				sid := "emit-cp"
				out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
					Provider:       provider.ClaudeCode,
					SessionID:      sid,
					TranscriptPath: transcript,
					SessionsRoot:   sessionsRoot,
					EmitOnly:       true,
				})
				require.NoError(t, err, "ImportSession emit_only_no_checkpoint")
				_, statErr := os.Stat(trajectory.ImportSidecarPath(sessionsRoot, "claude-code", sid))
				assert.True(t, os.IsNotExist(statErr))
				return out, sessionsRoot
			},
		},
		{
			name: "emit_only_no_hub_publish",
			run: func(t *testing.T) (importer.ImportSessionResult, string) {
				root := t.TempDir()
				sessionsRoot := filepath.Join(root, "sessions")
				hub := trajectory.NewHub(nil)
				ch, unregister := hub.Register(trajectory.SubscribeFilter{Provider: "claude-code"})
				defer unregister()
				defer hub.Close()

				out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
					Provider:       provider.ClaudeCode,
					SessionID:      "emit-hub",
					TranscriptPath: transcript,
					SessionsRoot:   sessionsRoot,
					EmitOnly:       true,
					Hub:            hub,
				})
				require.NoError(t, err, "ImportSession emit_only_no_hub_publish")

				select {
				case <-ch:
					t.Fatal("unexpected hub event")
				case <-time.After(50 * time.Millisecond):
				}
				return out, sessionsRoot
			},
		},
		{
			name: "emit_only_respects_start_index",
			run: func(t *testing.T) (importer.ImportSessionResult, string) {
				root := t.TempDir()
				sessionsRoot := filepath.Join(root, "sessions")
				sid := "emit-resume"

				first, err := importer.Import(provider.ClaudeCode, importer.ImportOptions{
					SessionID:      sid,
					TranscriptPath: transcript,
				})
				require.NoError(t, err, "Import first")
				require.NotEmpty(t, first.Events)
				mid := first.LastLineIndex / 2
				require.NoError(t, trajectory.SaveImportCheckpoint(
					trajectory.ImportSidecarPath(sessionsRoot, "claude-code", sid),
					trajectory.ImportCheckpoint{
						LastLineIndex: mid,
						SourcePath:    transcript,
					},
				))

				out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
					Provider:       provider.ClaudeCode,
					SessionID:      sid,
					TranscriptPath: transcript,
					SessionsRoot:   sessionsRoot,
					EmitOnly:       true,
				})
				require.NoError(t, err, "ImportSession emit_only_respects_start_index")
				assert.Less(t, out.Imported, len(first.Events))
				assert.Greater(t, out.Imported, 0)
				return out, sessionsRoot
			},
		},
		{
			name: "incremental_all_caught_up",
			run: func(t *testing.T) (importer.ImportSessionResult, string) {
				root := t.TempDir()
				sessionsRoot := filepath.Join(root, "sessions")
				sid := "emit-done"

				first, err := importer.Import(provider.ClaudeCode, importer.ImportOptions{
					SessionID:      sid,
					TranscriptPath: transcript,
				})
				require.NoError(t, err, "Import first")
				require.NoError(t, trajectory.SaveImportCheckpoint(
					trajectory.ImportSidecarPath(sessionsRoot, "claude-code", sid),
					trajectory.ImportCheckpoint{
						LastLineIndex: first.LastLineIndex,
						SourcePath:    transcript,
					},
				))

				out, err := importer.ImportSession(context.Background(), importer.ImportSessionOptions{
					Provider:       provider.ClaudeCode,
					SessionID:      sid,
					TranscriptPath: transcript,
					SessionsRoot:   sessionsRoot,
					EmitOnly:       true,
				})
				require.NoError(t, err, "ImportSession incremental_all_caught_up")
				assert.Equal(t, 0, out.Imported)
				assert.Empty(t, out.Events)
				return out, sessionsRoot
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out, _ := tt.run(t)
			if tt.name == "incremental_all_caught_up" {
				return
			}
			assert.NotEmpty(t, out.Events)
			assert.Equal(t, len(out.Events), out.Imported)
		})
	}
}

func TestImportSessionResultJSON(t *testing.T) {
	t.Parallel()

	t.Run("json_summary_excludes_events", func(t *testing.T) {
		t.Parallel()
		result := importer.ImportSessionResult{
			Provider:  "claude-code",
			SessionID: "s1",
			Imported:  1,
			Events: []trajectory.Event{{
				Type:   trajectory.TypeTranscriptMessage,
				Source: trajectory.SourceTranscript,
				Data:   json.RawMessage(`{"text":"secret-body"}`),
			}},
		}
		b, err := json.Marshal(result)
		require.NoError(t, err, "Marshal ImportSessionResult")
		s := string(b)
		assert.NotContains(t, s, "secret-body")
		assert.NotContains(t, s, `"events"`)
		assert.Contains(t, s, `"imported":1`)
	})
}
