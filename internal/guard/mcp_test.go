package guard_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/guard"
)

func TestAttachMCP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tool     agenthooks.ToolCall
		wantKind agenthooks.DecisionKind
	}{
		{
			name: "deny matched server",
			tool: agenthooks.ToolCall{
				Name:      "mcp__untrusted-foo__bar",
				Canonical: agenthooks.ToolMCP,
				MCP:       &agenthooks.MCPCall{Server: "untrusted-foo", Tool: "bar"},
				Input:     json.RawMessage(`{}`),
			},
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name: "allow other server",
			tool: agenthooks.ToolCall{
				Name:      "mcp__github__create_issue",
				Canonical: agenthooks.ToolMCP,
				MCP:       &agenthooks.MCPCall{Server: "github", Tool: "create_issue"},
				Input:     json.RawMessage(`{}`),
			},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name: "non-mcp ignored",
			tool: agenthooks.ToolCall{
				Name:      "Bash",
				Canonical: agenthooks.ToolShell,
				Input:     json.RawMessage(`{"command":"echo"}`),
			},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name: "empty server ignored",
			tool: agenthooks.ToolCall{
				Name:      "MCP:do",
				Canonical: agenthooks.ToolMCP,
				MCP:       &agenthooks.MCPCall{Tool: "do"},
				Input:     json.RawMessage(`{}`),
			},
			wantKind: agenthooks.DecisionNoDecision,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New()
			guard.AttachMCP(r, config.MCPGuard{
				Enabled:     true,
				DenyServers: []string{"untrusted-*"},
			})
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool:  tt.tool,
			}
			d, err := r.Decide(context.Background(), ev)
			require.NoError(t, err, "Decide(%q)", tt.name)
			require.NotNil(t, d, "Decide(%q)", tt.name)
			assert.Equal(t, tt.wantKind, d.Kind(), "Decide(%q)", tt.name)
		})
	}
}
