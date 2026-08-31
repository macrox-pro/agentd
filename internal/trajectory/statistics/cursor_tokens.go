package statistics

import (
	"github.com/macrox-pro/agentd/internal/trajectory/statistics/extract"
)

func computeCursorBillingDelta(last map[string]extract.TokenFields, sessionID string, t extract.TokenFields) extract.TokenFields {
	if !t.HasBilling() {
		return extract.TokenFields{}
	}
	if sessionID == "" {
		return t
	}
	prev := last[sessionID]
	var delta extract.TokenFields
	if t.HasInput {
		delta.Input = billingDelta(t.Input, prev.Input)
		delta.HasInput = true
	}
	if t.HasOutput {
		delta.Output = billingDelta(t.Output, prev.Output)
		delta.HasOutput = true
	}
	if t.HasCacheRead {
		delta.CacheRead = billingDelta(t.CacheRead, prev.CacheRead)
		delta.HasCacheRead = true
	}
	if t.HasCacheWrite {
		delta.CacheWrite = billingDelta(t.CacheWrite, prev.CacheWrite)
		delta.HasCacheWrite = true
	}
	last[sessionID] = t
	return delta
}

func billingDelta(current, prev uint64) uint64 {
	if current < prev {
		return 0
	}
	return current - prev
}
