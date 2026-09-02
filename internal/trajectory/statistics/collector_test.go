package statistics_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/decision"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

func enabledSnap(stats bool) *config.Snapshot {
	return &config.Snapshot{
		Trajectory: config.TrajectoryConfig{Enabled: true, Statistics: stats},
	}
}

func sampleRecordInput(snap *config.Snapshot) trajectory.RecordInput {
	return trajectory.RecordInput{
		Provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		Result: dispatch.InvokeResult{
			Meta: dispatch.InvokeMeta{EventKind: "tool.pre"},
			Decision: &agentdv1.Decision{
				Kind: agentdv1.DecisionKind_DECISION_KIND_ALLOW,
			},
			AsyncDispatchedCount: 2,
		},
		Snap: snap,
	}
}

func cursorStopInput(snap *config.Snapshot, sessionID, raw string) trajectory.RecordInput {
	return trajectory.RecordInput{
		Provider: agentdv1.Provider_PROVIDER_CURSOR,
		RawPayload: []byte(raw),
		Result: dispatch.InvokeResult{
			Meta: dispatch.InvokeMeta{
				EventKind: "agent.stop",
				SessionID: sessionID,
			},
			Decision: decision.Neutral(),
		},
		Snap: snap,
	}
}

func codexStopInput(t *testing.T, snap *config.Snapshot, sessionID, transcriptPath string) trajectory.RecordInput {
	t.Helper()
	raw := fmt.Sprintf(`{"hook_event_name":"Stop","session_id":%q,"transcript_path":%q}`, sessionID, transcriptPath)
	return trajectory.RecordInput{
		Provider:   agentdv1.Provider_PROVIDER_CODEX,
		RawPayload: []byte(raw),
		Result: dispatch.InvokeResult{
			Meta: dispatch.InvokeMeta{
				EventKind: "agent.stop",
				SessionID: sessionID,
			},
			Decision: decision.Neutral(),
		},
		Snap: snap,
	}
}

func writeCodexRollout(t *testing.T, input, cached, output uint64) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "rollout.jsonl")
	line := fmt.Sprintf(`{"type":"event_msg","payload":{"type":"token_count","info":{"last_token_usage":{"input_tokens":%d,"cached_input_tokens":%d,"cache_write_input_tokens":0,"output_tokens":%d}}}}`,
		input, cached, output)
	require.NoError(t, os.WriteFile(path, []byte(line+"\n"), 0o600))
	return path
}

func TestObserve(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		setup func(t *testing.T) (*statistics.Collector, trajectory.RecordInput)
		check func(t *testing.T, c *statistics.Collector)
	}{
		{
			name: "skips_when_disabled",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return statistics.NewCollector(), sampleRecordInput(&config.Snapshot{
					Trajectory: config.TrajectoryConfig{Enabled: false, Statistics: true},
				})
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Equal(t, uint64(0), c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).HooksByKind[agentdv1.EventKind_EVENT_KIND_TOOL_PRE])
			},
		},
		{
			name: "skips_when_statistics_off",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return statistics.NewCollector(), sampleRecordInput(enabledSnap(false))
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Empty(t, c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).HooksByKind)
			},
		},
		{
			name: "increments_hook_and_decision",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return statistics.NewCollector(), sampleRecordInput(enabledSnap(true))
			},
			check: func(t *testing.T, c *statistics.Collector) {
				r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
				assert.Equal(t, uint64(1), r.HooksByKind[agentdv1.EventKind_EVENT_KIND_TOOL_PRE])
				assert.Equal(t, uint64(1), r.DecisionsByKind[agentdv1.DecisionKind_DECISION_KIND_ALLOW])
			},
		},
		{
			name: "async_dispatched",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return statistics.NewCollector(), sampleRecordInput(enabledSnap(true))
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Equal(t, uint64(2), c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).AsyncDispatched)
			},
		},
		{
			name: "nil_collector",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return nil, sampleRecordInput(enabledSnap(true))
			},
			check: func(t *testing.T, _ *statistics.Collector) {},
		},
		{
			name: "many_sessions_same_counters",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				in := sampleRecordInput(enabledSnap(true))
				in.Result.Meta.SessionID = "a"
				return c, in
			},
			check: func(t *testing.T, c *statistics.Collector) {
				in := sampleRecordInput(enabledSnap(true))
				in.Result.Meta.SessionID = "b"
				c.Observe(in)
				assert.Equal(t, uint64(2), c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).HooksByKind[agentdv1.EventKind_EVENT_KIND_TOOL_PRE])
			},
		},
		{
			name: "counts_neutral_no_route",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				in := sampleRecordInput(enabledSnap(true))
				in.Result.Decision = decision.Neutral()
				in.Result.Meta.HasRoute = false
				return c, in
			},
			check: func(t *testing.T, c *statistics.Collector) {
				r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
				assert.Equal(t, uint64(1), r.HooksByKind[agentdv1.EventKind_EVENT_KIND_TOOL_PRE])
				assert.Equal(t, uint64(1), r.DecisionsByKind[agentdv1.DecisionKind_DECISION_KIND_NO_DECISION])
			},
		},
		{
			name: "malformed_json_still_counts_hooks",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				in := sampleRecordInput(enabledSnap(true))
				in.RawPayload = []byte(`{not json`)
				return c, in
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Equal(t, uint64(1), c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).HooksByKind[agentdv1.EventKind_EVENT_KIND_TOOL_PRE])
			},
		},
		{
			name: "partial_token_fields",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				in := sampleRecordInput(enabledSnap(true))
				in.RawPayload = []byte(`{"usage":{"input_tokens":4}}`)
				return c, in
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					return c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).InputTokensTotal == 4
				}, time.Second, 5*time.Millisecond)
				r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
				assert.Equal(t, uint64(4), r.InputTokensTotal)
				assert.Equal(t, uint64(0), r.OutputTokensTotal)
			},
		},
		{
			name: "cursor_stop_single",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return statistics.NewCollector(), cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":19582,"output_tokens":92,"cache_read_tokens":6272,"cache_write_tokens":0}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_CURSOR)
					return r.InputTokensTotal == 19582 && r.OutputTokensTotal == 92 && r.CacheReadTokens == 6272
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_two_stops_same_session",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":100,"output_tokens":10}`))
				return c, cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":250,"output_tokens":25}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 350 && r.OutputTokensTotal == 35
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_two_sessions",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "a", `{"input_tokens":10,"output_tokens":1}`))
				return c, cursorStopInput(enabledSnap(true), "b", `{"input_tokens":20,"output_tokens":2}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 30 && r.OutputTokensTotal == 3
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_regression_same_session",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":250,"output_tokens":25}`))
				return c, cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":100,"output_tokens":10}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 350 && r.OutputTokensTotal == 35
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_no_token_fields",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return statistics.NewCollector(), cursorStopInput(enabledSnap(true), "s1", `{"status":"completed","loop_count":0}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				time.Sleep(20 * time.Millisecond)
				r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
				assert.Equal(t, uint64(1), r.HooksByKind[agentdv1.EventKind_EVENT_KIND_AGENT_STOP])
				assert.Equal(t, uint64(0), r.InputTokensTotal)
			},
		},
		{
			name: "cursor_precompact_context_last",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				in := cursorStopInput(enabledSnap(true), "s1", `{"context_tokens":120000}`)
				in.Result.Meta.EventKind = "preCompact"
				return c, in
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					return c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).ContextTokensLast == 120000
				}, time.Second, 5*time.Millisecond)
				r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
				assert.Equal(t, uint64(0), r.InputTokensTotal)
			},
		},
		{
			name: "cursor_empty_session_id",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "", `{"input_tokens":10,"output_tokens":1}`))
				return c, cursorStopInput(enabledSnap(true), "", `{"input_tokens":20,"output_tokens":2}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 30 && r.OutputTokensTotal == 3
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "codex_stop_transcript_fallback",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				path := writeCodexRollout(t, 15156, 4352, 100)
				return statistics.NewCollector(), codexStopInput(t, enabledSnap(true), "s1", path)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					global := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					byCodex := c.Snapshot(agentdv1.Provider_PROVIDER_CODEX)
					return global.InputTokensTotal == 15156 && global.OutputTokensTotal == 100 &&
						byCodex.InputTokensTotal == 15156 && byCodex.CacheReadTokens == 4352
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "codex_stop_hook_raw_wins",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				path := writeCodexRollout(t, 999, 0, 999)
				raw := fmt.Sprintf(`{"hook_event_name":"Stop","session_id":"s1","transcript_path":%q,"usage":{"input_tokens":7,"output_tokens":3}}`, path)
				return statistics.NewCollector(), trajectory.RecordInput{
					Provider:   agentdv1.Provider_PROVIDER_CODEX,
					RawPayload: []byte(raw),
					Result: dispatch.InvokeResult{
						Meta:     dispatch.InvokeMeta{EventKind: "agent.stop", SessionID: "s1"},
						Decision: decision.Neutral(),
					},
					Snap: enabledSnap(true),
				}
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_CODEX)
					return r.InputTokensTotal == 7 && r.OutputTokensTotal == 3
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_two_generations",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":280027,"output_tokens":1523,"cache_read_tokens":270592}`))
				return c, cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":461839,"output_tokens":2813,"cache_read_tokens":439680}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					global := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					byCursor := c.Snapshot(agentdv1.Provider_PROVIDER_CURSOR)
					return global.InputTokensTotal == 741866 && global.OutputTokensTotal == 4336 && global.CacheReadTokens == 710272 &&
						byCursor.InputTokensTotal == 741866 && byCursor.OutputTokensTotal == 4336 && byCursor.CacheReadTokens == 710272
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_identical_payloads",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":100,"output_tokens":10}`))
				return c, cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":100,"output_tokens":10}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 200 && r.OutputTokensTotal == 20
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_partial_second",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":100,"output_tokens":10}`))
				return c, cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":42}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 142 && r.OutputTokensTotal == 10
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "cursor_stop_precompact_interleaved",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":100,"output_tokens":10}`))
				preCompact := cursorStopInput(enabledSnap(true), "s1", `{"context_tokens":120000}`)
				preCompact.Result.Meta.EventKind = "preCompact"
				c.Observe(preCompact)
				return c, cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":250,"output_tokens":25}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 350 && r.OutputTokensTotal == 35 && r.ContextTokensLast == 120000
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "codex_two_stops_same_session",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				pathA := writeCodexRollout(t, 10, 0, 1)
				c.Observe(codexStopInput(t, enabledSnap(true), "s1", pathA))
				pathB := writeCodexRollout(t, 20, 0, 2)
				return c, codexStopInput(t, enabledSnap(true), "s1", pathB)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_CODEX)
					return r.InputTokensTotal == 30 && r.OutputTokensTotal == 3
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "codex_two_sessions",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				pathA := writeCodexRollout(t, 10, 0, 1)
				c.Observe(codexStopInput(t, enabledSnap(true), "a", pathA))
				pathB := writeCodexRollout(t, 20, 0, 2)
				return c, codexStopInput(t, enabledSnap(true), "b", pathB)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_CODEX)
					return r.InputTokensTotal == 30 && r.OutputTokensTotal == 3
				}, time.Second, 5*time.Millisecond)
			},
		},
		{
			name: "codex_stop_no_transcript_path",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				return statistics.NewCollector(), trajectory.RecordInput{
					Provider:   agentdv1.Provider_PROVIDER_CODEX,
					RawPayload: []byte(`{"hook_event_name":"Stop","session_id":"s1"}`),
					Result: dispatch.InvokeResult{
						Meta:     dispatch.InvokeMeta{EventKind: "agent.stop", SessionID: "s1"},
						Decision: decision.Neutral(),
					},
					Snap: enabledSnap(true),
				}
			},
			check: func(t *testing.T, c *statistics.Collector) {
				time.Sleep(20 * time.Millisecond)
				r := c.Snapshot(agentdv1.Provider_PROVIDER_CODEX)
				assert.Equal(t, uint64(1), r.HooksByKind[agentdv1.EventKind_EVENT_KIND_AGENT_STOP])
				assert.Equal(t, uint64(0), r.InputTokensTotal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, in := tt.setup(t)
			if c != nil {
				c.Observe(in)
				tt.check(t, c)
				return
			}
			var nilCollector *statistics.Collector
			nilCollector.Observe(in)
			tt.check(t, nil)
		})
	}
}

func TestObserve_reload(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		on   bool
	}{
		{name: "statistics_on", on: true},
		{name: "statistics_off", on: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := statistics.NewCollector()
			c.Observe(sampleRecordInput(enabledSnap(true)))
			c.Observe(sampleRecordInput(enabledSnap(tt.on)))
			r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
			if tt.on {
				require.Equal(t, uint64(2), r.HooksByKind[agentdv1.EventKind_EVENT_KIND_TOOL_PRE])
			} else {
				require.Equal(t, uint64(1), r.HooksByKind[agentdv1.EventKind_EVENT_KIND_TOOL_PRE])
			}
		})
	}
}
