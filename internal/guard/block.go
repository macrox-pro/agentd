package guard

import (
	"context"
	"fmt"
	"time"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

// AttachBlocks registers temporary runtime block handlers (Deny) on r.
func AttachBlocks(r *agenthooks.Runner, blocks []config.TemporaryBlock) {
	if r == nil || len(blocks) == 0 {
		return
	}
	attachToolPre(r, blocksHandler(blocks))
}

func blocksHandler(blocks []config.TemporaryBlock) func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
	return func(ctx context.Context, e *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		haystack := blockHaystack(e)
		hit := config.MatchTemporaryBlock(blocks, e.Tool.Name, haystack, time.Now().UTC())
		if hit == nil {
			return agenthooks.NoDecision(), nil
		}
		reason := hit.Reason
		if reason == "" {
			reason = fmt.Sprintf("matched temporary block pattern %q", hit.Pattern)
		}
		agenthooks.Logger(ctx).Warn("temporary block: deny",
			"tool", e.Tool.Name, "pattern", hit.Pattern, "reason", reason)
		return agenthooks.Deny(fmt.Sprintf("tool call blocked: %s", reason)), nil
	}
}

func blockHaystack(e *agenthooks.ToolPreEvent) string {
	if e.Tool.Canonical == agenthooks.ToolShell {
		if cmd := shellCommand(e.Tool.Input); cmd != "" {
			return cmd
		}
	}
	return string(e.Tool.Input)
}
