package server

import (
	"context"
	"log/slog"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/decision"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

// SnapshotSource supplies project-aware compiled config for one Invoke.
type SnapshotSource interface {
	SnapshotFor(cwd, projectRoot string) *config.Snapshot
}

// Invoker is the dispatch sync/async pipeline port used by HookService.
type Invoker = dispatch.Invoker

var (
	_ SnapshotSource = (*config.Store)(nil)
	_ Invoker        = (*dispatch.Engine)(nil)
)

// NewHookService builds the HookService gRPC handler with injectable ports.
func NewHookService(snap SnapshotSource, inv Invoker, rec *trajectory.Recorder, collector *statistics.Collector, log *slog.Logger) agentdv1.HookServiceServer {
	return &hookService{
		snap:      snap,
		engine:    inv,
		recorder:  rec,
		collector: collector,
		log:       log,
	}
}

func (h *hookService) Invoke(ctx context.Context, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	resp := &agentdv1.InvokeResponse{
		Decision: decision.Neutral(),
		Config:   &agentdv1.ConfigGeneration{},
	}
	if h.snap == nil {
		return resp, nil
	}
	snap := h.snap.SnapshotFor(req.GetCwd(), req.GetProjectRoot())
	resp.Config = &agentdv1.ConfigGeneration{
		Generation:  snap.Generation,
		Fingerprint: snap.Fingerprint,
	}
	if h.engine == nil {
		return resp, nil
	}
	in := dispatch.InvokeInput{
		Provider:       req.GetProvider(),
		RawPayload:     req.GetRawPayload(),
		Snap:           snap,
		InvocationMode: normalizeInvocationMode(req.GetProvider(), req.GetInvocationMode()),
		CWD:            req.GetCwd(),
		ProjectRoot:    req.GetProjectRoot(),
	}
	if dl := req.GetDeadline(); dl != nil {
		in.Deadline = dl.AsTime()
	}
	result, err := h.engine.Invoke(ctx, in)
	if err != nil {
		// Match agenthooks wire: undecodable payloads become a neutral no-op.
		if log := h.log; log != nil {
			log.Warn("invoke failed; skipping trajectory record",
				"provider", req.GetProvider(),
				"error", err,
			)
		}
		return resp, nil
	}
	resp.Decision = result.Decision
	resp.AsyncDispatchedCount = result.AsyncDispatchedCount
	recIn := trajectory.RecordInput{
		Provider:       req.GetProvider(),
		InvocationMode: normalizeInvocationMode(req.GetProvider(), req.GetInvocationMode()),
		CWD:            req.GetCwd(),
		ProjectRoot:    req.GetProjectRoot(),
		RawPayload:     req.GetRawPayload(),
		Result:         result,
		Snap:           snap,
	}
	if h.collector != nil {
		h.collector.Observe(recIn)
	}
	if h.recorder != nil {
		h.recorder.Record(recIn)
	}
	return resp, nil
}
