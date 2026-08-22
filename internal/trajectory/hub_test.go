package trajectory_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func sampleEvent(typ, source, provider, session string) trajectory.Event {
	return trajectory.Event{
		Type:      typ,
		Source:    source,
		Provider:  provider,
		SessionID: session,
		TS:        time.Now().UTC(),
	}
}

func TestHubDeliverAfterAppend(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	ch, unregister := hub.Register(trajectory.SubscribeFilter{})
	t.Cleanup(unregister)

	ev := sampleEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "claude-code", "s1")
	hub.Publish([]trajectory.Event{ev})

	select {
	case got := <-ch:
		assert.Equal(t, trajectory.TypeHookInvoked, got.Type)
		assert.Equal(t, "claude-code", got.Provider)
	case <-time.After(time.Second):
		t.Fatal("expected delivered event")
	}
}

func TestHubSlowConsumerDrop(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	ch, unregister := hub.Register(trajectory.SubscribeFilter{})
	t.Cleanup(unregister)

	for i := 0; i < 100; i++ {
		hub.Publish([]trajectory.Event{sampleEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "cursor", "s1")})
	}

	done := make(chan struct{})
	go func() {
		<-ch
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected at least one event despite drops")
	}
}

func TestHubUnregister(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	ch, unregister := hub.Register(trajectory.SubscribeFilter{})

	hub.Publish([]trajectory.Event{sampleEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "cursor", "s1")})
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("expected event before unregister")
	}

	unregister()
	hub.Publish([]trajectory.Event{sampleEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "cursor", "s2")})

	select {
	case ev, ok := <-ch:
		if ok {
			assert.Equal(t, "s2", ev.SessionID)
		}
	case <-time.After(100 * time.Millisecond):
	}
}

func TestHubMultipleSubscribers(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	ch1, un1 := hub.Register(trajectory.SubscribeFilter{})
	ch2, un2 := hub.Register(trajectory.SubscribeFilter{})
	t.Cleanup(un1)
	t.Cleanup(un2)

	ev := sampleEvent(trajectory.TypeHookDecided, trajectory.SourceDecision, "codex", "s2")
	hub.Publish([]trajectory.Event{ev})

	for i, ch := range []<-chan trajectory.Event{ch1, ch2} {
		select {
		case got := <-ch:
			assert.Equal(t, trajectory.TypeHookDecided, got.Type, "sub %d", i)
		case <-time.After(time.Second):
			t.Fatalf("sub %d: timeout", i)
		}
	}
}

func TestHubPublishNoSubscribers(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	require.NotPanics(t, func() {
		hub.Publish([]trajectory.Event{sampleEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "gemini", "s1")})
	})
}

func TestHubCloseEndsSubscribers(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	ch, unregister := hub.Register(trajectory.SubscribeFilter{})
	t.Cleanup(unregister)

	hub.Close()

	_, ok := <-ch
	assert.False(t, ok, "channel closed after hub close")
}

func TestSubscribeDoesNotBlockEnqueue(t *testing.T) {
	t.Parallel()
	store := trajectory.NewStore()
	hub := trajectory.NewHub(nil)
	q := trajectory.NewQueue(8, store, nil, hub, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })

	_, unregister := hub.Register(trajectory.SubscribeFilter{})
	t.Cleanup(unregister)

	key := trajectory.SessionKey{Provider: "cursor", SessionID: "block-test"}
	done := make(chan struct{})
	go func() {
		for i := 0; i < 200; i++ {
			q.Enqueue(key, []trajectory.Event{sampleEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "cursor", "block-test")})
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Enqueue blocked")
	}
}

func TestSchemaVersionOnAppend(t *testing.T) {
	t.Parallel()
	store := trajectory.NewStore()
	key := trajectory.SessionKey{Provider: "claude-code", SessionID: "schema"}
	appended := store.Append(key, []trajectory.Event{{
		Type:   trajectory.TypeHookInvoked,
		Source: trajectory.SourceHook,
	}})
	require.Len(t, appended, 1)
	assert.Equal(t, trajectory.SchemaVersion, appended[0].SchemaVersion)

	root := t.TempDir()
	err := trajectory.AppendEvents(root, key, []trajectory.Event{{
		Type:   trajectory.TypeSessionFork,
		Source: trajectory.SourceSystem,
	}})
	require.NoError(t, err)
	events, err := trajectory.ReadEvents(trajectory.SessionFilePath(root, key))
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, trajectory.SchemaVersion, events[0].SchemaVersion)
}

func TestIgnorablePreserved(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	ch, unregister := hub.Register(trajectory.SubscribeFilter{Source: trajectory.SourceTranscript})
	t.Cleanup(unregister)

	ev := sampleEvent(trajectory.TypeTranscriptMessage, trajectory.SourceTranscript, "claude-code", "s1")
	ev.Ignorable = true
	hub.Publish([]trajectory.Event{ev})

	select {
	case got := <-ch:
		assert.True(t, got.Ignorable)
		assert.Equal(t, trajectory.TypeTranscriptMessage, got.Type)
	case <-time.After(time.Second):
		t.Fatal("expected ignorable transcript event")
	}
}

func TestHubConcurrentRegisterPublish(t *testing.T) {
	t.Parallel()
	hub := trajectory.NewHub(nil)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch, unregister := hub.Register(trajectory.SubscribeFilter{})
			defer unregister()
			hub.Publish([]trajectory.Event{sampleEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "opencode", "s1")})
			select {
			case <-ch:
			case <-time.After(time.Second):
			}
		}()
	}
	wg.Wait()
}
