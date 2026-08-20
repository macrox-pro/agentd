package guard

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

// AttachSecrets registers secrets scanning handlers on r.
func AttachSecrets(r *agenthooks.Runner, cfg config.SecretsGuard) {
	if r == nil || !cfg.Enabled {
		return
	}
	handler := secretsHandler(cfg)
	r.OnToolPre(handler)
	r.OnPermission(func(ctx context.Context, e *agenthooks.PermissionEvent) (agenthooks.ToolPreDecision, error) {
		return handler(ctx, &agenthooks.ToolPreEvent{
			Event: e.Event,
			Tool:  e.Tool,
		})
	})
}

func secretsHandler(cfg config.SecretsGuard) func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
	return func(ctx context.Context, e *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		findings := Scan(e.Tool.Input, cfg.Rules)
		if len(findings) == 0 {
			return agenthooks.NoDecision(), nil
		}
		found := describe(findings)
		agenthooks.Logger(ctx).Warn("secrets guard: credential-shaped strings in tool input",
			"tool", e.Tool.Name, "findings", len(findings))

		forceDeny := cfg.Action == config.GuardDeny || !e.Can(agenthooks.CapAsk)
		if forceDeny {
			return agenthooks.Deny(fmt.Sprintf(
				"tool call blocked: input contains credential-shaped strings: %s",
				found,
			)), nil
		}
		return agenthooks.AskUser(fmt.Sprintf(
			"tool call input contains credential-shaped strings: %s. Approve to accept the risk and continue; reject to block the call.",
			found,
		)).WithSystemMessage("secrets: " + found), nil
	}
}

func describe(findings []Finding) string {
	parts := make([]string, len(findings))
	for i, f := range findings {
		parts[i] = fmt.Sprintf("%s (%s)", f.Rule, f.Masked)
	}
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}
