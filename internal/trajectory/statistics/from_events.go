package statistics

import (
	"encoding/json"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics/extract"
)

// FromEvents folds ledger events into session statistics.
func FromEvents(events []trajectory.Event) Session {
	out := newSession()
	if len(events) == 0 {
		return out
	}
	out.Provider = events[0].Provider
	out.SessionID = events[0].SessionID
	for _, ev := range events {
		out.EventCount++
		out.EventsByType[ev.Type]++
		out.EventsBySource[ev.Source]++
		if out.FirstSeq == 0 || ev.Seq < out.FirstSeq {
			out.FirstSeq = ev.Seq
		}
		if ev.Seq > out.LastSeq {
			out.LastSeq = ev.Seq
		}
		if out.FirstTS.IsZero() || ev.TS.Before(out.FirstTS) {
			out.FirstTS = ev.TS
		}
		if ev.TS.After(out.LastTS) {
			out.LastTS = ev.TS
		}
		switch ev.Type {
		case trajectory.TypeHookInvoked:
			var data trajectory.HookInvokedData
			if err := json.Unmarshal(ev.Data, &data); err != nil {
				continue
			}
			out.HooksByKind[data.Kind]++
			if len(ev.Raw) > 0 {
				prov := agentdv1.Provider_PROVIDER_UNSPECIFIED
				if id, ok := provider.Lookup(ev.Provider); ok {
					if p, err := id.Proto(); err == nil {
						prov = p
					}
				}
				tokens := extract.Tokens(prov, ev.Raw)
				if !tokens.Any() {
					tokens = extract.TokensFromTranscript(prov, ev.Raw)
				}
				applySessionTokens(&out, billingTokensForRollup(tokens))
				applySessionTokens(&out, contextTokensForRollup(tokens))
			}
		case trajectory.TypeHookDecided:
			var data trajectory.HookDecidedData
			if err := json.Unmarshal(ev.Data, &data); err != nil {
				continue
			}
			out.DecisionsByKind[data.Decision]++
		case trajectory.TypeAsyncDispatched:
			var data trajectory.AsyncDispatchedData
			if err := json.Unmarshal(ev.Data, &data); err != nil {
				continue
			}
			out.AsyncDispatched += uint64(data.Count)
		case trajectory.TypeAsyncDropped:
			out.AsyncDropped++
		case trajectory.TypeTranscriptMessage:
			out.TranscriptMessages++
		case trajectory.TypeTranscriptThinking:
			out.TranscriptThinking++
		}
	}
	return out
}

func applySessionTokens(s *Session, t extract.TokenFields) {
	if t.HasInput {
		s.InputTokensTotal += t.Input
	}
	if t.HasOutput {
		s.OutputTokensTotal += t.Output
	}
	if t.HasCacheRead {
		s.CacheReadTokensTotal += t.CacheRead
	}
	if t.HasCacheWrite {
		s.CacheWriteTokensTotal += t.CacheWrite
	}
	if t.HasContext {
		s.ContextTokensLast = t.Context
	}
}
