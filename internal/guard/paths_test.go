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

func TestAttachPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tool     string
		input    map[string]string
		wantKind agenthooks.DecisionKind
	}{
		{
			name:     "deny read shadow",
			tool:     "Read",
			input:    map[string]string{"file_path": "/etc/shadow"},
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:     "allow other read",
			tool:     "Read",
			input:    map[string]string{"file_path": "/tmp/ok.txt"},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:     "deny write env doublestar",
			tool:     "Write",
			input:    map[string]string{"file_path": "repo/.env"},
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:     "deny write nested env",
			tool:     "Edit",
			input:    map[string]string{"path": "a/b/.env"},
			wantKind: agenthooks.DecisionDeny,
		},
		{
			name:     "shell ignored",
			tool:     "Bash",
			input:    map[string]string{"command": "cat /etc/shadow"},
			wantKind: agenthooks.DecisionNoDecision,
		},
		{
			name:     "write miss",
			tool:     "Write",
			input:    map[string]string{"file_path": "repo/readme.md"},
			wantKind: agenthooks.DecisionNoDecision,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := agenthooks.New()
			guard.AttachPaths(r, config.PathsGuard{
				Enabled:   true,
				DenyRead:  []string{"/etc/shadow"},
				DenyWrite: []string{"**/.env"},
			})
			input, err := json.Marshal(tt.input)
			require.NoError(t, err)
			ev := &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
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
