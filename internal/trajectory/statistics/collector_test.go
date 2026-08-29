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
