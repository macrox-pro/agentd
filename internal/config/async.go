package config

import (
	"fmt"
	"time"
)

type OverflowMode string

const (
	OverflowDrop OverflowMode = "drop"
	OverflowLog  OverflowMode = "log"
)

// AsyncConfig is the compiled async queue settings.
type AsyncConfig struct {
	QueueCapacity int
	WorkerLimit   int
	TargetTimeout time.Duration
	OnOverflow    OverflowMode
}

func parseAsync(fa *fileAsync, def AsyncConfig) (AsyncConfig, error) {
	out := def
	if fa == nil {
		return out, nil
	}
	if fa.QueueCapacity != 0 {
		if fa.QueueCapacity < 1 {
			return AsyncConfig{}, fmt.Errorf("async.queue_capacity must be >= 1")
		}
		out.QueueCapacity = fa.QueueCapacity
	}
	if fa.WorkerLimit != 0 {
		if fa.WorkerLimit < 1 {
			return AsyncConfig{}, fmt.Errorf("async.worker_limit must be >= 1")
		}
		out.WorkerLimit = fa.WorkerLimit
	}
	if fa.TargetTimeout != "" {
		d, err := time.ParseDuration(fa.TargetTimeout)
		if err != nil {
			return AsyncConfig{}, fmt.Errorf("async.target_timeout: %w", err)
		}
		if d <= 0 {
			return AsyncConfig{}, fmt.Errorf("async.target_timeout must be > 0")
		}
		out.TargetTimeout = d
	}
	if fa.OnOverflow != "" {
		switch OverflowMode(fa.OnOverflow) {
		case OverflowDrop, OverflowLog:
			out.OnOverflow = OverflowMode(fa.OnOverflow)
		default:
			return AsyncConfig{}, fmt.Errorf("async.on_overflow: unknown %q", fa.OnOverflow)
		}
	}
	return out, nil
}
