package targets

import (
	"context"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
)

// SyncOutcome is the sync pipeline result passed to after_sync async jobs.
type SyncOutcome struct {
	Kind   agentdv1.DecisionKind
	Reason string
}

// AsyncRequest is one async target invocation.
type AsyncRequest struct {
	Typed       any
	Raw         []byte
	Provider    string
	EventKind   string
	Target      config.CompiledTarget
	SyncOutcome *SyncOutcome // nil for parallel/async_only
}

// AsyncInvoker runs one async target without affecting the sync decision.
type AsyncInvoker interface {
	InvokeAsync(ctx context.Context, req AsyncRequest) error
}

// EventKindOf extracts unified kind from a typed agenthooks event.
func EventKindOf(typed any) string {
	if base := agenthooks.EventOf(typed); base != nil && base.Kind != "" {
		return string(base.Kind)
	}
	return string(agenthooks.KindOther)
}
