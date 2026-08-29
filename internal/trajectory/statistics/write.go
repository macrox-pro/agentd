package statistics

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

var protoJSON = protojson.MarshalOptions{
	EmitUnpopulated: false,
	Indent:          "  ",
	UseProtoNames:   true,
}

// WriteRollup formats daemon statistics RPC response for CLI output.
func WriteRollup(w io.Writer, resp *agentdv1.StatisticsResponse, jsonOut bool) error {
	if resp == nil {
		return fmt.Errorf("nil statistics response")
	}
	if jsonOut {
		b, err := protoJSON.Marshal(resp)
		if err != nil {
			return err
		}
		_, err = w.Write(append(b, '\n'))
		return err
	}
	since := ""
	if ts := resp.GetSince(); ts != nil {
		since = ts.AsTime().UTC().Format(time.RFC3339)
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
	for k, v := range r.GetHooksByKind() {
		fmt.Fprintf(w, "hooks_by_kind[%s]=%d\n", eventKindName(k), v)
	}
	for k, v := range r.GetDecisionsByKind() {
		fmt.Fprintf(w, "decisions_by_kind[%s]=%d\n", decisionKindName(k), v)
	}
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
	for k, v := range s.HooksByKind {
		fmt.Fprintf(w, "hooks_by_kind[%s]=%d\n", k, v)
	}
	for k, v := range s.DecisionsByKind {
		fmt.Fprintf(w, "decisions_by_kind[%s]=%d\n", k, v)
	}
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

func eventKindName(v int32) string {
	return strings.TrimPrefix(agentdv1.EventKind(v).String(), "EVENT_KIND_")
}

func decisionKindName(v int32) string {
	return strings.TrimPrefix(agentdv1.DecisionKind(v).String(), "DECISION_KIND_")
}
