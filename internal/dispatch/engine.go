package dispatch

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/decision"
	"github.com/macrox-pro/agentd/internal/dispatch/targets"
	"github.com/macrox-pro/agentd/internal/provider"
)

const syncPipelineFailReason = "sync pipeline failed (fail-closed policy)"

// Observer records invoke_sync and async_side histogram samples.
type Observer interface {
	ObserveInvoke(provider, eventKind, decision, outcome string, seconds float64)
	ObserveAsync(targetKind, result string, seconds float64)
}

// Engine routes hook invocations through sync and async pipelines.
type Engine struct {
	queue    *Queue
	log      *slog.Logger
	sessions *Sessions
	observer Observer
}

// NewEngine returns a dispatch engine backed by queue.
func NewEngine(queue *Queue, log *slog.Logger, observer Observer) *Engine {
	return &Engine{queue: queue, log: log, sessions: &Sessions{}, observer: observer}
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

// Invoker runs the sync/async dispatch pipeline for one hook invocation.
type Invoker interface {
	Invoke(ctx context.Context, in InvokeInput) (InvokeResult, error)
}

var _ Invoker = (*Engine)(nil)

// InvokeResult is the sync decision plus async enqueue count.
type InvokeResult struct {
	Decision             *agentdv1.Decision
	AsyncDispatchedCount uint32
	Meta                 InvokeMeta
}

func decisionLabel(d *agentdv1.Decision) string {
	if d == nil {
		return agentdv1.DecisionKind_DECISION_KIND_UNSPECIFIED.String()
	}
	return d.GetKind().String()
}

func invokeOutcome(ctx context.Context, err error) string {
	if err == nil {
		return "ok"
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "timeout"
	}
	return "error"
}

// Invoke decodes, matches a route, runs sync/async pipelines per mode.
func (e *Engine) Invoke(ctx context.Context, in InvokeInput) (InvokeResult, error) {
	start := time.Now()
	id, _ := provider.FromProto(in.Provider)
	providerName := string(id)
	eventKind := ""
	decisionStr := agentdv1.DecisionKind_DECISION_KIND_UNSPECIFIED.String()
	outcome := "ok"
	defer func() {
		if e != nil && e.observer != nil {
			e.observer.ObserveInvoke(providerName, eventKind, decisionStr, outcome, time.Since(start).Seconds())
		}
	}()

	if in.Snap == nil {
		outcome = "error"
		return InvokeResult{}, fmt.Errorf("dispatch: nil snapshot")
	}
	typed, err := DecodeTyped(ctx, in.Provider, in.InvocationMode, in.RawPayload)
	if err != nil {
		outcome = invokeOutcome(ctx, err)
		return InvokeResult{}, err
	}
	eventKind = targets.EventKindOf(typed)
	route := MatchRoute(in.Snap.Routes, typed)
	if route == nil {
		meta := MetaFromTyped(providerName, typed, false)
		res := InvokeResult{
			Decision: decision.Neutral(),
			Meta:     meta,
		}
		decisionStr = decisionLabel(res.Decision)
		e.logInvokeDebug(providerName, meta.EventKind, "", res.Decision)
		return res, nil
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
		res := InvokeResult{Decision: decision.Neutral(), AsyncDispatchedCount: n, Meta: meta}
		decisionStr = decisionLabel(res.Decision)
		e.logInvokeDebug(providerName, eventKind, route.Name, res.Decision)
		return res, nil

	default:
		var asyncN uint32
		if mode == config.ModeParallel {
			asyncN = e.enqueueAsync(builtin, route.Async, typed, in.RawPayload, providerName, eventKind, nil)
		}

		d, syncErr := e.runSync(ctx, builtin, route.Sync, typed, in.Provider, in.RawPayload)
		if syncErr != nil {
			outcome = invokeOutcome(ctx, syncErr)
			if e.log != nil {
				e.log.Warn("sync pipeline failed",
					"provider", providerName,
					"event_kind", eventKind,
					"route", route.Name,
					"policy_fail", in.Snap.Policy.Fail,
					"error", syncErr,
				)
			}
			d = syncFailureDecision(in.Snap.Policy.Fail, typed)
		}

		proto := decision.ToProto(d)
		if mode == config.ModeAfterSync {
			syncOutcome := &targets.SyncOutcome{Kind: proto.GetKind(), Reason: proto.GetReason()}
			asyncN = e.enqueueAsync(builtin, route.Async, typed, in.RawPayload, providerName, eventKind, syncOutcome)
		}

		res := InvokeResult{Decision: proto, AsyncDispatchedCount: asyncN, Meta: meta}
		decisionStr = decisionLabel(res.Decision)
		e.logInvokeDebug(providerName, eventKind, route.Name, res.Decision)
		return res, nil
	}
}

func syncFailureDecision(mode config.FailMode, typed any) agenthooks.Decision {
	if mode != config.FailClosed {
		return agenthooks.NoDecision()
	}
	base := agenthooks.EventOf(typed)
	if base == nil || !base.Can(agenthooks.CapDeny) {
		return agenthooks.NoDecision()
	}
	switch base.Kind {
	case agenthooks.KindToolPre, agenthooks.KindPermission:
		return agenthooks.Deny(syncPipelineFailReason)
	case agenthooks.KindPromptSubmitted:
		return agenthooks.BlockPrompt(syncPipelineFailReason)
	default:
		return agenthooks.NoDecision()
	}
}

func (e *Engine) logInvokeDebug(provider, eventKind, routeName string, decision *agentdv1.Decision) {
	if e == nil || e.log == nil {
		return
	}
	kind := agentdv1.DecisionKind_DECISION_KIND_NO_DECISION
	if decision != nil {
		kind = decision.GetKind()
	}
	e.log.Debug("invoke",
		"provider", provider,
		"event_kind", eventKind,
		"route", routeName,
		"decision_kind", kind.String(),
	)
}

func (e *Engine) runSync(ctx context.Context, b *targets.Builtin, syncTargets []config.CompiledTarget, typed any, provider agentdv1.Provider, raw []byte) (agenthooks.Decision, error) {
	var last agenthooks.Decision = agenthooks.NoDecision()
	for _, t := range syncTargets {
		inv, err := targets.NewSyncInvoker(t, b, e.log)
		if err != nil {
			if e.log != nil {
				e.log.Warn("skip sync target", "target", string(t.Kind), "error", err)
			}
			continue
		}
		d, err := inv.InvokeSync(ctx, targets.SyncRequest{
			Typed:    typed,
			Raw:      raw,
			Provider: provider,
			Target:   t,
		})
		if err != nil {
			return nil, err
		}
		last = d
		if FirstConclusive(d) {
			return d, nil
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
		targetKind := string(t.Kind)
		ok := e.queue.Enqueue(Job{Run: func(ctx context.Context) {
			start := time.Now()
			runErr := inv.InvokeAsync(ctx, req)
			result := "ok"
			if runErr != nil {
				result = "error"
				if e.log != nil {
					e.log.Warn("async target failed", "target", targetKind, "error", runErr)
				}
			}
			if e.observer != nil {
				e.observer.ObserveAsync(targetKind, result, time.Since(start).Seconds())
			}
		}})
		if ok {
			n++
		}
	}
	return n
}
