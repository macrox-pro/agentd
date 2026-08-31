package statistics_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

func TestFromEvents(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	tests := []struct {
		name         string
		events       []trajectory.Event
		buildEvents  func(t *testing.T) []trajectory.Event
		check        func(t *testing.T, s statistics.Session)
	}{
		{
			name:   "empty_file",
			events: nil,
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(0), s.EventCount)
			},
		},
		{
			name: "hooks_and_decisions",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "claude-code", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"tool.pre"}`)},
				{Seq: 2, Type: trajectory.TypeHookDecided, Source: trajectory.SourceDecision, Provider: "claude-code", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"tool.pre","decision":"DECISION_KIND_ALLOW"}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(2), s.EventCount)
				assert.Equal(t, uint64(1), s.HooksByKind["tool.pre"])
				assert.Equal(t, uint64(1), s.DecisionsByKind["DECISION_KIND_ALLOW"])
			},
		},
		{
			name: "transcript_counts",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeTranscriptMessage, Source: trajectory.SourceTranscript, Provider: "claude-code", SessionID: "s1", TS: now},
				{Seq: 2, Type: trajectory.TypeTranscriptThinking, Source: trajectory.SourceTranscript, Provider: "claude-code", SessionID: "s1", TS: now},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(1), s.TranscriptMessages)
				assert.Equal(t, uint64(1), s.TranscriptThinking)
			},
		},
		{
			name: "transcript_only",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeTranscriptMessage, Source: trajectory.SourceTranscript, Provider: "claude-code", SessionID: "s1", TS: now},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Empty(t, s.HooksByKind)
				assert.Equal(t, uint64(1), s.TranscriptMessages)
			},
		},
		{
			name: "tokens_from_raw",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "claude-code", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"tool.pre"}`),
					Raw:  json.RawMessage(`{"usage":{"input_tokens":9,"output_tokens":1}}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(9), s.InputTokensTotal)
				assert.Equal(t, uint64(1), s.OutputTokensTotal)
			},
		},
		{
			name: "no_raw_skips_tokens",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "claude-code", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"tool.pre"}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(0), s.InputTokensTotal)
			},
		},
		{
			name: "bad_data_skipped",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "claude-code", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{`)},
				{Seq: 2, Type: trajectory.TypeHookDecided, Source: trajectory.SourceDecision, Provider: "claude-code", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"decision":"x"}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(2), s.EventCount)
				assert.Equal(t, uint64(1), s.DecisionsByKind["x"])
			},
		},
		{
			name: "ignorable_counted",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeAsyncDropped, Source: trajectory.SourceSystem, Provider: "claude-code", SessionID: "s1", TS: now, Ignorable: true},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(1), s.EventsByType[trajectory.TypeAsyncDropped])
			},
		},
		{
			name: "cursor_stop_two_stops_from_events",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":100,"output_tokens":10}`)},
				{Seq: 2, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":250,"output_tokens":25}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(350), s.InputTokensTotal)
				assert.Equal(t, uint64(35), s.OutputTokensTotal)
			},
		},
		{
			name: "cursor_stop_two_sessions_from_events",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "a", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":10,"output_tokens":1}`)},
				{Seq: 2, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "b", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":20,"output_tokens":2}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(30), s.InputTokensTotal)
				assert.Equal(t, uint64(3), s.OutputTokensTotal)
			},
		},
		{
			name: "cursor_stop_regression_from_events",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":250,"output_tokens":25}`)},
				{Seq: 2, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":100,"output_tokens":10}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(350), s.InputTokensTotal)
				assert.Equal(t, uint64(35), s.OutputTokensTotal)
			},
		},
		{
			name: "cursor_stop_identical_from_events",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":100,"output_tokens":10}`)},
				{Seq: 2, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":100,"output_tokens":10}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(200), s.InputTokensTotal)
				assert.Equal(t, uint64(20), s.OutputTokensTotal)
			},
		},
		{
			name: "cursor_stop_partial_from_events",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":100,"output_tokens":10}`)},
				{Seq: 2, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":42}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(142), s.InputTokensTotal)
				assert.Equal(t, uint64(10), s.OutputTokensTotal)
			},
		},
		{
			name: "cursor_stop_precompact_from_events",
			events: []trajectory.Event{
				{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":100,"output_tokens":10}`)},
				{Seq: 2, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"preCompact"}`),
					Raw:  json.RawMessage(`{"context_tokens":120000}`)},
				{Seq: 3, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor", SessionID: "s1", TS: now,
					Data: json.RawMessage(`{"kind":"agent.stop"}`),
					Raw:  json.RawMessage(`{"input_tokens":250,"output_tokens":25}`)},
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(350), s.InputTokensTotal)
				assert.Equal(t, uint64(35), s.OutputTokensTotal)
				assert.Equal(t, uint64(120000), s.ContextTokensLast)
			},
		},
		{
			name: "offline_from_events_codex_fallback",
			buildEvents: func(t *testing.T) []trajectory.Event {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "rollout.jsonl")
				line := `{"type":"event_msg","payload":{"type":"token_count","info":{"last_token_usage":{"input_tokens":15156,"cached_input_tokens":4352,"cache_write_input_tokens":0,"output_tokens":100}}}}`
				require.NoError(t, os.WriteFile(path, []byte(line+"\n"), 0o600))
				raw := fmt.Sprintf(`{"hook_event_name":"Stop","session_id":"s1","transcript_path":%q}`, path)
				return []trajectory.Event{
					{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "codex", SessionID: "s1", TS: now,
						Data: json.RawMessage(`{"kind":"agent.stop"}`),
						Raw:  json.RawMessage(raw)},
				}
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(15156), s.InputTokensTotal)
				assert.Equal(t, uint64(100), s.OutputTokensTotal)
				assert.Equal(t, uint64(4352), s.CacheReadTokensTotal)
			},
		},
		{
			name: "offline_no_raw_skips_transcript",
			buildEvents: func(t *testing.T) []trajectory.Event {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "rollout.jsonl")
				line := `{"type":"event_msg","payload":{"type":"token_count","info":{"last_token_usage":{"input_tokens":999,"output_tokens":1}}}}`
				require.NoError(t, os.WriteFile(path, []byte(line+"\n"), 0o600))
				return []trajectory.Event{
					{Seq: 1, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "codex", SessionID: "s1", TS: now,
						Data: json.RawMessage(`{"kind":"agent.stop"}`)},
				}
			},
			check: func(t *testing.T, s statistics.Session) {
				assert.Equal(t, uint64(0), s.InputTokensTotal)
				assert.Equal(t, uint64(1), s.HooksByKind["agent.stop"])
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			events := tt.events
			if tt.buildEvents != nil {
				events = tt.buildEvents(t)
			}
			got := statistics.FromEvents(events)
			tt.check(t, got)
		})
	}
}
