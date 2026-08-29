package trajectory_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestAssignImportedSeq(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T, root string, key trajectory.SessionKey)
		events  []trajectory.Event
		wantSeq []uint64
		wantErr bool
	}{
		{
			name: "empty_ledger",
			events: []trajectory.Event{{
				Type:   trajectory.TypeTranscriptMessage,
				Source: trajectory.SourceTranscript,
				TS:     time.Now().UTC(),
			}, {
				Type:   trajectory.TypeTranscriptThinking,
				Source: trajectory.SourceTranscript,
				TS:     time.Now().UTC(),
			}},
			wantSeq: []uint64{1, 2},
		},
		{
			name: "after_existing",
			setup: func(t *testing.T, root string, key trajectory.SessionKey) {
				t.Helper()
				path := trajectory.SessionFilePath(root, key)
				require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
				hookLine := `{"seq":3,"type":"hook/invoked","source":"hook","session_id":"s1","provider":"claude-code"}`
				require.NoError(t, os.WriteFile(path, []byte(hookLine+"\n"), 0o600))
			},
			events: []trajectory.Event{{
				Type:   trajectory.TypeTranscriptMessage,
				Source: trajectory.SourceTranscript,
				TS:     time.Now().UTC(),
			}},
			wantSeq: []uint64{4},
		},
		{
			name: "corrupt_ledger_read",
			setup: func(t *testing.T, root string, key trajectory.SessionKey) {
				t.Helper()
				path := trajectory.SessionFilePath(root, key)
				require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
				require.NoError(t, os.WriteFile(path, []byte("not-json\n"), 0o600))
			},
			events: []trajectory.Event{{
				Type:   trajectory.TypeTranscriptMessage,
				Source: trajectory.SourceTranscript,
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			key := trajectory.SessionKey{Provider: provider.ClaudeCode, SessionID: "s1"}
			if tt.setup != nil {
				tt.setup(t, root, key)
			}
			events := append([]trajectory.Event(nil), tt.events...)
			err := trajectory.AssignImportedSeq(root, key, events)
			if tt.wantErr {
				require.Error(t, err, "AssignImportedSeq(%q)", tt.name)
				return
			}
			require.NoError(t, err, "AssignImportedSeq(%q)", tt.name)
			require.Len(t, events, len(tt.wantSeq))
			for i, want := range tt.wantSeq {
				assert.Equal(t, want, events[i].Seq, "seq[%d]", i)
				assert.Equal(t, "claude-code", events[i].Provider)
				assert.Equal(t, "s1", events[i].SessionID)
			}
		})
	}
}

func TestWriteEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		events     []trajectory.Event
		wantLines  int
		wantSchema bool
	}{
		{
			name:      "empty",
			events:    nil,
			wantLines: 0,
		},
		{
			name: "single_event",
			events: []trajectory.Event{{
				Seq:       1,
				Type:      trajectory.TypeTranscriptMessage,
				Source:    trajectory.SourceTranscript,
				TS:        time.Now().UTC(),
				Provider:  "claude-code",
				SessionID: "s1",
				Data:      mustJSON(trajectory.TranscriptMessageData{Text: "hi", TranscriptLineIndex: 0}),
			}},
			wantLines:  1,
			wantSchema: true,
		},
		{
			name: "schema_version_stamped",
			events: []trajectory.Event{{
				Seq:       2,
				Type:      trajectory.TypeTranscriptThinking,
				Source:    trajectory.SourceTranscript,
				TS:        time.Now().UTC(),
				Provider:  "claude-code",
				SessionID: "s1",
			}},
			wantLines:  1,
			wantSchema: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			require.NoError(t, trajectory.WriteEvents(&buf, tt.events), "WriteEvents(%q)", tt.name)
			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			if tt.wantLines == 0 {
				assert.Empty(t, buf.String())
				return
			}
			require.Len(t, lines, tt.wantLines)
			if tt.wantSchema {
				assert.Contains(t, lines[0], `"schema_version":1`)
			}
		})
	}
}

func TestWriteEventsToFile(t *testing.T) {
	t.Parallel()

	t.Run("writes_file", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		path := filepath.Join(root, "out.jsonl")
		events := []trajectory.Event{{
			Seq:       1,
			Type:      trajectory.TypeTranscriptMessage,
			Source:    trajectory.SourceTranscript,
			TS:        time.Now().UTC(),
			Provider:  "claude-code",
			SessionID: "s1",
		}}
		require.NoError(t, trajectory.WriteEventsToFile(path, events))
		b, err := os.ReadFile(path)
		require.NoError(t, err, "ReadFile")
		assert.Contains(t, string(b), "transcript/message")
	})

	t.Run("permission_denied", func(t *testing.T) {
		t.Parallel()
		if os.Getuid() == 0 {
			t.Skip("root can write anywhere")
		}
		root := t.TempDir()
		readOnly := filepath.Join(root, "ro")
		require.NoError(t, os.MkdirAll(readOnly, 0o500))
		path := filepath.Join(readOnly, "out.jsonl")
		err := trajectory.WriteEventsToFile(path, []trajectory.Event{{Type: trajectory.TypeTranscriptMessage}})
		require.Error(t, err, "WriteEventsToFile permission_denied")
	})
}
