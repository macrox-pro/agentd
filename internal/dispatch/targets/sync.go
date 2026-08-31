package targets

import (
	"context"
	"errors"
	"log/slog"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
)

// SyncRequest is one sync target invocation.
type SyncRequest struct {
	Typed    any
	Raw      []byte
	Provider agentdv1.Provider
	Target   config.CompiledTarget
}

// SyncInvoker runs one sync target and returns a decision that may affect the wire response.
type SyncInvoker interface {
	InvokeSync(ctx context.Context, req SyncRequest) (agenthooks.Decision, error)
}

// builtinSync adapts Builtin.Decide as a SyncInvoker.
type builtinSync struct {
	inner *Builtin
}

// InvokeSync runs in-process guards for the target's guard list.
func (b *builtinSync) InvokeSync(ctx context.Context, req SyncRequest) (agenthooks.Decision, error) {
	if b == nil || b.inner == nil {
		return agenthooks.NoDecision(), nil
	}
	return b.inner.Decide(ctx, req.Typed, req.Target.Guards)
}

// GRPCSync wraps GRPC.InvokeSync with OnError fail_open / fail_closed for the sync factory path.
// GRPC.InvokeSync itself still returns raw errors (tests and callers that need them).
type GRPCSync struct {
	Inner *GRPC
	Log   *slog.Logger
}

// InvokeSync forwards to Inner and maps errors per Target.OnError.
func (w *GRPCSync) InvokeSync(ctx context.Context, req SyncRequest) (agenthooks.Decision, error) {
	g := w.Inner
	if g == nil {
		g = &GRPC{Logger: w.Log}
	}
	d, err := g.InvokeSync(ctx, req)
	if err == nil {
		return d, nil
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.Canceled) {
		return nil, err
	}
	if req.Target.OnError == config.FailOpen {
		if w.Log != nil {
			w.Log.Warn("grpc sync target failed (fail_open)", "error", err)
		}
		return agenthooks.NoDecision(), nil
	}
	return agenthooks.Deny("grpc forward failed"), nil
}
