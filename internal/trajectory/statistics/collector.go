package statistics

import (
	"maps"
	"sync"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics/extract"
)

// Collector holds daemon-lifetime trajectory statistics counters.
type Collector struct {
	mu         sync.Mutex
	global     StatisticsRollup
	byProvider map[agentdv1.Provider]StatisticsRollup
}

// NewCollector returns an empty statistics collector.
func NewCollector() *Collector {
	return &Collector{
		global:     newStatisticsRollup(),
		byProvider: map[agentdv1.Provider]StatisticsRollup{},
	}
}

// Observe increments counters from one successful HookService.Invoke.
func (c *Collector) Observe(in trajectory.RecordInput) {
	if c == nil || in.Snap == nil {
		return
	}
	if err := Gate(in.Snap.Trajectory); err != nil {
		return
	}
	meta := in.Result.Meta
	hookKind := HookKind(meta.EventKind)
	decisionKind := agentdv1.DecisionKind_DECISION_KIND_NO_DECISION
	if d := in.Result.Decision; d != nil {
		decisionKind = d.GetKind()
	}
	asyncN := uint64(in.Result.AsyncDispatchedCount)

	c.mu.Lock()
	applyObserve(&c.global, hookKind, decisionKind, asyncN)
	if in.Provider != agentdv1.Provider_PROVIDER_UNSPECIFIED {
		r := c.byProvider[in.Provider]
		if r.HooksByKind == nil {
			r = newStatisticsRollup()
		}
		applyObserve(&r, hookKind, decisionKind, asyncN)
		c.byProvider[in.Provider] = r
	}
	c.mu.Unlock()

	raw := append([]byte(nil), in.RawPayload...)
	prov := in.Provider
	kind := meta.EventKind
	go c.extractAsync(prov, raw, kind)
}

func applyObserve(r *StatisticsRollup, hook agentdv1.EventKind, decision agentdv1.DecisionKind, asyncN uint64) {
	r.HooksByKind[hook]++
	r.DecisionsByKind[decision]++
	r.AsyncDispatched += asyncN
}

func (c *Collector) extractAsync(prov agentdv1.Provider, raw []byte, eventKind string) {
	if c == nil || len(raw) == 0 {
		return
	}
	stop := agenthooks.EventKind(eventKind) == agenthooks.KindStop
	tokens := extract.Tokens(prov, raw)
	if !tokens.Any() && stop {
		tokens = extract.TokensFromTranscript(prov, raw)
	}
	if !tokens.Any() {
		return
	}
	var billing extract.TokenFields
	if stop {
		billing = billingTokensForRollup(tokens)
	}
	contextTokens := contextTokensForRollup(tokens)
	if !billing.HasBilling() && !contextTokens.HasContext {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	applyTokens(&c.global, billing)
	applyTokens(&c.global, contextTokens)
	if prov != agentdv1.Provider_PROVIDER_UNSPECIFIED {
		r := c.byProvider[prov]
		if r.HooksByKind == nil {
			r = newStatisticsRollup()
		}
		applyTokens(&r, billing)
		applyTokens(&r, contextTokens)
		c.byProvider[prov] = r
	}
}

func billingTokensForRollup(tokens extract.TokenFields) extract.TokenFields {
	out := tokens
	out.Context = 0
	out.HasContext = false
	return out
}

func contextTokensForRollup(tokens extract.TokenFields) extract.TokenFields {
	if !tokens.HasContext {
		return extract.TokenFields{}
	}
	return extract.TokenFields{Context: tokens.Context, HasContext: true}
}

func applyTokens(r *StatisticsRollup, t extract.TokenFields) {
	if t.HasInput {
		r.InputTokensTotal += t.Input
	}
	if t.HasOutput {
		r.OutputTokensTotal += t.Output
	}
	if t.HasCacheRead {
		r.CacheReadTokens += t.CacheRead
	}
	if t.HasCacheWrite {
		r.CacheWriteTokens += t.CacheWrite
	}
	if t.HasContext {
		r.ContextTokensLast = t.Context
	}
}

// Snapshot returns a copy of counters for global or one provider filter.
func (c *Collector) Snapshot(p agentdv1.Provider) StatisticsRollup {
	if c == nil {
		return newStatisticsRollup()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if p == agentdv1.Provider_PROVIDER_UNSPECIFIED {
		return cloneRollup(c.global)
	}
	return cloneRollup(c.byProvider[p])
}

func cloneRollup(in StatisticsRollup) StatisticsRollup {
	out := newStatisticsRollup()
	maps.Copy(out.HooksByKind, in.HooksByKind)
	maps.Copy(out.DecisionsByKind, in.DecisionsByKind)
	out.AsyncDispatched = in.AsyncDispatched
	out.InputTokensTotal = in.InputTokensTotal
	out.OutputTokensTotal = in.OutputTokensTotal
	out.CacheReadTokens = in.CacheReadTokens
	out.CacheWriteTokens = in.CacheWriteTokens
	out.ContextTokensLast = in.ContextTokensLast
	return out
}
