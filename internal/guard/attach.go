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
func AttachSecrets(r *agenthooks.Runner, cfg config.SecretsGuard, dctx DecisionContext) {
	if r == nil || !cfg.Enabled {
		return
	}
	attachToolPre(r, secretsHandler(cfg, dctx))
}

// AttachShell registers shell deny/ask handlers on r.
func AttachShell(r *agenthooks.Runner, cfg config.ShellGuard, dctx DecisionContext) {
	if r == nil || !cfg.Enabled {
		return
	}
	attachToolPre(r, shellHandler(cfg, dctx))
}

// AttachMCP registers MCP server deny handlers on r.
func AttachMCP(r *agenthooks.Runner, cfg config.MCPGuard) {
	if r == nil || !cfg.Enabled {
		return
	}
	attachToolPre(r, mcpHandler(cfg))
}

// AttachPaths registers filesystem path deny handlers on r.
func AttachPaths(r *agenthooks.Runner, cfg config.PathsGuard) {
	if r == nil || !cfg.Enabled {
		return
	}
	attachToolPre(r, pathsHandler(cfg))
}

func attachToolPre(r *agenthooks.Runner, handler func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error)) {
	r.OnToolPre(handler)
	r.OnPermission(func(ctx context.Context, e *agenthooks.PermissionEvent) (agenthooks.ToolPreDecision, error) {
		return handler(ctx, &agenthooks.ToolPreEvent{
			Event: e.Event,
			Tool:  e.Tool,
		})
	})
}

func secretsHandler(cfg config.SecretsGuard, dctx DecisionContext) func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
	return func(ctx context.Context, e *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		findings := Scan(e.Tool.Input, cfg.Rules)
		if len(findings) == 0 {
			return agenthooks.NoDecision(), nil
		}
		found := describe(findings)
		fp := config.ApprovalFingerprint(
			config.ApprovalKindSecrets,
			e.Tool.Name,
			config.SecretsStableKey(findingRuleIDs(findings)),
		)
		if approved(dctx, config.ApprovalKindSecrets, fp, e.Session.ID) {
			agenthooks.Logger(ctx).Info("secrets guard: approval matched, skipping ask",
				"tool", e.Tool.Name, "fingerprint", fp)
			return agenthooks.NoDecision(), nil
		}

		agenthooks.Logger(ctx).Warn("secrets guard: credential-shaped strings in tool input",
			"tool", e.Tool.Name, "findings", len(findings))

		if cfg.Action == config.GuardDeny {
			return agenthooks.Deny(fmt.Sprintf(
				"tool call blocked: input contains credential-shaped strings: %s",
				found,
			)), nil
		}
		if !e.Can(agenthooks.CapAsk) {
			return askUnsupported(dctx.AskFallback, fmt.Sprintf(
				"tool call blocked: input contains credential-shaped strings: %s",
				found,
			)), nil
		}
		return agenthooks.AskUser(fmt.Sprintf(
			"tool call input contains credential-shaped strings: %s. Approve to accept the risk and continue; reject to block the call.",
			found,
		)).WithSystemMessage(fmt.Sprintf("secrets: %s; approval_fingerprint=%s", found, fp)), nil
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
