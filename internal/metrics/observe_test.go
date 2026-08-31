package metrics_test

import (
	"strings"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/metrics"
)

func TestGatherForbiddenLabels(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		"session_id",
		"tool_name",
		"route_name",
		"cwd",
		"fingerprint",
		"url",
		"command",
	}

	reg := metrics.NewRegistry()
	metrics.RegisterGoAndProcess(reg)
	rec := metrics.NewRecorder(reg)
	metrics.RegisterReloadCounter(reg)
	metrics.RegisterRuntime(reg, metrics.RuntimeStats{
		AsyncQueueDepth:    func() float64 { return 1 },
		AsyncQueueCapacity: func() float64 { return 8 },
	})

	rec.ObserveInvoke("claude_code", "tool.pre", "DECISION_KIND_ASK", "ok", 0.01)
	rec.ObserveAsync("http", "ok", 0.02)

	families, err := reg.Gather()
	require.NoError(t, err, "Gather")

	tests := []struct {
		name  string
		label string
	}{
		{name: "session_id", label: "session_id"},
		{name: "tool_name", label: "tool_name"},
		{name: "route_name", label: "route_name"},
		{name: "cwd", label: "cwd"},
		{name: "fingerprint", label: "fingerprint"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, fam := range families {
				if !strings.HasPrefix(fam.GetName(), "agentd_") {
					continue
				}
				for _, m := range fam.GetMetric() {
					for _, lp := range labelPairs(m) {
						for _, bad := range forbidden {
							assert.NotEqual(t, bad, lp.GetName(), "forbidden label on %s", fam.GetName())
						}
					}
				}
			}
		})
	}
}

func labelPairs(m *dto.Metric) []*dto.LabelPair {
	if m == nil {
		return nil
	}
	return m.GetLabel()
}
