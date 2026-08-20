package dispatch_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestQueueEnqueueNonBlocking(t *testing.T) {
	t.Parallel()
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 2,
		WorkerLimit:   1,
		TargetTimeout: time.Second,
		OnOverflow:    config.OverflowDrop,
	}, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })

	started := make(chan struct{})
	release := make(chan struct{})
	require.True(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {
		close(started)
		<-release
	}}), "first job")
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not start")
	}
	require.True(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {}}), "queue slot 1")
	require.True(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {}}), "queue slot 2")
	assert.False(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {}}), "overflow")
	assert.Equal(t, uint64(1), q.Dropped(), "Dropped()")
	close(release)
}

func TestQueueDrain(t *testing.T) {
	t.Parallel()
	var ran atomic.Int32
	var mu sync.Mutex
	q := dispatch.NewQueue(config.AsyncConfig{
		QueueCapacity: 8,
		WorkerLimit:   2,
		TargetTimeout: time.Second,
	}, nil)
	for i := 0; i < 5; i++ {
		require.True(t, q.Enqueue(dispatch.Job{Run: func(context.Context) {
			mu.Lock()
			ran.Add(1)
			mu.Unlock()
		}}))
	}
	q.Close(2 * time.Second)
	assert.Equal(t, int32(5), ran.Load(), "all jobs drained")
}
