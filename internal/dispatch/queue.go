package dispatch

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/macrox-pro/agentd/internal/config"
)

// Job is one async dispatch unit.
type Job struct {
	Run func(ctx context.Context)
}

// Queue is a bounded async worker pool.
type Queue struct {
	cfg     config.AsyncConfig
	log     *slog.Logger
	ch      chan Job
	dropped atomic.Uint64
	closed  atomic.Bool
	wg      sync.WaitGroup
}

// NewQueue starts workers for the given async config.
func NewQueue(cfg config.AsyncConfig, log *slog.Logger) *Queue {
	if cfg.QueueCapacity < 1 {
		cfg.QueueCapacity = 1024
	}
	if cfg.WorkerLimit < 1 {
		cfg.WorkerLimit = 8
	}
	if cfg.TargetTimeout <= 0 {
		cfg.TargetTimeout = 30 * time.Second
	}
	q := &Queue{
		cfg: cfg,
		log: log,
		ch:  make(chan Job, cfg.QueueCapacity),
	}
	for i := 0; i < cfg.WorkerLimit; i++ {
		q.wg.Add(1)
		go q.worker()
	}
	return q
}

// Enqueue submits a job. Never blocks the caller; on overflow drops (or logs).
// Returns true if accepted.
func (q *Queue) Enqueue(job Job) bool {
	if q == nil || q.closed.Load() {
		return false
	}
	select {
	case q.ch <- job:
		return true
	default:
		q.dropped.Add(1)
		if q.cfg.OnOverflow == config.OverflowLog && q.log != nil {
			q.log.Warn("async queue overflow; dropping job")
		}
		return false
	}
}

// Depth returns queued jobs waiting for a worker.
func (q *Queue) Depth() int {
	if q == nil {
		return 0
	}
	return len(q.ch)
}

// Dropped returns the overflow drop count.
func (q *Queue) Dropped() uint64 {
	if q == nil {
		return 0
	}
	return q.dropped.Load()
}

// Close stops accepting jobs, closes the channel, and waits for workers up to timeout.
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
		return
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
	case <-timer.C:
		if q.log != nil {
			q.log.Warn("async queue drain timed out")
		}
	}
}

func (q *Queue) worker() {
	defer q.wg.Done()
	for job := range q.ch {
		q.run(job)
	}
}

func (q *Queue) run(job Job) {
	if job.Run == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(context.Background()), q.cfg.TargetTimeout)
	defer cancel()
	job.Run(ctx)
}
