package server

import (
	"context"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func (h *hookService) Invoke(ctx context.Context, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	snap := h.store.SnapshotFor(req.GetCwd(), req.GetProjectRoot())
	resp := &agentdv1.InvokeResponse{
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
		Decision: dispatch.NeutralDecision(),
	}
	if h.engine == nil {
		return resp, nil
	}
	in := dispatch.InvokeInput{
		Provider:       req.GetProvider(),
		RawPayload:     req.GetRawPayload(),
		Snap:           snap,
		InvocationMode: req.GetInvocationMode(),
		CWD:            req.GetCwd(),
		ProjectRoot:    req.GetProjectRoot(),
	}
	if dl := req.GetDeadline(); dl != nil {
		in.Deadline = dl.AsTime()
	}
	result, err := h.engine.Invoke(ctx, in)
	if err != nil {
		// Match agenthooks wire: undecodable payloads become a neutral no-op.
		return resp, nil
	}
	resp.Decision = result.Decision
	resp.AsyncDispatchedCount = result.AsyncDispatchedCount
	if h.recorder != nil {
		h.recorder.Record(trajectory.RecordInput{
			Provider:       req.GetProvider(),
			InvocationMode: req.GetInvocationMode(),
			CWD:            req.GetCwd(),
			ProjectRoot:    req.GetProjectRoot(),
			RawPayload:     req.GetRawPayload(),
			Result:         result,
			Snap:           snap,
		})
	}
	return resp, nil
}
