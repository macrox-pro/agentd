package trajectory_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestPersisterConcurrentSchedule(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	persist := trajectory.NewPersister(root, nil)

	providers := []string{"claude-code", "cursor", "codex", "gemini", "kimi-code", "opencode"}
	ev := []trajectory.Event{{
		Type:   trajectory.TypeHookInvoked,
		Source: trajectory.SourceHook,
		TS:     time.Now().UTC(),
	}}

	var wg sync.WaitGroup
	for range 50 {
		for _, prov := range providers {
			wg.Add(1)
			go func(provider string) {
				defer wg.Done()
				key := trajectory.SessionKey{Provider: provider, SessionID: "s-" + provider}
				persist.Schedule(key, ev)
			}(prov)
		}
	}
	wg.Wait()
	require.NoError(t, persist.Flush(context.Background()))

	for _, prov := range providers {
		path := filepath.Join(root, prov, "s-"+prov+".jsonl")
		_, err := os.Stat(path)
		require.NoError(t, err, "provider %s ledger", prov)
	}
}

func TestAppendEventsNoLeakedGoroutine(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	key := trajectory.SessionKey{Provider: "cursor", SessionID: "append-sync"}
	ev := []trajectory.Event{{
		Type:      trajectory.TypeHookInvoked,
		Source:    trajectory.SourceHook,
		Provider:  key.Provider,
		SessionID: key.SessionID,
		TS:        time.Now().UTC(),
		Seq:       1,
	}}
	for range 20 {
		require.NoError(t, trajectory.AppendEvents(root, key, ev))
	}
	path := trajectory.SessionFilePath(root, key)
	events, err := trajectory.ReadEvents(path)
	require.NoError(t, err)
	require.Len(t, events, 20)
	assert.Equal(t, trajectory.TypeHookInvoked, events[0].Type)
}

func TestPersisterRequeuesFailedFlush(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	root := filepath.Join(dir, "root")
	require.NoError(t, os.WriteFile(root, []byte("x"), 0o600))

	p := trajectory.NewPersister(root, nil)
	key := trajectory.SessionKey{Provider: "cursor", SessionID: "s1"}
	p.Schedule(key, []trajectory.Event{{
		Type:   trajectory.TypeHookInvoked,
		Source: trajectory.SourceHook,
		TS:     time.Now().UTC(),
	}})

	require.Error(t, p.Flush(context.Background()))

	require.NoError(t, os.Remove(root))
	require.NoError(t, os.MkdirAll(root, 0o700))

	require.NoError(t, p.Flush(context.Background()))
	path := trajectory.SessionFilePath(root, key)
	_, err := os.Stat(path)
	require.NoError(t, err)
	events, err := trajectory.ReadEvents(path)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, trajectory.TypeHookInvoked, events[0].Type)
}
