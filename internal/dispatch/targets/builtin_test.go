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
