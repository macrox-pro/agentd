package guard

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

func shellHandler(cfg config.ShellGuard, dctx DecisionContext) func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
	return func(ctx context.Context, e *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		if e.Tool.Canonical != agenthooks.ToolShell {
			return agenthooks.NoDecision(), nil
		}
		cmd := shellCommand(e.Tool.Input)
		if cmd == "" {
			return agenthooks.NoDecision(), nil
		}
		if hit := firstSubstring(cmd, cfg.DenyPatterns); hit != "" {
			agenthooks.Logger(ctx).Warn("shell guard: deny pattern matched",
				"tool", e.Tool.Name, "pattern", hit)
			return agenthooks.Deny(fmt.Sprintf(
				"shell command blocked: matched deny pattern %q", hit,
			)), nil
		}
		if hit := firstSubstring(cmd, cfg.AskOn); hit != "" {
			fp := config.ApprovalFingerprint(config.ApprovalKindShell, e.Tool.Name, hit)
			if approved(dctx, config.ApprovalKindShell, fp, e.Session.ID) {
				agenthooks.Logger(ctx).Info("shell guard: approval matched, skipping ask",
					"tool", e.Tool.Name, "fingerprint", fp)
				return agenthooks.NoDecision(), nil
			}
			agenthooks.Logger(ctx).Warn("shell guard: ask_on matched",
				"tool", e.Tool.Name, "pattern", hit)
			if !e.Can(agenthooks.CapAsk) {
				return agenthooks.Deny(fmt.Sprintf(
					"shell command blocked: matched ask_on pattern %q (ask unsupported)", hit,
				)), nil
			}
			return agenthooks.AskUser(fmt.Sprintf(
				"shell command matched ask_on pattern %q. Approve to continue; reject to block.",
				hit,
			)).WithSystemMessage(fmt.Sprintf("shell: %s; approval_fingerprint=%s", hit, fp)), nil
		}
		return agenthooks.NoDecision(), nil
	}
}

func shellCommand(input json.RawMessage) string {
	var obj map[string]any
	if err := json.Unmarshal(input, &obj); err != nil {
		return ""
	}
	for _, key := range []string{"command", "cmd"} {
		v, ok := obj[key]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if ok {
			return s
		}
	}
	return ""
}

func firstSubstring(haystack string, patterns []string) string {
	for _, p := range patterns {
		if p == "" {
			continue
		}
		if strings.Contains(haystack, p) {
			return p
		}
	}
	return ""
}
