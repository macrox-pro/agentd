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
	queue    *Queue
	log      *slog.Logger
	sessions *Sessions
}

// NewEngine returns a dispatch engine backed by queue.
func NewEngine(queue *Queue, log *slog.Logger) *Engine {
	return &Engine{queue: queue, log: log, sessions: &Sessions{}}
}

// Queue returns the async queue (may be nil in tests).
func (e *Engine) Queue() *Queue {
	if e == nil {
		return nil
	}
	return e.queue
}

// Sessions returns the per-session lock registry.
func (e *Engine) Sessions() *Sessions {
	if e == nil {
		return nil
	}
	return e.sessions
}

// InvokeInput is one HookService.Invoke.
type InvokeInput struct {
	Provider       agentdv1.Provider
	RawPayload     []byte
	Deadline       time.Time
	Snap           *config.Snapshot
	InvocationMode agentdv1.InvocationMode
	CWD            string
	ProjectRoot    string
}

// InvokeResult is the sync decision plus async enqueue count.
type InvokeResult struct {
	Decision             *agentdv1.Decision
	AsyncDispatchedCount uint32
	Meta                 InvokeMeta
}

// Invoke decodes, matches a route, runs sync/async pipelines per mode.
func (e *Engine) Invoke(ctx context.Context, in InvokeInput) (InvokeResult, error) {
	if in.Snap == nil {
		return InvokeResult{}, fmt.Errorf("dispatch: nil snapshot")
	}
	typed, err := DecodeTyped(ctx, in.Provider, in.InvocationMode, in.RawPayload)
	if err != nil {
		return InvokeResult{}, err
	}
	providerName, _ := providerName(in.Provider)
	route := MatchRoute(in.Snap.Routes, typed)
	if route == nil {
		return InvokeResult{
			Decision: NeutralDecision(),
			Meta:     MetaFromTyped(providerName, typed, false),
		}, nil
	}

	mode := config.NormalizeMode(route.Mode)
	builtin := &targets.Builtin{
		Guards:          in.Snap.Guards,
		Policy:          in.Snap.Policy,
		Approvals:       in.Snap.Approvals,
		TemporaryBlocks: in.Snap.TemporaryBlocks,
		ProjectRoot:     targets.ProjectRootOf(in.Snap),
		Log:             e.log,
	}
	eventKind := targets.EventKindOf(typed)

	unlock := e.sessions.Lock(SessionIDOf(typed))
	defer unlock()

	budget := SyncBudget(time.Now(), in.Deadline, eventKind, route.SyncTimeout)
	if budget > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, budget)
		defer cancel()
	}

	meta := MetaFromTyped(providerName, typed, true)

	switch mode {
	case config.ModeAsyncOnly:
		n := e.enqueueAsync(builtin, route.Async, typed, in.RawPayload, providerName, eventKind, nil)
		return InvokeResult{Decision: NeutralDecision(), AsyncDispatchedCount: n, Meta: meta}, nil

	case config.ModeParallel:
		n := e.enqueueAsync(builtin, route.Async, typed, in.RawPayload, providerName, eventKind, nil)
		d, err := e.runSync(ctx, builtin, route.Sync, typed, in.Provider, in.RawPayload)
		if err != nil {
			return InvokeResult{}, err
		}
		return InvokeResult{Decision: DecisionToProto(d), AsyncDispatchedCount: n, Meta: meta}, nil

	case config.ModeAfterSync:
		d, err := e.runSync(ctx, builtin, route.Sync, typed, in.Provider, in.RawPayload)
		if err != nil {
			return InvokeResult{}, err
		}
		proto := DecisionToProto(d)
		outcome := &targets.SyncOutcome{Kind: proto.GetKind(), Reason: proto.GetReason()}
		n := e.enqueueAsync(builtin, route.Async, typed, in.RawPayload, providerName, eventKind, outcome)
		return InvokeResult{Decision: proto, AsyncDispatchedCount: n, Meta: meta}, nil

	default: // sync_only
		d, err := e.runSync(ctx, builtin, route.Sync, typed, in.Provider, in.RawPayload)
		if err != nil {
			return InvokeResult{}, err
		}
		return InvokeResult{Decision: DecisionToProto(d), AsyncDispatchedCount: 0, Meta: meta}, nil
	}
}

func (e *Engine) runSync(ctx context.Context, b *targets.Builtin, syncTargets []config.CompiledTarget, typed any, provider agentdv1.Provider, raw []byte) (agenthooks.Decision, error) {
	var last agenthooks.Decision = agenthooks.NoDecision()
	for _, t := range syncTargets {
		var (
			d   agenthooks.Decision
			err error
		)
		switch t.Kind {
		case config.TargetBuiltin:
			d, err = b.Decide(ctx, typed, t.Guards)
		case config.TargetGRPC:
			g := &targets.GRPC{Logger: e.log}
			d, err = g.InvokeSync(ctx, targets.SyncRequest{
				Provider: provider,
				Raw:      raw,
				Target:   t,
			})
			if err != nil {
				if t.OnError == config.FailOpen {
					if e.log != nil {
						e.log.Warn("grpc sync target failed (fail_open)", "error", err)
					}
					d, err = agenthooks.NoDecision(), nil
				} else {
					d, err = agenthooks.Deny("grpc forward failed"), nil
				}
			}
		default:
			continue
		}
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

func (e *Engine) enqueueAsync(
	b *targets.Builtin,
	asyncTargets []config.CompiledTarget,
	typed any,
	raw []byte,
	provider, eventKind string,
	outcome *targets.SyncOutcome,
) uint32 {
	if e == nil || e.queue == nil {
		return 0
	}
	var n uint32
	for _, t := range asyncTargets {
		inv, err := targets.NewAsyncInvoker(t, b, e.log)
		if err != nil {
			if e.log != nil {
				e.log.Warn("skip async target", "error", err)
			}
			continue
		}
		req := targets.AsyncRequest{
			Typed:       typed,
			Raw:         raw,
			Provider:    provider,
			EventKind:   eventKind,
			Target:      t,
			SyncOutcome: outcome,
		}
		ok := e.queue.Enqueue(Job{Run: func(ctx context.Context) {
			if err := inv.InvokeAsync(ctx, req); err != nil && e.log != nil {
				e.log.Warn("async target failed", "target", string(t.Kind), "error", err)
			}
		}})
		if ok {
			n++
		}
	}
	return n
}
