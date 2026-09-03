package statistics_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

func TestWriteRollup(t *testing.T) {
	t.Parallel()
	since := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	resp := statistics.Response(since, statistics.StatisticsRollup{
		HooksByKind: map[agentdv1.EventKind]uint64{agentdv1.EventKind_EVENT_KIND_TOOL_PRE: 1},
	})
	respWithKinds := statistics.Response(since, statistics.StatisticsRollup{
		HooksByKind: map[agentdv1.EventKind]uint64{
			agentdv1.EventKind_EVENT_KIND_SESSION_START:    1,
			agentdv1.EventKind_EVENT_KIND_PROMPT_SUBMITTED: 1,
		},
		InputTokensTotal: 19582,
	})
	tests := []struct {
		name     string
		resp     *agentdv1.StatisticsResponse
		jsonOut  bool
		contains []string
		not      []string
	}{
		{
			name:     "human",
			resp:     resp,
			jsonOut:  false,
			contains: []string{"since=2026-01-02T03:04:05Z"},
		},
		{
			name:     "human_sorted_hooks",
			resp:     respWithKinds,
			jsonOut:  false,
			contains: []string{"hooks_by_kind[PROMPT_SUBMITTED]=1", "hooks_by_kind[SESSION_START]=1"},
		},
		{
			name:     "json_enum_names",
			resp:     respWithKinds,
			jsonOut:  true,
			contains: []string{`"PROMPT_SUBMITTED"`, `"SESSION_START"`},
			not:      []string{`"1": 1`},
		},
		{
			name:     "json_numeric_counters",
			resp:     respWithKinds,
			jsonOut:  true,
			contains: []string{`"input_tokens_total": 19582`},
		},
		{
			name:    "json_emit_zero_scalars",
			resp:    statistics.Response(since, statistics.StatisticsRollup{}),
			jsonOut: true,
			contains: []string{
				`"input_tokens_total": 0`,
				`"output_tokens_total": 0`,
				`"cache_read_tokens_total": 0`,
				`"cache_write_tokens_total": 0`,
				`"context_tokens_last": 0`,
				`"hooks_by_kind": {}`,
				`"decisions_by_kind": {}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			require.NoError(t, statistics.WriteRollup(&buf, tt.resp, tt.jsonOut), "WriteRollup(%s)", tt.name)
			out := buf.String()
			for _, want := range tt.contains {
				assert.Contains(t, out, want, "WriteRollup(%s)", tt.name)
			}
			for _, bad := range tt.not {
				assert.NotContains(t, out, bad, "WriteRollup(%s)", tt.name)
			}
			if tt.name == "human_sorted_hooks" {
				assert.True(t, strings.Index(out, "PROMPT_SUBMITTED") < strings.Index(out, "SESSION_START"), "WriteRollup(%s)", tt.name)
			}
			if tt.jsonOut && tt.name == "json_numeric_counters" {
				var parsed map[string]any
				require.NoError(t, json.Unmarshal([]byte(out), &parsed), "WriteRollup(%s)", tt.name)
			}
		})
	}
}

func TestWriteSession(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		session  statistics.Session
		contains []string
	}{
		{
			name: "human_sorted_hooks",
			session: statistics.Session{
				Provider:  "cursor",
				SessionID: "s1",
				HooksByKind: map[string]uint64{
					"prompt.submitted": 1,
					"session.start":    1,
				},
				EventsByType: map[string]uint64{
					"hook.invoked": 2,
					"hook.decided": 1,
				},
				EventsBySource: map[string]uint64{
					"daemon": 1,
					"hook":   2,
				},
			},
			contains: []string{
				"hooks_by_kind[prompt.submitted]=1",
				"hooks_by_kind[session.start]=1",
				"events_by_type[hook.decided]=1",
				"events_by_type[hook.invoked]=2",
				"events_by_source[daemon]=1",
				"events_by_source[hook]=2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			require.NoError(t, statistics.WriteSession(&buf, tt.session, false), "WriteSession(%s)", tt.name)
			out := buf.String()
			for _, want := range tt.contains {
				assert.Contains(t, out, want, "WriteSession(%s)", tt.name)
			}
			assert.True(t, strings.Index(out, "prompt.submitted") < strings.Index(out, "session.start"), "WriteSession(%s)", tt.name)
			if tt.name == "human_sorted_hooks" {
				assert.True(t, strings.Index(out, "hook.decided") < strings.Index(out, "hook.invoked"), "WriteSession(%s) events_by_type", tt.name)
				assert.True(t, strings.Index(out, "events_by_source[daemon]") < strings.Index(out, "events_by_source[hook]"), "WriteSession(%s) events_by_source", tt.name)
			}
		})
	}
}
