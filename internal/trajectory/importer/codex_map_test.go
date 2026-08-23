package importer_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestMapCodexRolloutLineTable(t *testing.T) {
	t.Parallel()
	// ImportCodex exercises mapCodexRolloutLine; table covers skip / edge shapes via temp files.
	tests := []struct {
		name      string
		lines     []string
		wantTypes []string
		wantTools []string
	}{
		{
			name: "agent_reasoning empty",
			lines: []string{
				`{"type":"event_msg","payload":{"type":"agent_reasoning","text":"  "}}`,
			},
			wantTypes: nil,
		},
		{
			name: "agent_reasoning present",
			lines: []string{
				`{"type":"event_msg","payload":{"type":"agent_reasoning","text":"**Plan**"}}`,
			},
			wantTypes: []string{string(trajectory.TypeTranscriptThinking)},
		},
		{
			name: "skip meta telemetry",
			lines: []string{
				`{"type":"session_meta","payload":{"session_id":"s1"}}`,
				`{"type":"turn_context","payload":{"turn_id":"t1"}}`,
				`{"type":"world_state","payload":{"full":{}}}`,
				`{"type":"event_msg","payload":{"type":"token_count","info":{}}}`,
				`{"type":"event_msg","payload":{"type":"task_started","turn_id":"t1"}}`,
				`{"type":"event_msg","payload":{"type":"task_complete","turn_id":"t1"}}`,
			},
			wantTypes: nil,
		},
		{
			name: "skip developer response_item",
			lines: []string{
				`{"type":"response_item","payload":{"type":"message","role":"developer","content":[{"type":"input_text","text":"x"}]}}`,
			},
			wantTypes: nil,
		},
		{
			name: "malformed line",
			lines: []string{
				`not-json`,
				`{"type":"event_msg","payload":{"type":"user_message","message":"ok"}}`,
			},
			wantTypes: []string{string(trajectory.TypeTranscriptMessage)},
		},
		{
			name: "function_call empty call_id",
			lines: []string{
				`{"type":"response_item","payload":{"type":"function_call","name":"exec_command","call_id":"","arguments":"{}"}}`,
			},
			wantTypes: nil,
		},
		{
			name: "custom_tool_call",
			lines: []string{
				`{"type":"response_item","payload":{"type":"custom_tool_call","name":"mcp","call_id":"c1","arguments":"{}"}}`,
				`{"type":"response_item","payload":{"type":"custom_tool_call_output","call_id":"c1","output":"out"}}`,
			},
			wantTypes: []string{string(trajectory.TypeTranscriptMessage), string(trajectory.TypeTranscriptMessage)},
			wantTools: []string{"c1", "c1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "t.jsonl")
			var body strings.Builder
			for _, line := range tt.lines {
				body.WriteString(line)
				body.WriteByte('\n')
			}
			require.NoError(t, os.WriteFile(path, []byte(body.String()), 0o600), "write %s", tt.name)

			result, err := importer.ImportCodex(importer.ImportOptions{
				SessionID:      "map-" + tt.name,
				TranscriptPath: path,
				Cfg:            config.TrajectoryConfig{},
			})
			require.NoError(t, err, "ImportCodex(%q)", tt.name)
			require.Len(t, result.Events, len(tt.wantTypes), "ImportCodex(%q) event count", tt.name)
			for i, e := range result.Events {
				assert.Equal(t, tt.wantTypes[i], string(e.Type), "ImportCodex(%q) type[%d]", tt.name, i)
				if i < len(tt.wantTools) {
					var d trajectory.TranscriptMessageData
					require.NoError(t, json.Unmarshal(e.Data, &d))
					assert.Equal(t, tt.wantTools[i], d.ToolUseID, "ImportCodex(%q) tool[%d]", tt.name, i)
				}
			}
		})
	}
}
