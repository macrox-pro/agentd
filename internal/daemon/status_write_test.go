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
			rep: daemon.StatusReport{
				Running: false,
				Socket:  "/tmp/agentd.sock",
				Autostart: daemon.AutostartReport{
					Enabled: false,
					Backend: daemon.BackendLaunchd,
				},
			},
			want: "agentd: not running\n",
		},
		{
			name: "human running",
			rep: daemon.StatusReport{
				Running:    true,
				Version:    "dev",
				Generation: 3,
				Autostart:  daemon.AutostartReport{Enabled: false},
			},
			want: "agentd: running (version dev, generation 3)\n",
		},
		{
			name:   "json not running",
			asJSON: true,
			rep: daemon.StatusReport{
				Running: false,
				Socket:  "/tmp/agentd.sock",
				Autostart: daemon.AutostartReport{
					Enabled: false,
					Backend: daemon.BackendLaunchd,
				},
			},
			check: func(t *testing.T, raw []byte) {
				t.Helper()
				var got map[string]any
				require.NoError(t, json.Unmarshal(raw, &got), "json")
				assert.Equal(t, false, got["running"])
				assert.Equal(t, "/tmp/agentd.sock", got["socket"])
				assert.NotContains(t, got, "version")
				auto, ok := got["autostart"].(map[string]any)
				require.True(t, ok, "autostart")
				assert.Equal(t, false, auto["enabled"])
			},
		},
		{
			name:   "json_autostart_disabled",
			asJSON: true,
			rep: daemon.StatusReport{
				Running: false,
				Socket:  "/tmp/agentd.sock",
				Autostart: daemon.AutostartReport{
					Enabled: false,
					Backend: daemon.BackendSchtasks,
				},
			},
			check: func(t *testing.T, raw []byte) {
				t.Helper()
				var got map[string]any
				require.NoError(t, json.Unmarshal(raw, &got), "json")
				auto := got["autostart"].(map[string]any)
				assert.Equal(t, false, auto["enabled"])
				assert.Equal(t, "schtasks", auto["backend"])
			},
		},
		{
			name:   "json_autostart_enabled",
			asJSON: true,
			rep: daemon.StatusReport{
				Running: true,
				Socket:  "/tmp/agentd.sock",
				Version: "dev",
				Autostart: daemon.AutostartReport{
					Enabled:       true,
					Backend:       daemon.BackendLaunchd,
					ManifestPath:  "/tmp/plist",
					RegisteredExe: "/usr/bin/agentd",
					Stale:         false,
				},
			},
			check: func(t *testing.T, raw []byte) {
				t.Helper()
				var got map[string]any
				require.NoError(t, json.Unmarshal(raw, &got), "json")
				auto := got["autostart"].(map[string]any)
				assert.Equal(t, true, auto["enabled"])
				assert.Equal(t, "launchd", auto["backend"])
				assert.Equal(t, "/tmp/plist", auto["manifest_path"])
			},
		},
		{
			name:   "json_when_daemon_not_running",
			asJSON: true,
			rep: daemon.StatusReport{
				Running: false,
				Socket:  "/tmp/agentd.sock",
				Autostart: daemon.AutostartReport{
					Enabled: true,
					Backend: daemon.BackendSystemd,
				},
			},
			check: func(t *testing.T, raw []byte) {
				t.Helper()
				var got map[string]any
				require.NoError(t, json.Unmarshal(raw, &got), "json")
				assert.Equal(t, false, got["running"])
				auto := got["autostart"].(map[string]any)
				assert.Equal(t, true, auto["enabled"])
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
				Autostart:              daemon.AutostartReport{Enabled: false},
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
