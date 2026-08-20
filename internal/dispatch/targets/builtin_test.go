package targets_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch/targets"
)

const fakeAWSKey = "AKIAIOSFODNN7EXAMPLE"

func TestBuiltinDecide(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		action   config.GuardAction
		wantKind agenthooks.DecisionKind
	}{
		{
			name:     "clean",
			input:    `{"command":"go test"}`,
			action:   config.GuardAsk,
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:     "secret ask",
			input:    `{"command":"export K=` + fakeAWSKey + `"}`,
			action:   config.GuardAsk,
			wantKind: agenthooks.DecisionAsk,
		},
		{
			name:     "secret deny",
			input:    `{"command":"export K=` + fakeAWSKey + `"}`,
			action:   config.GuardDeny,
			wantKind: agenthooks.DecisionDeny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := &targets.Builtin{
				Guards: config.Guards{Secrets: config.SecretsGuard{
					Enabled: true,
					Action:  tt.action,
					Rules:   config.DefaultSecretsRules,
				}},
				Policy: config.Policy{
					Fail:        config.FailClosed,
					AskFallback: config.AskFallbackDeny,
				},
			}
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool:  agenthooks.ToolCall{Name: "Bash", Input: json.RawMessage(tt.input)},
			}
			d, err := b.Decide(context.Background(), ev, []string{"secrets"})
			require.NoError(t, err, "Decide(%q)", tt.name)
			require.NotNil(t, d, "Decide(%q)", tt.name)
			assert.Equal(t, tt.wantKind, d.Kind(), "Decide(%q)", tt.name)
			if tt.wantKind != agenthooks.DecisionNoDecision {
				assert.NotContains(t, d.Reason(), fakeAWSKey, "Decide(%q) reason", tt.name)
			}
		})
	}
}

func TestBuiltinObserve(t *testing.T) {
	t.Parallel()
	b := &targets.Builtin{
		Guards: config.Guards{},
		Policy: config.Policy{Fail: config.FailOpen},
	}
	ev := &agenthooks.ToolPreEvent{
		Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
		Tool:  agenthooks.ToolCall{Name: "Bash", Input: json.RawMessage(`{}`)},
	}
	require.NotPanics(t, func() {
		b.Observe(context.Background(), ev)
	})
}

func TestBuiltinDecideGuardSubset(t *testing.T) {
	t.Parallel()

	b := &targets.Builtin{
		Guards: config.Guards{
			Secrets: config.SecretsGuard{
				Enabled: true,
				Action:  config.GuardDeny,
				Rules:   config.DefaultSecretsRules,
			},
			Shell: config.ShellGuard{
				Enabled:      true,
				DenyPatterns: []string{"rm -rf /"},
				AskOn:        []string{"curl"},
			},
			MCP: config.MCPGuard{
				Enabled:     true,
				DenyServers: []string{"untrusted-*"},
			},
			Paths: config.PathsGuard{
				Enabled:  true,
				DenyRead: []string{"/etc/shadow"},
			},
		},
		Policy: config.Policy{
			Fail:        config.FailClosed,
			AskFallback: config.AskFallbackDeny,
		},
	}

	tests := []struct {
		name       string
		guards     []string
		tool       agenthooks.ToolCall
		wantKind   agenthooks.DecisionKind
		wantReason string
	}{
		{
			name:   "shell subset ignores secrets",
			guards: []string{"shell"},
			tool: agenthooks.ToolCall{
				Name:      "Bash",
				Canonical: agenthooks.ToolShell,
				Input:     json.RawMessage(`{"command":"export K=` + fakeAWSKey + `"}`),
			},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:   "shell deny via subset",
			guards: []string{"shell"},
			tool: agenthooks.ToolCall{
				Name:      "Bash",
				Canonical: agenthooks.ToolShell,
				Input:     json.RawMessage(`{"command":"rm -rf /"}`),
			},
			wantKind:   agenthooks.DecisionDeny,
			wantReason: "deny pattern",
		},
		{
			name:   "mcp deny",
			guards: []string{"mcp"},
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
			name:   "paths deny",
			guards: []string{"paths"},
			tool: agenthooks.ToolCall{
				Name:      "Read",
				Canonical: agenthooks.ToolFileRead,
				Input:     json.RawMessage(`{"file_path":"/etc/shadow"}`),
			},
			wantKind:   agenthooks.DecisionDeny,
			wantReason: "path",
		},
		{
			name:   "multi guard shell ask",
			guards: []string{"secrets", "shell", "mcp", "paths"},
			tool: agenthooks.ToolCall{
				Name:      "Bash",
				Canonical: agenthooks.ToolShell,
				Input:     json.RawMessage(`{"command":"curl https://example.com"}`),
			},
			wantKind: agenthooks.DecisionAsk,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool:  tt.tool,
			}
			d, err := b.Decide(context.Background(), ev, tt.guards)
			require.NoError(t, err, "Decide(%q)", tt.name)
			require.NotNil(t, d, "Decide(%q)", tt.name)
			assert.Equal(t, tt.wantKind, d.Kind(), "Decide(%q)", tt.name)
			if tt.wantReason != "" {
				assert.Contains(t, d.Reason(), tt.wantReason, "Decide(%q) reason", tt.name)
			}
		})
	}
}

