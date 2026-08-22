package trajectory_test

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestStoreContiguousSeq(t *testing.T) {
	t.Parallel()
	store := trajectory.NewStore()
	key := trajectory.SessionKey{Provider: "claude-code", SessionID: "s1"}
	ev := []trajectory.Event{{Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, TS: time.Now().UTC()}}
	for i := range 3 {
		got := store.Append(key, ev)
		require.Len(t, got, 1, "append %d", i)
		assert.Equal(t, uint64(i+1), got[0].Seq, "append %d", i)
	}
	all := store.Events(key)
	require.Len(t, all, 3)
	for i, e := range all {
		assert.Equal(t, uint64(i+1), e.Seq, "stored seq")
	}
}

func TestStoreAppendImmutable(t *testing.T) {
	t.Parallel()
	store := trajectory.NewStore()
	key := trajectory.SessionKey{Provider: "cursor", SessionID: "s2"}
	appended := store.Append(key, []trajectory.Event{{
		Type:   trajectory.TypeHookInvoked,
		Source: trajectory.SourceHook,
		TS:     time.Now().UTC(),
		Data:   json.RawMessage(`{"kind":"tool.pre"}`),
	}})
	require.Len(t, appended, 1)
	first := store.Events(key)[0]
	copied := first
	copied.Data = json.RawMessage(`{"kind":"mutated"}`)
	_ = copied
	second := store.Events(key)[0]
	assert.JSONEq(t, `{"kind":"tool.pre"}`, string(second.Data))
}

func TestTruncateMaxEventBytes(t *testing.T) {
	t.Parallel()
	raw := make([]byte, 300)
	for i := range raw {
		raw[i] = 'x'
	}
	cfg := config.TrajectoryConfig{Enabled: true, IncludeRaw: true, MaxEventBytes: 128}
	got := trajectory.PrepareRaw(raw, cfg)
	require.Len(t, got, 128)
}

func TestSessionOpenOnce(t *testing.T) {
	t.Parallel()
	store := trajectory.NewStore()
	key := trajectory.SessionKey{Provider: "gemini", SessionID: "s3"}
	open := trajectory.Event{Type: trajectory.TypeSessionOpen, Source: trajectory.SourceSystem, TS: time.Now().UTC()}
	store.Append(key, []trajectory.Event{open, {Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, TS: time.Now().UTC()}})
	store.Append(key, []trajectory.Event{open, {Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, TS: time.Now().UTC()}})
	events := store.Events(key)
	var opens int
	for _, e := range events {
		if e.Type == trajectory.TypeSessionOpen {
			opens++
		}
	}
	assert.Equal(t, 1, opens, "session/open once")
}

func TestQueueOverflowDrop(t *testing.T) {
	t.Parallel()
	store := trajectory.NewStore()
	q := trajectory.NewQueue(1, store, nil, nil, nil)
	defer q.Close(0)
	key := trajectory.SessionKey{Provider: "codex", SessionID: "s4"}
	ev := []trajectory.Event{{Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, TS: time.Now().UTC()}}
	var dropped atomic.Bool
	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 64 {
				if !q.Enqueue(key, ev) {
					dropped.Store(true)
				}
			}
		}()
	}
	wg.Wait()
	require.True(t, dropped.Load(), "expected overflow drop")
	assert.GreaterOrEqual(t, q.Dropped(), uint64(1))
}

func TestResolveSessionKeyWeakID(t *testing.T) {
	t.Parallel()
	k1 := trajectory.ResolveSessionKey("kimi-code", "", "/proj", "/cwd")
	k2 := trajectory.ResolveSessionKey("kimi-code", "", "/proj", "/cwd")
	assert.Equal(t, k1.SessionID, k2.SessionID)
	assert.True(t, len(k1.SessionID) > 0)
	assert.Equal(t, "kimi-code", k1.Provider)
}
