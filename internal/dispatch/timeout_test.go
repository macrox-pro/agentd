package dispatch_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestSyncBudget(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name         string
		deadline     time.Time
		kind         string
		routeTimeout time.Duration
		want         time.Duration
	}{
		{
			name: "tool.pre default with margin",
			kind: "tool.pre",
			want: 27 * time.Second, // 30s - 10%
		},
		{
			name: "short kind default with margin",
			kind: "agent.stop",
			want: time.Duration(5*time.Second) * 9 / 10,
		},
		{
			name: "subagent.start short budget",
			kind: "subagent.start",
			want: time.Duration(5*time.Second) * 9 / 10,
		},
		{
			name: "compact.pre short budget",
			kind: "compact.pre",
			want: time.Duration(5*time.Second) * 9 / 10,
		},
		{
			name: "prompt.submitted long budget",
			kind: "prompt.submitted",
			want: 27 * time.Second,
		},
		{
			name:     "invoke deadline with margin",
			deadline: now.Add(10 * time.Second),
			kind:     "tool.pre",
			want:     9 * time.Second,
		},
		{
			name:         "route sync_timeout caps budget",
			deadline:     now.Add(20 * time.Second),
			kind:         "tool.pre",
			routeTimeout: 5 * time.Second,
			want:         5 * time.Second,
		},
		{
			name:         "route larger than budget unused",
			deadline:     now.Add(10 * time.Second),
			kind:         "tool.pre",
			routeTimeout: 30 * time.Second,
			want:         9 * time.Second,
		},
		{
			name:     "expired deadline",
			deadline: now.Add(-time.Second),
			kind:     "tool.pre",
			want:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := dispatch.SyncBudget(now, tt.deadline, tt.kind, tt.routeTimeout)
			assert.Equal(t, tt.want, got, "SyncBudget(%v, %q, %v)", tt.deadline, tt.kind, tt.routeTimeout)
		})
	}
}
