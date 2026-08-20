package dispatch

import (
	"time"

	"github.com/speakeasy-api/agenthooks"
)

// Timeout margin keeps 90% of the provider budget (same ratio as agenthooks).
const timeoutMarginDenom = 10

// Kind timeouts aligned with internal/install HookSpec defaults.
const (
	kindTimeoutLong  = 30 * time.Second // tool.pre, prompt.submitted
	kindTimeoutShort = 5 * time.Second  // stop, session, notification, …
)

// SyncBudget returns the sync pipeline duration:
// min(provider_timeout - margin, route.sync_timeout).
// invokeDeadline zero means derive provider_timeout from eventKind defaults.
// routeSyncTimeout zero means no route cap.
func SyncBudget(now time.Time, invokeDeadline time.Time, eventKind string, routeSyncTimeout time.Duration) time.Duration {
	providerTimeout := defaultKindTimeout(eventKind)
	if !invokeDeadline.IsZero() {
		d := invokeDeadline.Sub(now)
		if d <= 0 {
			return 0
		}
		providerTimeout = d
	}
	budget := providerTimeout - providerTimeout/timeoutMarginDenom
	if budget <= 0 {
		return 0
	}
	if routeSyncTimeout > 0 && routeSyncTimeout < budget {
		return routeSyncTimeout
	}
	return budget
}

func defaultKindTimeout(eventKind string) time.Duration {
	switch agenthooks.EventKind(eventKind) {
	case agenthooks.KindToolPre, agenthooks.KindPromptSubmitted:
		return kindTimeoutLong
	default:
		return kindTimeoutShort
	}
}
