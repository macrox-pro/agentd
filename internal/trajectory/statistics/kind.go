package statistics

import (
	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/speakeasy-api/agenthooks"
)

// HookKind maps agenthooks ledger kind strings to proto EventKind.
func HookKind(kind string) agentdv1.EventKind {
	switch agenthooks.EventKind(kind) {
	case agenthooks.KindSessionStart:
		return agentdv1.EventKind_EVENT_KIND_SESSION_START
	case agenthooks.KindSessionEnd:
		return agentdv1.EventKind_EVENT_KIND_SESSION_END
	case agenthooks.KindPromptSubmitted:
		return agentdv1.EventKind_EVENT_KIND_PROMPT_SUBMITTED
	case agenthooks.KindToolPre:
		return agentdv1.EventKind_EVENT_KIND_TOOL_PRE
	case agenthooks.KindToolPost:
		return agentdv1.EventKind_EVENT_KIND_TOOL_POST
	case agenthooks.KindToolError:
		return agentdv1.EventKind_EVENT_KIND_TOOL_ERROR
	case agenthooks.KindPermission:
		return agentdv1.EventKind_EVENT_KIND_PERMISSION
	case agenthooks.KindStop:
		return agentdv1.EventKind_EVENT_KIND_AGENT_STOP
	case agenthooks.KindNotification:
		return agentdv1.EventKind_EVENT_KIND_NOTIFICATION
	default:
		if kind == "permission.request" {
			return agentdv1.EventKind_EVENT_KIND_PERMISSION
		}
		return agentdv1.EventKind_EVENT_KIND_OTHER
	}
}
