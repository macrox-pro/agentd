package guard

import (
	"context"
	"fmt"
	"path"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

func mcpHandler(cfg config.MCPGuard) func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
	return func(ctx context.Context, e *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		if e.Tool.MCP == nil {
			return agenthooks.NoDecision(), nil
		}
		server := e.Tool.MCP.Server
		if server == "" {
			return agenthooks.NoDecision(), nil
		}
		hit := matchServer(server, cfg.DenyServers)
		if hit == "" {
			return agenthooks.NoDecision(), nil
		}
		agenthooks.Logger(ctx).Warn("mcp guard: deny_servers matched",
			"tool", e.Tool.Name, "server", server, "pattern", hit)
		return agenthooks.Deny(fmt.Sprintf(
			"mcp tool blocked: server %q matched deny pattern %q", server, hit,
		)), nil
	}
}

func matchServer(server string, patterns []string) string {
	for _, p := range patterns {
		if p == "" {
			continue
		}
		ok, err := path.Match(p, server)
		if err != nil {
			continue
		}
		if ok {
			return p
		}
	}
	return ""
}
