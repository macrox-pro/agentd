package dispatch_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/speakeasy-api/agenthooks/agenthookstest"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestClaudeThenCursorDecode(t *testing.T) {
	t.Parallel()
	rawClaude, err := json.Marshal(map[string]any{
		"session_id":      "s",
		"cwd":             "/w",
		"hook_event_name": "PreToolUse",
		"tool_name":       "Bash",
		"tool_use_id":     "t1",
		"tool_input":      map[string]any{"command": "go test"},
	})
	require.NoError(t, err)
	_, err = dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CLAUDE_CODE, agentdv1.InvocationMode_INVOCATION_MODE_STDIN, rawClaude)
	require.NoError(t, err, "DecodeTyped(claude)")
	rawCursor := agenthookstest.Fixture(t, "cursor/pre_tool_use.json")
	_, err = dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CURSOR, agentdv1.InvocationMode_INVOCATION_MODE_ARGV, rawCursor)
	require.NoError(t, err, "DecodeTyped(cursor after claude)")
}
