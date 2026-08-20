package hookedge

import (
	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

// fromProto maps a daemon Decision to an agenthooks Decision for the wire path.
func fromProto(d *agentdv1.Decision) agenthooks.Decision {
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

func withToolExtras(out agenthooks.ToolPreDecision, d *agentdv1.Decision) agenthooks.ToolPreDecision {
	if msg := d.GetSystemMessage(); msg != "" {
		out = out.WithSystemMessage(msg)
	}
	if c := d.GetContext(); c != "" {
		out = out.WithContext(c)
	}
	return out
}

func toolPreFromProto(d *agentdv1.Decision) agenthooks.ToolPreDecision {
	dec := fromProto(d)
	if tp, ok := dec.(agenthooks.ToolPreDecision); ok {
		return tp
	}
	return agenthooks.NoDecision()
}

func promptFromProto(d *agentdv1.Decision) agenthooks.PromptDecision {
	if d == nil {
		return agenthooks.AcceptPrompt()
	}
	switch d.GetKind() {
	case agentdv1.DecisionKind_DECISION_KIND_DENY, agentdv1.DecisionKind_DECISION_KIND_BLOCK_PROMPT:
		return agenthooks.BlockPrompt(d.GetReason())
	default:
		return agenthooks.AcceptPrompt()
	}
}
