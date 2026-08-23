package guard_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/guard"
)

func TestSecretsApprovalSkipsAsk(t *testing.T) {
	t.Parallel()
	input := `{"command":"export AWS_ACCESS_KEY_ID=` + fakeAWSKey + `"}`
	findings := guard.Scan(json.RawMessage(input), config.DefaultSecretsRules)
	require.NotEmpty(t, findings)
	ids := make([]string, len(findings))
	for i, f := range findings {
		ids[i] = f.ID
	}
	fp := config.ApprovalFingerprint(config.ApprovalKindSecrets, "Bash", config.SecretsStableKey(ids))

	tests := []struct {
		name     string
		dctx     guard.DecisionContext
		wantKind agenthooks.DecisionKind
		wantFP   bool
	}{
		{
			name:     "no approval asks",
			dctx:     guard.DecisionContext{ProjectRoot: "/repo"},
			wantKind: agenthooks.DecisionAsk,
			wantFP:   true,
		},
		{
			name: "matching project approval skips",
			dctx: guard.DecisionContext{
				ProjectRoot: "/repo",
				Approvals: config.Approvals{
					Secrets: []config.Approval{{
						Kind:        config.ApprovalKindSecrets,
						Fingerprint: fp,
						Scope:       config.ApprovalScopeProject,
						Project:     "/repo",
						ExpiresAt:   time.Now().UTC().Add(time.Hour),
						GrantedBy:   "ask_user",
					}},
				},
			},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name: "wrong project still asks",
			dctx: guard.DecisionContext{
				ProjectRoot: "/other",
				Approvals: config.Approvals{
					Secrets: []config.Approval{{
						Kind:        config.ApprovalKindSecrets,
						Fingerprint: fp,
						Scope:       config.ApprovalScopeProject,
						Project:     "/repo",
						ExpiresAt:   time.Now().UTC().Add(time.Hour),
					}},
				},
			},
			wantKind: agenthooks.DecisionAsk,
			wantFP:   true,
		},
		{
			name: "expired asks again",
			dctx: guard.DecisionContext{
				ProjectRoot: "/repo",
				Approvals: config.Approvals{
					Secrets: []config.Approval{{
						Kind:        config.ApprovalKindSecrets,
						Fingerprint: fp,
						Scope:       config.ApprovalScopeProject,
						Project:     "/repo",
						ExpiresAt:   time.Now().UTC().Add(-time.Hour),
					}},
				},
			},
			wantKind: agenthooks.DecisionAsk,
			wantFP:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New()
			guard.AttachSecrets(r, config.SecretsGuard{
				Enabled: true,
				Action:  config.GuardAsk,
				Rules:   config.DefaultSecretsRules,
			}, tt.dctx)
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{
					Provider: agenthooks.ProviderClaudeCode,
					Kind:     agenthooks.KindToolPre,
					Session:  agenthooks.SessionInfo{ID: "s1"},
				},
				Tool: agenthooks.ToolCall{Name: "Bash", Input: json.RawMessage(input)},
			}
			d, err := r.Decide(context.Background(), ev)
			require.NoError(t, err)
			assert.Equal(t, tt.wantKind, d.Kind())
			if tt.wantFP {
				assert.Contains(t, d.SystemMessage(), "approval_fingerprint="+fp)
			}
		})
	}
}

func TestShellApprovalSkipsAsk(t *testing.T) {
	t.Parallel()
	fp := config.ApprovalFingerprint(config.ApprovalKindShell, "Bash", "curl")
	shellInput := json.RawMessage(`{"command":"curl https://example.com"}`)

	tests := []struct {
		name     string
		dctx     guard.DecisionContext
		session  string
		wantKind agenthooks.DecisionKind
	}{
		{
			name: "matching session skips",
			dctx: guard.DecisionContext{
				Approvals: config.Approvals{
					Shell: []config.Approval{{
						Kind:        config.ApprovalKindShell,
						Fingerprint: fp,
						Scope:       config.ApprovalScopeSession,
						SessionID:   "sess",
					}},
				},
			},
			session:  "sess",
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:     "no approval asks",
			dctx:     guard.DecisionContext{},
			session:  "sess",
			wantKind: agenthooks.DecisionAsk,
		},
		{
			name: "wrong session still asks",
			dctx: guard.DecisionContext{
				Approvals: config.Approvals{
					Shell: []config.Approval{{
						Kind:        config.ApprovalKindShell,
						Fingerprint: fp,
						Scope:       config.ApprovalScopeSession,
						SessionID:   "sess",
					}},
				},
			},
			session:  "other",
			wantKind: agenthooks.DecisionAsk,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New()
			guard.AttachShell(r, config.ShellGuard{
				Enabled: true,
				AskOn:   []string{"curl"},
			}, tt.dctx)
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{
					Provider: agenthooks.ProviderClaudeCode,
					Kind:     agenthooks.KindToolPre,
					Session:  agenthooks.SessionInfo{ID: tt.session},
				},
				Tool: agenthooks.ToolCall{
					Name:      "Bash",
					Canonical: agenthooks.ToolShell,
					Input:     shellInput,
				},
			}
			d, err := r.Decide(context.Background(), ev)
			require.NoError(t, err, "Decide(%q)", tt.name)
			assert.Equal(t, tt.wantKind, d.Kind(), "Decide(%q)", tt.name)
		})
	}
}

func TestTemporaryBlockDenies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		until    time.Time
		wantKind agenthooks.DecisionKind
	}{
		{
			name:     "active",
			until:    time.Now().UTC().Add(time.Hour),
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:     "expired ignored",
			until:    time.Now().UTC().Add(-time.Hour),
			wantKind: agenthooks.DecisionNoDecision,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New()
			guard.AttachBlocks(r, []config.TemporaryBlock{{
				Tool:    "Bash",
				Pattern: "dangerous",
				Reason:  "auto-block",
				Until:   tt.until,
			}})
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool: agenthooks.ToolCall{
					Name:      "Bash",
					Canonical: agenthooks.ToolShell,
					Input:     json.RawMessage(`{"command":"do dangerous things"}`),
				},
			}
			d, err := r.Decide(context.Background(), ev)
			require.NoError(t, err)
			assert.Equal(t, tt.wantKind, d.Kind())
		})
	}
}
