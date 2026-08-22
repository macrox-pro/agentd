package daemon_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestWriteStatus(t *testing.T) {
	t.Parallel()

	started := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		rep    daemon.StatusReport
		asJSON bool
		want   string
		check  func(t *testing.T, raw []byte)
	}{
		{
			name: "human not running",
			rep:  daemon.StatusReport{Running: false, Socket: "/tmp/agentd.sock"},
			want: "agentd: not running\n",
		},
		{
			name: "human running",
			rep: daemon.StatusReport{
				Running:    true,
				Version:    "dev",
				Generation: 3,
			},
			want: "agentd: running (version dev, generation 3)\n",
		},
		{
			name:   "json not running",
			asJSON: true,
			rep:    daemon.StatusReport{Running: false, Socket: "/tmp/agentd.sock"},
			check: func(t *testing.T, raw []byte) {
				t.Helper()
				var got map[string]any
				require.NoError(t, json.Unmarshal(raw, &got), "json")
				assert.Equal(t, false, got["running"])
				assert.Equal(t, "/tmp/agentd.sock", got["socket"])
				assert.NotContains(t, got, "version")
			},
		},
		{
			name:   "json running",
			asJSON: true,
			rep: daemon.StatusReport{
				Running:            true,
				Socket:             "/tmp/agentd.sock",
				Version:            "dev",
				StartedAt:          started,
				Generation:         2,
				Fingerprint:        "abc",
				AsyncQueueDepth:    4,
				AsyncDroppedCount:      9,
				TrajectoryDroppedCount: 2,
				CompiledRouteCount:     7,
			},
			check: func(t *testing.T, raw []byte) {
				t.Helper()
				var got map[string]any
				require.NoError(t, json.Unmarshal(raw, &got), "json")
				assert.Equal(t, true, got["running"])
				assert.Equal(t, "dev", got["version"])
				assert.Equal(t, "2026-08-20T12:00:00Z", got["started_at"])
				assert.Equal(t, float64(2), got["generation"])
				assert.Equal(t, "abc", got["fingerprint"])
				assert.Equal(t, float64(4), got["async_queue_depth"])
				assert.Equal(t, float64(9), got["async_dropped_count"])
				assert.Equal(t, float64(2), got["trajectory_dropped_count"])
				assert.Equal(t, float64(7), got["compiled_route_count"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := daemon.WriteStatus(&buf, tt.rep, tt.asJSON)
			require.NoError(t, err, "WriteStatus(%s)", tt.name)
			if tt.check != nil {
				tt.check(t, buf.Bytes())
				return
			}
			assert.Equal(t, tt.want, buf.String(), "WriteStatus(%s)", tt.name)
		})
	}
}
