package targets

import (
	"context"
	"io"
	"log/slog"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/guard"
)

// Builtin runs in-process agenthooks guards / observers.
type Builtin struct {
	Guards config.Guards
	Policy config.Policy
	Log    *slog.Logger
}

// Decide runs sync guards via Runner.Decide.
func (b *Builtin) Decide(ctx context.Context, typed any, guardNames []string) (agenthooks.Decision, error) {
	r := b.newRunner(guardNames, false)
	return r.Decide(ctx, typed)
}

// Observe runs async OnAny observers; errors are logged only.
func (b *Builtin) Observe(ctx context.Context, typed any) {
	r := b.newRunner(nil, true)
	// Decide still runs OnAny observers; we ignore the decision.
	_, err := r.Decide(ctx, typed)
	if err != nil && b.logger() != nil {
		b.logger().Warn("builtin observe failed", "error", err)
	}
}

func (b *Builtin) newRunner(guardNames []string, observeOnly bool) *agenthooks.Runner {
	opts := []agenthooks.Option{
		agenthooks.WithLogger(b.logger()),
		agenthooks.WithPolicy(toHooksPolicy(b.Policy)),
	}
	r := agenthooks.New(opts...)
	if observeOnly {
		r.OnAny(func(context.Context, *agenthooks.Event) error { return nil })
		return r
	}
	for _, name := range guardNames {
		switch name {
		case "secrets":
			guard.AttachSecrets(r, b.Guards.Secrets)
		}
	}
	return r
}

func (b *Builtin) logger() *slog.Logger {
	if b.Log != nil {
		return b.Log
	}
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func toHooksPolicy(p config.Policy) agenthooks.Policy {
	out := agenthooks.Policy{}
	switch p.Fail {
	case config.FailClosed:
		out.Fail = agenthooks.FailClosed
	default:
		out.Fail = agenthooks.FailOpen
	}
	switch p.AskFallback {
	case config.AskFallbackDeny:
		out.AskFallback = agenthooks.FallbackDeny
	default:
		out.AskFallback = agenthooks.FallbackNoDecision
	}
	switch p.Unsupported {
	case config.UnsupportedStrict:
		out.Unsupported = agenthooks.Strict
	default:
		out.Unsupported = agenthooks.Degrade
	}
	return out
}
