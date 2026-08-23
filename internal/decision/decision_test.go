package decision_test

import (
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/decision"
)

func TestFromProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		in         *agentdv1.Decision
		wantKind   agenthooks.DecisionKind
		wantReason string
		wantMsg    string
		wantCtx    string
	}{
		{name: "nil", in: nil, wantKind: agenthooks.DecisionNoDecision},
		{
			name:     "no decision",
			in:       &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name: "deny+extras",
			in: &agentdv1.Decision{
				Kind:          agentdv1.DecisionKind_DECISION_KIND_DENY,
				Reason:        "blocked",
				SystemMessage: "sys",
				Context:       "ctx",
			},
			wantKind:   agenthooks.DecisionDeny,
			wantReason: "blocked",
			wantMsg:    "sys",
			wantCtx:    "ctx",
		},
		{
			name: "ask+extras",
			in: &agentdv1.Decision{
				Kind:          agentdv1.DecisionKind_DECISION_KIND_ASK,
				Reason:        "confirm",
				SystemMessage: "ask-sys",
				Context:       "ask-ctx",
			},
			wantKind:   agenthooks.DecisionAsk,
			wantReason: "confirm",
			wantMsg:    "ask-sys",
			wantCtx:    "ask-ctx",
		},
		{
			name: "allow+extras",
			in: &agentdv1.Decision{
				Kind:          agentdv1.DecisionKind_DECISION_KIND_ALLOW,
				SystemMessage: "ok-sys",
				Context:       "ok-ctx",
			},
			wantKind: agenthooks.DecisionAllow,
			wantMsg:  "ok-sys",
			wantCtx:  "ok-ctx",
		},
		{
			name: "block_prompt",
			in: &agentdv1.Decision{
				Kind:   agentdv1.DecisionKind_DECISION_KIND_BLOCK_PROMPT,
				Reason: "blocked prompt",
			},
			wantKind:   agenthooks.DecisionBlockPrompt,
			wantReason: "blocked prompt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := decision.FromProto(tt.in)
			require.NotNil(t, got, "FromProto(%q)", tt.name)
			assert.Equal(t, tt.wantKind, got.Kind(), "FromProto(%q) kind", tt.name)
			if tt.wantReason != "" {
				assert.Equal(t, tt.wantReason, got.Reason(), "FromProto(%q) reason", tt.name)
			}
			if tt.wantMsg != "" {
				assert.Equal(t, tt.wantMsg, got.SystemMessage(), "FromProto(%q) system", tt.name)
			}
			if tt.wantCtx != "" {
				assert.Equal(t, []string{tt.wantCtx}, got.Context(), "FromProto(%q) context", tt.name)
			}
		})
	}
}

func TestToProtoFromProtoRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   agenthooks.Decision
	}{
		{name: "deny", in: agenthooks.Deny("blocked").WithSystemMessage("sys").WithContext("ctx")},
		{name: "ask", in: agenthooks.AskUser("confirm").WithSystemMessage("sys").WithContext("ctx")},
		{name: "allow", in: agenthooks.Allow().WithSystemMessage("sys").WithContext("ctx")},
		{name: "block_prompt", in: agenthooks.BlockPrompt("nope")},
		{name: "no decision", in: agenthooks.NoDecision()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			proto := decision.ToProto(tt.in)
			require.NotNil(t, proto, "ToProto(%q)", tt.name)
			got := decision.FromProto(proto)
			require.NotNil(t, got, "FromProto(%q)", tt.name)
			assert.Equal(t, tt.in.Kind(), got.Kind(), "round-trip(%q) kind", tt.name)
			assert.Equal(t, tt.in.Reason(), got.Reason(), "round-trip(%q) reason", tt.name)
			assert.Equal(t, tt.in.SystemMessage(), got.SystemMessage(), "round-trip(%q) system", tt.name)
			assert.Equal(t, tt.in.Context(), got.Context(), "round-trip(%q) context", tt.name)
		})
	}
}

func TestToProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   agenthooks.Decision
		want agentdv1.DecisionKind
	}{
		{name: "nil", in: nil, want: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION},
		{name: "no decision", in: agenthooks.NoDecision(), want: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION},
		{name: "deny", in: agenthooks.Deny("blocked"), want: agentdv1.DecisionKind_DECISION_KIND_DENY},
		{name: "ask", in: agenthooks.AskUser("confirm"), want: agentdv1.DecisionKind_DECISION_KIND_ASK},
		{name: "allow", in: agenthooks.Allow(), want: agentdv1.DecisionKind_DECISION_KIND_ALLOW},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := decision.ToProto(tt.in)
			require.NotNil(t, got, "ToProto(%q)", tt.name)
			assert.Equal(t, tt.want, got.Kind, "ToProto(%q)", tt.name)
			if tt.want == agentdv1.DecisionKind_DECISION_KIND_DENY {
				assert.Equal(t, "blocked", got.Reason, "ToProto(%q) reason", tt.name)
			}
		})
	}
}
