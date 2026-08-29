package daemon

import (
	"context"
	"time"

	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/transport"
)

// StatusReport is a snapshot of daemon liveness for CLI status.
type StatusReport struct {
	Running                bool
	Socket                 string
	Version                string // process that answered Status, not this CLI binary
	StartedAt              time.Time
	Generation             uint64
	Fingerprint            string
	AsyncQueueDepth        uint32
	AsyncDroppedCount      uint64
	TrajectoryDroppedCount uint64
	CompiledRouteCount     uint32
	Autostart              AutostartReport
}

// Status probes the daemon and returns a StatusReport. When the daemon is
// unreachable, Running is false and other fields are zero (Socket is still set).
func Status(ctx context.Context, socket string) (StatusReport, error) {
	if socket == "" {
		socket = transport.DefaultSocketPath()
	}
	rep := StatusReport{Socket: socket}

	cli, err := hookclient.Dial(ctx, socket)
	if err != nil {
		autostart, _ := AutostartStatus()
		rep.Autostart = autostart
		return rep, nil
	}
	defer cli.Close()

	resp, err := cli.Status(ctx)
	if err != nil {
		autostart, _ := AutostartStatus()
		rep.Autostart = autostart
		return rep, nil
	}

	rep.Running = true
	rep.Version = resp.GetVersion()
	if ts := resp.GetStartedAt(); ts != nil {
		rep.StartedAt = ts.AsTime().UTC()
	}
	if cfg := resp.GetConfig(); cfg != nil {
		rep.Generation = cfg.GetGeneration()
		rep.Fingerprint = cfg.GetFingerprint()
	}
	rep.AsyncQueueDepth = resp.GetAsyncQueueDepth()
	rep.AsyncDroppedCount = resp.GetAsyncDroppedCount()
	rep.TrajectoryDroppedCount = resp.GetTrajectoryDroppedCount()
	rep.CompiledRouteCount = resp.GetCompiledRouteCount()
	rep.Autostart, _ = AutostartStatus()
	return rep, nil
}
