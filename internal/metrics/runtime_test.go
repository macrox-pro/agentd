package metrics_test

import (
	"fmt"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/metrics"
)

func TestRegisterRuntime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		stats    metrics.RuntimeStats
		wantName string
		wantVal  float64
		kind     dto.MetricType
	}{
		{
			name: "queue depth",
			stats: metrics.RuntimeStats{
				AsyncQueueDepth: func() float64 { return 7 },
			},
			wantName: "agentd_async_queue_depth",
			wantVal:  7,
			kind:     dto.MetricType_GAUGE,
		},
		{
			name: "dropped total",
			stats: metrics.RuntimeStats{
				AsyncQueueDropped: func() float64 { return 42 },
			},
			wantName: "agentd_async_queue_dropped_total",
			wantVal:  42,
			kind:     dto.MetricType_COUNTER,
		},
		{
			name:     "nil closures",
			stats:    metrics.RuntimeStats{},
			wantName: "agentd_async_queue_depth",
			wantVal:  0,
			kind:     dto.MetricType_GAUGE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reg := metrics.NewRegistry()
			if tt.name == "nil closures" {
				metrics.RegisterRuntime(reg, tt.stats)
				_, err := gatherMetric(reg, tt.wantName)
				require.Error(t, err, "TestRegisterRuntime(%q)", tt.name)
				return
			}
			metrics.RegisterRuntime(reg, tt.stats)
			m, err := gatherMetric(reg, tt.wantName)
			require.NoError(t, err, "TestRegisterRuntime(%q)", tt.name)
			require.NotNil(t, m, "TestRegisterRuntime(%q)", tt.name)
			switch tt.kind {
			case dto.MetricType_GAUGE:
				require.NotNil(t, m.Gauge, "TestRegisterRuntime(%q)", tt.name)
				assert.Equal(t, tt.wantVal, m.GetGauge().GetValue(), "TestRegisterRuntime(%q)", tt.name)
			case dto.MetricType_COUNTER:
				require.NotNil(t, m.Counter, "TestRegisterRuntime(%q)", tt.name)
				assert.Equal(t, tt.wantVal, m.GetCounter().GetValue(), "TestRegisterRuntime(%q)", tt.name)
			}
		})
	}
}

func gatherMetric(reg *prometheus.Registry, name string) (*dto.Metric, error) {
	families, err := reg.Gather()
	if err != nil {
		return nil, err
	}
	for _, fam := range families {
		if fam.GetName() != name {
			continue
		}
		if len(fam.Metric) == 0 {
			return nil, nil
		}
		return fam.Metric[0], nil
	}
		return nil, fmt.Errorf("not found")
}
