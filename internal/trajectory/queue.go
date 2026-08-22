package trajectory

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

const persistDebounce = 50 * time.Millisecond

// appendJob is one async store write unit.
type appendJob struct {
	key    SessionKey
	events []Event
}

// Queue is a bounded worker pool for trajectory append + persist.
type Queue struct {
	capacity int
	ch       chan appendJob
	dropped  atomic.Uint64
	closed   atomic.Bool
	wg       sync.WaitGroup
	store    *Store
	persist  *Persister
	log      *slog.Logger
}

// NewQueue starts workers for trajectory side effects.
func NewQueue(capacity int, store *Store, persist *Persister, log *slog.Logger) *Queue {
	if capacity < 1 {
		capacity = 1024
	}
	q := &Queue{
		capacity: capacity,
		ch:       make(chan appendJob, capacity),
		store:    store,
		persist:  persist,
		log:      log,
	}
	workers := 2
	if workers > capacity {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go q.worker()
	}
	return q
}

// Enqueue submits events for async append. Never blocks the caller.
func (q *Queue) Enqueue(key SessionKey, events []Event) bool {
	if q == nil || q.closed.Load() || len(events) == 0 {
		return false
	}
	select {
	case q.ch <- appendJob{key: key, events: events}:
		return true
	default:
		q.dropped.Add(1)
		if q.log != nil {
			q.log.Warn("trajectory queue overflow; dropping events")
		}
		return false
	}
}

// Dropped returns overflow drop count.
func (q *Queue) Dropped() uint64 {
	if q == nil {
		return 0
	}
	return q.dropped.Load()
}

// Close drains pending jobs up to timeout then stops workers.
func (q *Queue) Close(timeout time.Duration) {
	if q == nil || !q.closed.CompareAndSwap(false, true) {
		return
	}
	close(q.ch)
	done := make(chan struct{})
	go func() {
		q.wg.Wait()
		close(done)
	}()
	if timeout <= 0 {
		<-done
		q.flushPersist()
		return
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
	case <-timer.C:
		if q.log != nil {
			q.log.Warn("trajectory queue drain timed out")
		}
	}
	q.flushPersist()
}

func (q *Queue) flushPersist() {
	if q.persist == nil {
		return
	}
	if err := q.persist.Flush(context.Background()); err != nil && q.log != nil {
		q.log.Warn("trajectory persist flush failed", "error", err)
	}
}

func (q *Queue) worker() {
	defer q.wg.Done()
	for job := range q.ch {
		q.run(job)
	}
}

func (q *Queue) run(job appendJob) {
	if q.store == nil {
		return
	}
	appended := q.store.Append(job.key, job.events)
	if q.persist != nil && len(appended) > 0 {
		q.persist.Schedule(job.key, appended)
	}
}
