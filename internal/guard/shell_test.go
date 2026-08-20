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

func TestAttachShell(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tool     string
		command  string
		provider agenthooks.Provider
		wantKind agenthooks.DecisionKind
	}{
		{
			name:     "clean",
			tool:     "Bash",
			command:  "go test ./...",
			provider: agenthooks.ProviderClaudeCode,
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:     "deny",
			tool:     "Bash",
			command:  "sudo rm -rf /",
			provider: agenthooks.ProviderClaudeCode,
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:     "ask",
			tool:     "Bash",
			command:  "curl https://example.com",
			provider: agenthooks.ProviderClaudeCode,
			wantKind: agenthooks.DecisionAsk,
		},
		{
			name:     "ask falls back to deny without CapAsk",
			tool:     "Bash",
			command:  "curl https://example.com",
			provider: agenthooks.ProviderCodex,
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:     "non-shell ignored",
			tool:     "Read",
			command:  "rm -rf /",
			provider: agenthooks.ProviderClaudeCode,
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:     "deny wins over ask",
			tool:     "Bash",
			command:  "curl http://x && rm -rf /",
			provider: agenthooks.ProviderClaudeCode,
			wantKind: agenthooks.DecisionDeny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New()
			guard.AttachShell(r, config.ShellGuard{
				Enabled:      true,
				DenyPatterns: []string{"rm -rf /"},
				AskOn:        []string{"curl"},
			})
			input, err := json.Marshal(map[string]string{"command": tt.command})
			require.NoError(t, err)
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{
					Provider: tt.provider,
					Kind:     agenthooks.KindToolPre,
				},
				Tool: agenthooks.ToolCall{
					Name:      tt.tool,
					Canonical: agenthooks.CanonicalToolFor(tt.tool),
					Input:     input,
				},
			}
			d, err := r.Decide(context.Background(), ev)
			require.NoError(t, err, "Decide(%q)", tt.name)
			require.NotNil(t, d, "Decide(%q)", tt.name)
			assert.Equal(t, tt.wantKind, d.Kind(), "Decide(%q)", tt.name)
		})
	}
}

func TestAttachShellDisabled(t *testing.T) {
	t.Parallel()
	r := agenthooks.New()
	guard.AttachShell(r, config.ShellGuard{
		Enabled:      false,
		DenyPatterns: []string{"rm -rf /"},
	})
	ev := &agenthooks.ToolPreEvent{
		Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
		Tool: agenthooks.ToolCall{
			Name:      "Bash",
			Canonical: agenthooks.ToolShell,
			Input:     json.RawMessage(`{"command":"rm -rf /"}`),
		},
	}
	d, err := r.Decide(context.Background(), ev)
	require.NoError(t, err)
	assert.Equal(t, agenthooks.DecisionNoDecision, d.Kind())
}
