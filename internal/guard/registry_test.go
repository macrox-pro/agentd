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

func TestAttachCheckers(t *testing.T) {
	t.Parallel()

	policy := agenthooks.Policy{
		Fail:        agenthooks.FailClosed,
		AskFallback: agenthooks.FallbackDeny,
	}

	tests := []struct {
		name       string
		names      []string
		guards     config.Guards
		tool       agenthooks.ToolCall
		wantKind   agenthooks.DecisionKind
		wantReason string
	}{
		{
			name:  "secrets_deny",
			names: []string{"secrets"},
			guards: config.Guards{Secrets: config.SecretsGuard{
				Enabled: true,
				Action:  config.GuardDeny,
				Rules:   config.DefaultSecretsRules,
			}},
			tool: agenthooks.ToolCall{
				Name:  "Bash",
				Input: json.RawMessage(`{"command":"export K=` + fakeAWSKey + `"}`),
			},
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:  "shell_deny",
			names: []string{"shell"},
			guards: config.Guards{Shell: config.ShellGuard{
				Enabled:      true,
				DenyPatterns: []string{"rm -rf /"},
			}},
			tool: agenthooks.ToolCall{
				Name:      "Bash",
				Canonical: agenthooks.ToolShell,
				Input:     json.RawMessage(`{"command":"rm -rf /"}`),
			},
			wantKind:   agenthooks.DecisionDeny,
			wantReason: "deny pattern",
		},
		{
			name:  "mcp_deny",
			names: []string{"mcp"},
			guards: config.Guards{MCP: config.MCPGuard{
				Enabled:     true,
				DenyServers: []string{"untrusted-*"},
			}},
			tool: agenthooks.ToolCall{
				Name:      "mcp__untrusted-x__t",
				Canonical: agenthooks.ToolMCP,
				MCP:       &agenthooks.MCPCall{Server: "untrusted-x", Tool: "t"},
				Input:     json.RawMessage(`{}`),
			},
			wantKind:   agenthooks.DecisionDeny,
			wantReason: "mcp",
		},
		{
			name:  "paths_deny",
			names: []string{"paths"},
			guards: config.Guards{Paths: config.PathsGuard{
				Enabled:  true,
				DenyRead: []string{"/etc/shadow"},
			}},
			tool: agenthooks.ToolCall{
				Name:      "Read",
				Canonical: agenthooks.ToolFileRead,
				Input:     json.RawMessage(`{"file_path":"/etc/shadow"}`),
			},
			wantKind:   agenthooks.DecisionDeny,
			wantReason: "path",
		},
		{
			name:  "secrets_disabled",
			names: []string{"secrets"},
			guards: config.Guards{Secrets: config.SecretsGuard{
				Enabled: false,
				Action:  config.GuardDeny,
				Rules:   config.DefaultSecretsRules,
			}},
			tool: agenthooks.ToolCall{
				Name:  "Bash",
				Input: json.RawMessage(`{"command":"export K=` + fakeAWSKey + `"}`),
			},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:  "unknown_name_skipped",
			names: []string{"nope"},
			guards: config.Guards{Secrets: config.SecretsGuard{
				Enabled: true,
				Action:  config.GuardDeny,
				Rules:   config.DefaultSecretsRules,
			}},
			tool: agenthooks.ToolCall{
				Name:  "Bash",
				Input: json.RawMessage(`{"command":"export K=` + fakeAWSKey + `"}`),
			},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:  "multi_guard_order",
			names: []string{"shell", "secrets"},
			guards: config.Guards{
				Secrets: config.SecretsGuard{
					Enabled: true,
					Action:  config.GuardDeny,
					Rules:   config.DefaultSecretsRules,
				},
				Shell: config.ShellGuard{
					Enabled:      true,
					DenyPatterns: []string{"rm -rf /"},
				},
			},
			tool: agenthooks.ToolCall{
				Name:      "Bash",
				Canonical: agenthooks.ToolShell,
				Input:     json.RawMessage(`{"command":"rm -rf /"}`),
			},
			wantKind:   agenthooks.DecisionDeny,
			wantReason: "deny pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New(agenthooks.WithPolicy(policy))
			guard.AttachCheckers(r, tt.guards, guard.DecisionContext{}, tt.names)

			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{
					Provider: agenthooks.ProviderClaudeCode,
					Kind:     agenthooks.KindToolPre,
				},
				Tool: tt.tool,
			}
			d, err := r.Decide(context.Background(), ev)
			require.NoError(t, err, "Decide(%q)", tt.name)
			require.NotNil(t, d, "Decide(%q)", tt.name)
			assert.Equal(t, tt.wantKind, d.Kind(), "Decide(%q)", tt.name)
			if tt.wantReason != "" {
				assert.Contains(t, d.Reason(), tt.wantReason, "Decide(%q) reason", tt.name)
			}
		})
	}
}
