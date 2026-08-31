package statistics

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

type statisticsResponseJSON struct {
	Since  time.Time  `json:"since"`
	Rollup rollupJSON `json:"rollup"`
}

type rollupJSON struct {
	HooksByKind           map[string]uint64 `json:"hooks_by_kind"`
	DecisionsByKind       map[string]uint64 `json:"decisions_by_kind"`
	AsyncDispatchedTotal  uint64            `json:"async_dispatched_total"`
	InputTokensTotal      uint64            `json:"input_tokens_total"`
	OutputTokensTotal     uint64            `json:"output_tokens_total"`
	CacheReadTokensTotal  uint64            `json:"cache_read_tokens_total"`
	CacheWriteTokensTotal uint64            `json:"cache_write_tokens_total"`
	ContextTokensLast     uint64            `json:"context_tokens_last"`
}

// WriteRollup formats daemon statistics RPC response for CLI output.
func WriteRollup(w io.Writer, resp *agentdv1.StatisticsResponse, jsonOut bool) error {
	if resp == nil {
		return fmt.Errorf("nil statistics response")
	}
	if jsonOut {
		out := statisticsResponseJSON{
			Since:  time.Time{},
			Rollup: rollupJSON{HooksByKind: map[string]uint64{}, DecisionsByKind: map[string]uint64{}},
		}
		if ts := resp.GetSince(); ts != nil {
			out.Since = ts.AsTime().UTC()
		}
		if r := resp.GetRollup(); r != nil {
			out.Rollup = rollupToJSON(r)
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	since := ""
	if ts := resp.GetSince(); ts != nil {
		since = ts.AsTime().UTC().Format(time.RFC3339Nano)
	}
	r := resp.GetRollup()
	if r == nil {
		fmt.Fprintf(w, "since=%s\n", since)
		return nil
	}
	fmt.Fprintf(w, "since=%s\n", since)
	fmt.Fprintf(w, "async_dispatched_total=%d\n", r.GetAsyncDispatchedTotal())
	fmt.Fprintf(w, "input_tokens_total=%d\n", r.GetInputTokensTotal())
	fmt.Fprintf(w, "output_tokens_total=%d\n", r.GetOutputTokensTotal())
	fmt.Fprintf(w, "cache_read_tokens_total=%d\n", r.GetCacheReadTokensTotal())
	fmt.Fprintf(w, "cache_write_tokens_total=%d\n", r.GetCacheWriteTokensTotal())
	fmt.Fprintf(w, "context_tokens_last=%d\n", r.GetContextTokensLast())
	writeSortedUint64Map(w, "hooks_by_kind", r.GetHooksByKind(), eventKindName)
	writeSortedUint64Map(w, "decisions_by_kind", r.GetDecisionsByKind(), decisionKindName)
	return nil
}

// WriteSession formats offline session stats for CLI output.
func WriteSession(w io.Writer, s Session, jsonOut bool) error {
	if jsonOut {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(s)
	}
	fmt.Fprintf(w, "provider=%s session_id=%s events=%d\n", s.Provider, s.SessionID, s.EventCount)
	if !s.FirstTS.IsZero() {
		fmt.Fprintf(w, "first_ts=%s last_ts=%s\n", s.FirstTS.UTC().Format(time.RFC3339), s.LastTS.UTC().Format(time.RFC3339))
	}
	fmt.Fprintf(w, "async_dispatched=%d async_dropped=%d\n", s.AsyncDispatched, s.AsyncDropped)
	fmt.Fprintf(w, "transcript_messages=%d transcript_thinking=%d\n", s.TranscriptMessages, s.TranscriptThinking)
	fmt.Fprintf(w, "input_tokens_total=%d output_tokens_total=%d\n", s.InputTokensTotal, s.OutputTokensTotal)
	fmt.Fprintf(w, "cache_read_tokens_total=%d cache_write_tokens_total=%d context_tokens_last=%d\n",
		s.CacheReadTokensTotal, s.CacheWriteTokensTotal, s.ContextTokensLast)
	writeSortedStringMap(w, "hooks_by_kind", s.HooksByKind)
	writeSortedStringMap(w, "decisions_by_kind", s.DecisionsByKind)
	writeSortedStringMap(w, "events_by_type", s.EventsByType)
	writeSortedStringMap(w, "events_by_source", s.EventsBySource)
	return nil
}

// RollupToProto maps domain rollup counters to the wire message.
func RollupToProto(r StatisticsRollup) *agentdv1.StatisticsRollup {
	out := &agentdv1.StatisticsRollup{
		HooksByKind:           map[int32]uint64{},
		DecisionsByKind:       map[int32]uint64{},
		AsyncDispatchedTotal:  r.AsyncDispatched,
		InputTokensTotal:      r.InputTokensTotal,
		OutputTokensTotal:     r.OutputTokensTotal,
		CacheReadTokensTotal:  r.CacheReadTokens,
		CacheWriteTokensTotal: r.CacheWriteTokens,
		ContextTokensLast:     r.ContextTokensLast,
	}
	for k, v := range r.HooksByKind {
		out.HooksByKind[int32(k)] = v
	}
	for k, v := range r.DecisionsByKind {
		out.DecisionsByKind[int32(k)] = v
	}
	return out
}

// Response builds a StatisticsResponse from daemon start time and rollup.
func Response(since time.Time, r StatisticsRollup) *agentdv1.StatisticsResponse {
	return &agentdv1.StatisticsResponse{
		Since:  timestamppb.New(since),
		Rollup: RollupToProto(r),
	}
}

func rollupToJSON(r *agentdv1.StatisticsRollup) rollupJSON {
	out := rollupJSON{
		HooksByKind:           map[string]uint64{},
		DecisionsByKind:       map[string]uint64{},
		AsyncDispatchedTotal:  r.GetAsyncDispatchedTotal(),
		InputTokensTotal:      r.GetInputTokensTotal(),
		OutputTokensTotal:     r.GetOutputTokensTotal(),
		CacheReadTokensTotal:  r.GetCacheReadTokensTotal(),
		CacheWriteTokensTotal: r.GetCacheWriteTokensTotal(),
		ContextTokensLast:     r.GetContextTokensLast(),
	}
	for k, v := range r.GetHooksByKind() {
		out.HooksByKind[eventKindName(k)] = v
	}
	for k, v := range r.GetDecisionsByKind() {
		out.DecisionsByKind[decisionKindName(k)] = v
	}
	return out
}

func writeSortedUint64Map(w io.Writer, prefix string, m map[int32]uint64, name func(int32) string) {
	if len(m) == 0 {
		return
	}
	keys := make([]int32, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b int32) int {
		return strings.Compare(name(a), name(b))
	})
	for _, k := range keys {
		fmt.Fprintf(w, "%s[%s]=%d\n", prefix, name(k), m[k])
	}
}

func writeSortedStringMap(w io.Writer, prefix string, m map[string]uint64) {
	if len(m) == 0 {
		return
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s[%s]=%d\n", prefix, k, m[k])
	}
}

func eventKindName(v int32) string {
	return strings.TrimPrefix(agentdv1.EventKind(v).String(), "EVENT_KIND_")
}

func decisionKindName(v int32) string {
	return strings.TrimPrefix(agentdv1.DecisionKind(v).String(), "DECISION_KIND_")
}
