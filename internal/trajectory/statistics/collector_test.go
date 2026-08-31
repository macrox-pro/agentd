package statistics_test

import (
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
					return r.InputTokensTotal == 250 && r.OutputTokensTotal == 25
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
			name: "cursor_stop_regression_saturates",
			setup: func(t *testing.T) (*statistics.Collector, trajectory.RecordInput) {
				c := statistics.NewCollector()
				c.Observe(cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":250,"output_tokens":25}`))
				return c, cursorStopInput(enabledSnap(true), "s1", `{"input_tokens":100,"output_tokens":10}`)
			},
			check: func(t *testing.T, c *statistics.Collector) {
				assert.Eventually(t, func() bool {
					r := c.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED)
					return r.InputTokensTotal == 250 && r.OutputTokensTotal == 25
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
