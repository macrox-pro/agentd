package guard_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/guard"
)

const fakeAWSKey = "AKIAIOSFODNN7EXAMPLE"

func TestScan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		sample  string
		wantRule string
		rules   []string
		wantHit bool
	}{
		{name: "aws", sample: "key=" + fakeAWSKey, wantRule: "AWS access key ID", wantHit: true},
		{name: "github", sample: "ghp_" + strings.Repeat("a1", 18), wantRule: "GitHub token", wantHit: true},
		{name: "clean", sample: "ls -la", wantHit: false},
		{name: "aws filtered out", sample: "key=" + fakeAWSKey, rules: []string{"jwt"}, wantHit: false},
		{name: "aws filtered in", sample: "key=" + fakeAWSKey, rules: []string{"aws_key"}, wantRule: "AWS access key ID", wantHit: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input, err := json.Marshal(map[string]string{"v": tt.sample})
			require.NoError(t, err, "Marshal")
			findings := guard.Scan(input, tt.rules)
			if !tt.wantHit {
				assert.Empty(t, findings, "Scan(%q)", tt.sample)
				return
			}
			require.NotEmpty(t, findings, "Scan(%q)", tt.sample)
			found := false
			for _, f := range findings {
				if f.Rule == tt.wantRule {
					found = true
					assert.NotContains(t, f.Masked, fakeAWSKey, "mask must not leak full secret")
				}
			}
			assert.True(t, found, "Scan(%q) want rule %q got %v", tt.sample, tt.wantRule, findings)
		})
	}
}

func TestAttachSecrets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		action     config.GuardAction
		input      string
		wantKind   agenthooks.DecisionKind
		wantSecret bool
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
			})

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
	guard.AttachSecrets(r, config.SecretsGuard{Enabled: false, Action: config.GuardAsk})
	ev := &agenthooks.ToolPreEvent{
		Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
		Tool:  agenthooks.ToolCall{Name: "Bash", Input: json.RawMessage(`{"command":"AKIAIOSFODNN7EXAMPLE"}`)},
	}
	d, err := r.Decide(context.Background(), ev)
	require.NoError(t, err)
	assert.Equal(t, agenthooks.DecisionNoDecision, d.Kind())
}
