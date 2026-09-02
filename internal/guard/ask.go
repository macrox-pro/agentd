package guard

import (
	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

// askUnsupported returns the decision when a guard would Ask but the event cannot express Ask.
func askUnsupported(fallback config.AskFallback, reason string) agenthooks.ToolPreDecision {
	if fallback == config.AskFallbackNoDecision {
		return agenthooks.NoDecision()
	}
	return agenthooks.Deny(reason)
}
