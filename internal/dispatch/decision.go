package dispatch

import (
	"strings"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

// DecisionToProto maps an agenthooks Decision to the IPC Decision message.
func DecisionToProto(d agenthooks.Decision) *agentdv1.Decision {
	if d == nil {
		return &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION}
	}
	out := &agentdv1.Decision{
		Reason:        d.Reason(),
		SystemMessage: d.SystemMessage(),
	}
	if ctx := d.Context(); len(ctx) > 0 {
		out.Context = strings.Join(ctx, "\n")
	}
	switch d.Kind() {
	case agenthooks.DecisionDeny:
		out.Kind = agentdv1.DecisionKind_DECISION_KIND_DENY
	case agenthooks.DecisionAsk:
		out.Kind = agentdv1.DecisionKind_DECISION_KIND_ASK
	case agenthooks.DecisionAllow:
		out.Kind = agentdv1.DecisionKind_DECISION_KIND_ALLOW
	case agenthooks.DecisionBlockPrompt:
		out.Kind = agentdv1.DecisionKind_DECISION_KIND_BLOCK_PROMPT
	default:
		out.Kind = agentdv1.DecisionKind_DECISION_KIND_NO_DECISION
	}
	return out
}

// NeutralDecision is a NO_DECISION proto value.
func NeutralDecision() *agentdv1.Decision {
	return &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION}
}
