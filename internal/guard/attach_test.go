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

func TestAttachSecrets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		action     config.GuardAction
		fallback   config.AskFallback
		input      string
		wantKind   agenthooks.DecisionKind
		wantSecret bool
		noCapAsk   bool
	}{
		{
			name:     "clean no decision",
			action:   config.GuardAsk,
			input:    `{"command":"go test ./..."}`,
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:       "ask on hit with CapAsk",
			action:     config.GuardAsk,
			input:      `{"command":"export AWS_ACCESS_KEY_ID=` + fakeAWSKey + `"}`,
			wantKind:   agenthooks.DecisionAsk,
			wantSecret: true,
		},
		{
			name:       "deny action",
			action:     config.GuardDeny,
			input:      `{"command":"export AWS_ACCESS_KEY_ID=` + fakeAWSKey + `"}`,
			wantKind:   agenthooks.DecisionDeny,
			wantSecret: true,
		},
		{
			name:       "unsupported ask fallback deny",
			action:     config.GuardAsk,
			fallback:   config.AskFallbackDeny,
			input:      `{"command":"export AWS_ACCESS_KEY_ID=` + fakeAWSKey + `"}`,
			wantKind:   agenthooks.DecisionDeny,
			wantSecret: true,
			noCapAsk:   true,
		},
		{
			name:       "unsupported ask fallback no_decision",
			action:     config.GuardAsk,
			fallback:   config.AskFallbackNoDecision,
			input:      `{"command":"export AWS_ACCESS_KEY_ID=` + fakeAWSKey + `"}`,
			wantKind:   agenthooks.DecisionNoDecision,
			wantSecret: false,
			noCapAsk:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New(agenthooks.WithPolicy(agenthooks.Policy{
				Fail:        agenthooks.FailClosed,
				AskFallback: agenthooks.FallbackDeny,
			}))
			guard.AttachSecrets(r, config.SecretsGuard{
				Enabled: true,
				Action:  tt.action,
				Rules:   config.DefaultSecretsRules,
			}, guard.DecisionContext{AskFallback: tt.fallback})

			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{
					Provider: agenthooks.ProviderClaudeCode,
					Kind:     agenthooks.KindToolPre,
				},
				Tool: agenthooks.ToolCall{
					Name:  "Bash",
					Input: json.RawMessage(tt.input),
				},
			}
			if tt.noCapAsk {
				ev.Provider = agenthooks.ProviderCodex
			}
			d, err := r.Decide(context.Background(), ev)
			require.NoError(t, err, "Decide")
			require.NotNil(t, d, "Decide")
			assert.Equal(t, tt.wantKind, d.Kind(), "Decide kind")
			if tt.wantSecret {
				assert.NotContains(t, d.Reason(), fakeAWSKey, "reason must not leak secret")
			}
		})
	}
}

func TestAttachSecretsDisabled(t *testing.T) {
	t.Parallel()
	r := agenthooks.New()
	guard.AttachSecrets(r, config.SecretsGuard{Enabled: false, Action: config.GuardAsk}, guard.DecisionContext{})
	ev := &agenthooks.ToolPreEvent{
		Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
		Tool:  agenthooks.ToolCall{Name: "Bash", Input: json.RawMessage(`{"command":"AKIAIOSFODNN7EXAMPLE"}`)},
	}
	d, err := r.Decide(context.Background(), ev)
	require.NoError(t, err)
	assert.Equal(t, agenthooks.DecisionNoDecision, d.Kind())
}
