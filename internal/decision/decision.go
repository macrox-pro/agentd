// Package decision maps between agenthooks Decision and the agentd.v1 Decision proto.
//
// Owns: proto↔agenthooks Decision Kind mapping (ToProto, FromProto, Neutral).
// Must not: routing, wire I/O, guards, config compile.
//
// Invariants:
//   - Deny/Ask/Allow copy Reason, SystemMessage, and Context.
//   - BLOCK_PROMPT maps to BlockPrompt(reason) without tool extras.
//   - nil / unknown Kind → NoDecision / NO_DECISION.
//
// Entry: ToProto, FromProto, Neutral.
// See DESIGN.md §1.5 (invoke_sync).
package decision

import (
	"strings"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

// ToProto maps an agenthooks Decision to the IPC Decision message.
func ToProto(d agenthooks.Decision) *agentdv1.Decision {
	if d == nil {
		return Neutral()
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

// FromProto maps a daemon Decision to an agenthooks Decision.
func FromProto(d *agentdv1.Decision) agenthooks.Decision {
	if d == nil {
		return agenthooks.NoDecision()
	}
	switch d.GetKind() {
	case agentdv1.DecisionKind_DECISION_KIND_DENY:
		return withToolExtras(agenthooks.Deny(d.GetReason()), d)
	case agentdv1.DecisionKind_DECISION_KIND_ASK:
		return withToolExtras(agenthooks.AskUser(d.GetReason()), d)
	case agentdv1.DecisionKind_DECISION_KIND_ALLOW:
		return withToolExtras(agenthooks.Allow(), d)
	case agentdv1.DecisionKind_DECISION_KIND_BLOCK_PROMPT:
		return agenthooks.BlockPrompt(d.GetReason())
	default:
		return agenthooks.NoDecision()
	}
}

// Neutral is a NO_DECISION proto value.
func Neutral() *agentdv1.Decision {
	return &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION}
}

func withToolExtras(out agenthooks.ToolPreDecision, d *agentdv1.Decision) agenthooks.ToolPreDecision {
	if msg := d.GetSystemMessage(); msg != "" {
		out = out.WithSystemMessage(msg)
	}
	if c := d.GetContext(); c != "" {
		out = out.WithContext(c)
	}
	return out
}
