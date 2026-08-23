package trajectory_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestForkSessionTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		provider  string
		sessionID string
		events    int
		newID     string
		atSeq     uint64
		wantErr   error
	}{
		{
			name:      "ok",
			provider:  "claude-code",
			sessionID: "src-sess",
			events:    3,
			newID:     "forked-sess",
			atSeq:     2,
		},
		{
			name:      "duplicate rejected",
			provider:  "cursor",
			sessionID: "s1",
			events:    2,
			newID:     "dup",
			atSeq:     0,
			wantErr:   trajectory.ErrSessionAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			writeForkSource(t, root, tt.provider, tt.sessionID, tt.events)

			key := trajectory.ResolveSessionKey(tt.provider, tt.sessionID, "", "")

			if tt.wantErr != nil {
				_, err := trajectory.ForkSession(root, key, tt.newID, tt.atSeq)
				require.NoError(t, err, "ForkSession(%q)", tt.name)
				_, err = trajectory.ForkSession(root, key, tt.newID, tt.atSeq)
				require.ErrorIs(t, err, tt.wantErr, "ForkSession(%q)", tt.name)
				return
			}

			srcBefore, err := os.ReadFile(filepath.Join(root, tt.provider, tt.sessionID+".jsonl"))
			require.NoError(t, err, "ReadFile(%q)", tt.name)

			result, err := trajectory.ForkSession(root, key, tt.newID, tt.atSeq)
			require.NoError(t, err, "ForkSession(%q)", tt.name)
			assert.Equal(t, tt.newID, result.NewSessionID, "ForkSession(%q)", tt.name)
			assert.Equal(t, tt.sessionID, result.ParentSession, "ForkSession(%q)", tt.name)
			assert.Equal(t, tt.atSeq, result.BoundarySeq, "ForkSession(%q)", tt.name)
			assert.Equal(t, 2, result.Copied, "ForkSession(%q)", tt.name)

			srcAfter, err := os.ReadFile(filepath.Join(root, tt.provider, tt.sessionID+".jsonl"))
			require.NoError(t, err, "ReadFile(%q)", tt.name)
			assert.Equal(t, srcBefore, srcAfter, "ForkSession(%q)", tt.name)

			events, err := trajectory.ReadEvents(filepath.Join(root, tt.provider, tt.newID+".jsonl"))
			require.NoError(t, err, "ReadEvents(%q)", tt.name)
			require.Len(t, events, 4, "ForkSession(%q)", tt.name)
			assert.Equal(t, trajectory.TypeSessionFork, events[2].Type, "ForkSession(%q)", tt.name)
			assert.Equal(t, trajectory.TypeSessionEndSeed, events[3].Type, "ForkSession(%q)", tt.name)
			var forkData trajectory.SessionForkData
			require.NoError(t, json.Unmarshal(events[2].Data, &forkData), "ForkSession(%q)", tt.name)
			assert.Equal(t, tt.sessionID, forkData.ParentSession, "ForkSession(%q)", tt.name)
			assert.Equal(t, tt.atSeq, forkData.BoundarySeq, "ForkSession(%q)", tt.name)
		})
	}
}

func writeForkSource(t *testing.T, root, provider, sessionID string, n int) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	path := filepath.Join(dir, sessionID+".jsonl")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()
	enc := json.NewEncoder(f)
	now := time.Now().UTC()
	for i := 1; i <= n; i++ {
		require.NoError(t, enc.Encode(trajectory.Event{
			Seq:       uint64(i),
			Type:      trajectory.TypeHookInvoked,
			Source:    trajectory.SourceHook,
			TS:        now,
			Provider:  provider,
			SessionID: sessionID,
			Data:      json.RawMessage(`{"kind":"tool.pre"}`),
		}))
	}
}
