package hookedge

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/decision"
)

// encodeDecision encodes provider wire for a daemon Decision.
// agenthooks does not export Encode; Runner.Run is the supported wire path.
// Handlers return the daemon decision so the edge cannot invent policy —
// Decide remains in the daemon.
func encodeDecision(ctx context.Context, provider string, argvPayload bool, payload []byte, d *agentdv1.Decision, stdout, stderr io.Writer) int {
	r := agenthooks.New(agenthooks.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	dec := decision.FromProto(d)
	toolPre := agenthooks.NoDecision()
	if tp, ok := dec.(agenthooks.ToolPreDecision); ok {
		toolPre = tp
	}
	r.OnToolPre(func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		return toolPre, nil
	})
	r.OnPermission(func(context.Context, *agenthooks.PermissionEvent) (agenthooks.ToolPreDecision, error) {
		return toolPre, nil
	})
	r.OnPromptSubmitted(func(context.Context, *agenthooks.PromptEvent) (agenthooks.PromptDecision, error) {
		return promptDecisionFromProto(d), nil
	})
	r.OnStop(func(context.Context, *agenthooks.StopEvent) (agenthooks.StopDecision, error) {
		return agenthooks.Finish(), nil
	})
	r.OnToolPost(func(context.Context, *agenthooks.ToolPostEvent) (agenthooks.ToolPostDecision, error) {
		return agenthooks.Observed(), nil
	})
	r.OnToolError(func(context.Context, *agenthooks.ToolPostEvent) (agenthooks.ToolPostDecision, error) {
		return agenthooks.Observed(), nil
	})

	args := []string{"run", "--provider=" + provider}
	var stdin io.Reader
	if argvPayload {
		args = append(args, "--argv-payload", string(payload))
		stdin = bytes.NewReader(nil)
	} else {
		stdin = bytes.NewReader(payload)
	}
	return r.Run(ctx, args, stdin, stdout, stderr)
}

func promptDecisionFromProto(d *agentdv1.Decision) agenthooks.PromptDecision {
	if d == nil {
		return agenthooks.AcceptPrompt()
	}
	switch d.GetKind() {
	case agentdv1.DecisionKind_DECISION_KIND_DENY, agentdv1.DecisionKind_DECISION_KIND_BLOCK_PROMPT:
		return agenthooks.BlockPrompt(d.GetReason())
	default:
		return agenthooks.AcceptPrompt()
	}
}
