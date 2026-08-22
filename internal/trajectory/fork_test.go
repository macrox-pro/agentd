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

func TestForkSession(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	prov, sid := "claude-code", "src-sess"
	writeForkSource(t, root, prov, sid, 3)

	srcBefore, err := os.ReadFile(filepath.Join(root, prov, sid+".jsonl"))
	require.NoError(t, err)

	key := trajectory.ResolveSessionKey(prov, sid, "", "")
	result, err := trajectory.ForkSession(root, key, "forked-sess", 2)
	require.NoError(t, err, "ForkSession")
	assert.Equal(t, "forked-sess", result.NewSessionID)
	assert.Equal(t, sid, result.ParentSession)
	assert.Equal(t, uint64(2), result.BoundarySeq)
	assert.Equal(t, 2, result.Copied)

	srcAfter, err := os.ReadFile(filepath.Join(root, prov, sid+".jsonl"))
	require.NoError(t, err)
	assert.Equal(t, srcBefore, srcAfter, "source must be immutable")

	events, err := trajectory.ReadEvents(filepath.Join(root, prov, "forked-sess.jsonl"))
	require.NoError(t, err)
	require.Len(t, events, 4) // 2 copied + fork + end-seed
	assert.Equal(t, trajectory.TypeSessionFork, events[2].Type)
	assert.Equal(t, trajectory.TypeSessionEndSeed, events[3].Type)
	var forkData trajectory.SessionForkData
	require.NoError(t, json.Unmarshal(events[2].Data, &forkData))
	assert.Equal(t, sid, forkData.ParentSession)
	assert.Equal(t, uint64(2), forkData.BoundarySeq)
}

func TestForkDuplicateRejected(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	prov, sid := "cursor", "s1"
	writeForkSource(t, root, prov, sid, 2)
	key := trajectory.ResolveSessionKey(prov, sid, "", "")
	_, err := trajectory.ForkSession(root, key, "dup", 0)
	require.NoError(t, err)
	_, err = trajectory.ForkSession(root, key, "dup", 0)
	require.ErrorIs(t, err, trajectory.ErrSessionAlreadyExists)
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
