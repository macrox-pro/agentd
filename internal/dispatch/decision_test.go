package dispatch_test

import (
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestDecisionToProto(t *testing.T) {
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
			got := dispatch.DecisionToProto(tt.in)
			require.NotNil(t, got, "DecisionToProto(%q)", tt.name)
			assert.Equal(t, tt.want, got.Kind, "DecisionToProto(%q)", tt.name)
			if tt.want == agentdv1.DecisionKind_DECISION_KIND_DENY {
				assert.Equal(t, "blocked", got.Reason, "DecisionToProto(%q) reason", tt.name)
			}
		})
	}
}
