package dispatch_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestDecodeTypedTable(t *testing.T) {
	t.Parallel()

	claudeRaw, err := json.Marshal(map[string]any{
		"session_id":      "s",
		"cwd":             "/w",
		"hook_event_name": "PreToolUse",
		"tool_name":       "Bash",
		"tool_use_id":     "t1",
		"tool_input":      map[string]any{"command": "echo"},
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		provider    agentdv1.Provider
		mode        agentdv1.InvocationMode
		raw         []byte
		checkNotify bool
	}{
		{
			name:     "cursor argv",
			provider: agentdv1.Provider_PROVIDER_CURSOR,
			mode:     agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
			raw:      agenthookstest.Fixture(t, "cursor/pre_tool_use.json"),
		},
		{
			name:        "codex notify",
			provider:    agentdv1.Provider_PROVIDER_CODEX,
			mode:        agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY,
			raw:         agenthookstest.Fixture(t, "codex/pre_tool_use.json"),
			checkNotify: true,
		},
		{
			name:     "cursor unspecified",
			provider: agentdv1.Provider_PROVIDER_CURSOR,
			mode:     agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
			raw:      agenthookstest.Fixture(t, "cursor/pre_tool_use.json"),
		},
		{
			name:     "claude stdin",
			provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			mode:     agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			raw:      claudeRaw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			typed, err := dispatch.DecodeTyped(context.Background(), tt.provider, tt.mode, tt.raw)
			require.NoError(t, err, "DecodeTyped(%q)", tt.name)
			if tt.checkNotify {
				base := agenthooks.EventOf(typed)
				require.NotNil(t, base, "DecodeTyped(%q)", tt.name)
				assert.Equal(t, agenthooks.KindNotification, base.Kind, "DecodeTyped(%q)", tt.name)
				return
			}
			_, ok := typed.(*agenthooks.ToolPreEvent)
			assert.True(t, ok, "DecodeTyped(%q)", tt.name)
		})
	}
}
