package dispatch

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch/targets"
)

// Engine routes hook invocations through sync and async pipelines.
type Engine struct {
	queue *Queue
	log   *slog.Logger
}

// NewEngine returns a dispatch engine backed by queue.
func NewEngine(queue *Queue, log *slog.Logger) *Engine {
	return &Engine{queue: queue, log: log}
}

// Queue returns the async queue (may be nil in tests).
func (e *Engine) Queue() *Queue {
	if e == nil {
		return nil
	}
	return e.queue
}

// InvokeInput is one HookService.Invoke.
type InvokeInput struct {
	Provider   agentdv1.Provider
	RawPayload []byte
	Deadline   time.Time
	Snap       *config.Snapshot
}

// InvokeResult is the sync decision plus async enqueue count.
type InvokeResult struct {
	Decision             *agentdv1.Decision
	AsyncDispatchedCount uint32
}

// Invoke decodes, matches a route, runs sync/async pipelines per mode.
func (e *Engine) Invoke(ctx context.Context, in InvokeInput) (InvokeResult, error) {
	if in.Snap == nil {
		return InvokeResult{}, fmt.Errorf("dispatch: nil snapshot")
	}
	typed, err := DecodeTyped(ctx, in.Provider, in.RawPayload)
	if err != nil {
		return InvokeResult{}, err
	}
	route := MatchRoute(in.Snap.Routes, typed)
	if route == nil {
		return InvokeResult{Decision: NeutralDecision()}, nil
	}

	mode := config.NormalizeMode(route.Mode)
	builtin := &targets.Builtin{
		Guards: in.Snap.Guards,
		Policy: in.Snap.Policy,
		Log:    e.log,
	}

	switch mode {
	case config.ModeAsyncOnly:
		n := e.enqueueAsync(builtin, route.Async, typed, nil)
		return InvokeResult{Decision: NeutralDecision(), AsyncDispatchedCount: n}, nil

	case config.ModeParallel:
		n := e.enqueueAsync(builtin, route.Async, typed, nil)
		d, err := e.runSync(ctx, builtin, route.Sync, typed)
		if err != nil {
			return InvokeResult{}, err
		}
		return InvokeResult{Decision: DecisionToProto(d), AsyncDispatchedCount: n}, nil

	case config.ModeAfterSync:
		d, err := e.runSync(ctx, builtin, route.Sync, typed)
		if err != nil {
			return InvokeResult{}, err
		}
		n := e.enqueueAsync(builtin, route.Async, typed, d)
		return InvokeResult{Decision: DecisionToProto(d), AsyncDispatchedCount: n}, nil

	default: // sync_only
		d, err := e.runSync(ctx, builtin, route.Sync, typed)
		if err != nil {
			return InvokeResult{}, err
		}
		return InvokeResult{Decision: DecisionToProto(d), AsyncDispatchedCount: 0}, nil
	}
}

func (e *Engine) runSync(ctx context.Context, b *targets.Builtin, syncTargets []config.CompiledTarget, typed any) (agenthooks.Decision, error) {
	var last agenthooks.Decision = agenthooks.NoDecision()
	for _, t := range syncTargets {
		if t.Kind != config.TargetBuiltin {
			continue
		}
		d, err := b.Decide(ctx, typed, t.Guards)
		if err != nil {
			return nil, err
		}
		last = d
		if d != nil && d.Kind() != agenthooks.DecisionNoDecision {
			return d, nil // first_conclusive
		}
	}
	return last, nil
}

func (e *Engine) enqueueAsync(b *targets.Builtin, asyncTargets []config.CompiledTarget, typed any, _ agenthooks.Decision) uint32 {
	if e == nil || e.queue == nil {
		return 0
	}
	var n uint32
	for _, t := range asyncTargets {
		if t.Kind != config.TargetBuiltin || !t.Observe {
			continue
		}
		jobTyped := typed
		ok := e.queue.Enqueue(Job{Run: func(ctx context.Context) {
			b.Observe(ctx, jobTyped)
		}})
		if ok {
			n++
		}
	}
	return n
}
