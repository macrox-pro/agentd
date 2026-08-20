package server

import (
	"context"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

func (h *hookService) Invoke(ctx context.Context, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	snap := h.store.Current()
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
		Provider:   req.GetProvider(),
		RawPayload: req.GetRawPayload(),
		Snap:       snap,
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
	return resp, nil
}
