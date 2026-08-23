package dispatch_test

import (
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestFirstConclusive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    agenthooks.Decision
		want bool
	}{
		{name: "nil_continues", d: nil, want: false},
		{name: "no_decision_continues", d: agenthooks.NoDecision(), want: false},
		{name: "deny_stops", d: agenthooks.Deny("blocked"), want: true},
		{name: "ask_stops", d: agenthooks.AskUser("confirm"), want: true},
		{name: "allow_stops", d: agenthooks.Allow(), want: true},
		{name: "block_prompt_stops", d: agenthooks.BlockPrompt("stop"), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, dispatch.FirstConclusive(tt.d))
		})
	}
}
