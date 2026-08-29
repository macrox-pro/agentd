package statistics_test

import (
	"bytes"
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
	tests := []struct {
		name     string
		jsonOut  bool
		contains string
	}{
		{name: "human", jsonOut: false, contains: "since=2026-01-02T03:04:05Z"},
		{name: "json", jsonOut: true, contains: `"since"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			require.NoError(t, statistics.WriteRollup(&buf, resp, tt.jsonOut))
			assert.Contains(t, buf.String(), tt.contains)
		})
	}
}
