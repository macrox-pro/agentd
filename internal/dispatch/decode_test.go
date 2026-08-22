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

func TestDecodeTypedCursorArgv(t *testing.T) {
	t.Parallel()
	raw := agenthookstest.Fixture(t, "cursor/pre_tool_use.json")
	typed, err := dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CURSOR, agentdv1.InvocationMode_INVOCATION_MODE_ARGV, raw)
	require.NoError(t, err, "DecodeTyped(cursor, argv)")
	_, ok := typed.(*agenthooks.ToolPreEvent)
	assert.True(t, ok, "typed event kind")
}

func TestDecodeTypedCodexNotify(t *testing.T) {
	t.Parallel()
	raw := agenthookstest.Fixture(t, "codex/pre_tool_use.json")
	typed, err := dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CODEX, agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY, raw)
	require.NoError(t, err, "DecodeTyped(codex, notify)")
	base := agenthooks.EventOf(typed)
	require.NotNil(t, base, "EventOf")
	assert.Equal(t, agenthooks.KindNotification, base.Kind)
}

func TestDecodeTypedCursorUnspecified(t *testing.T) {
	t.Parallel()
	raw := agenthookstest.Fixture(t, "cursor/pre_tool_use.json")
	typed, err := dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CURSOR, agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED, raw)
	require.NoError(t, err, "DecodeTyped(cursor, unspecified)")
	_, ok := typed.(*agenthooks.ToolPreEvent)
	assert.True(t, ok, "typed event kind")
}

func TestDecodeTypedClaudeStdin(t *testing.T) {
	t.Parallel()
	raw, err := json.Marshal(map[string]any{
		"session_id":      "s",
		"cwd":             "/w",
		"hook_event_name": "PreToolUse",
		"tool_name":       "Bash",
		"tool_use_id":     "t1",
		"tool_input":      map[string]any{"command": "echo"},
	})
	require.NoError(t, err)
	typed, err := dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CLAUDE_CODE, agentdv1.InvocationMode_INVOCATION_MODE_STDIN, raw)
	require.NoError(t, err, "DecodeTyped(claude, stdin)")
	_, ok := typed.(*agenthooks.ToolPreEvent)
	assert.True(t, ok, "typed event kind")
}
