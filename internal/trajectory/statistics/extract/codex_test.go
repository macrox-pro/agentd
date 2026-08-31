package extract_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics/extract"
)

func writeCodexRollout(t *testing.T, lines ...string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "rollout.jsonl")
	var b []byte
	for _, line := range lines {
		b = append(b, line...)
		b = append(b, '\n')
	}
	require.NoError(t, os.WriteFile(path, b, 0o600))
	return path
}

func codexTokenCountLine(input, cached, cacheWrite, output uint64) string {
	return `{"type":"event_msg","payload":{"type":"token_count","info":{"last_token_usage":{"input_tokens":` +
		strconv.FormatUint(input, 10) + `,"cached_input_tokens":` + strconv.FormatUint(cached, 10) +
		`,"cache_write_input_tokens":` + strconv.FormatUint(cacheWrite, 10) +
		`,"output_tokens":` + strconv.FormatUint(output, 10) + `}}}}`
}

func codexStopHook(transcriptPath string) string {
	return `{"hook_event_name":"Stop","session_id":"s1","transcript_path":"` + transcriptPath + `"}`
}

func codexPromptHook(transcriptPath string) string {
	return `{"hook_event_name":"UserPromptSubmit","session_id":"s1","transcript_path":"` + transcriptPath + `"}`
}

func TestScanCodexTranscript(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		hookRaw string
		rollout []string
		check   func(t *testing.T, got extract.TokenFields)
	}{
		{
			name:    "codex_stop_full_fields",
			hookRaw: "",
			rollout: []string{codexTokenCountLine(15156, 4352, 0, 100)},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.Equal(t, uint64(15156), got.Input)
				assert.Equal(t, uint64(4352), got.CacheRead)
				assert.Equal(t, uint64(0), got.CacheWrite)
				assert.Equal(t, uint64(100), got.Output)
			},
		},
		{
			name:    "codex_prompt_submit_skipped",
			hookRaw: "",
			rollout: []string{codexTokenCountLine(999, 0, 0, 1)},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name:    "codex_malformed_hook_raw",
			hookRaw: `{not json`,
			rollout: []string{codexTokenCountLine(10, 0, 0, 1)},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name:    "codex_empty_transcript_path",
			hookRaw: `{"hook_event_name":"Stop","transcript_path":""}`,
			rollout: []string{codexTokenCountLine(10, 0, 0, 1)},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name:    "codex_missing_file",
			hookRaw: `{"hook_event_name":"Stop","transcript_path":"/no/such/rollout.jsonl"}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name:    "codex_empty_transcript_file",
			hookRaw: "",
			rollout: []string{},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name:    "codex_malformed_tail_line",
			hookRaw: "",
			rollout: []string{
				`{"type":"event_msg","payload":{"type":"token_count","broken"}`,
				codexTokenCountLine(5, 0, 0, 1),
			},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasInput)
				assert.Equal(t, uint64(5), got.Input)
			},
		},
		{
			name:    "codex_no_token_count_event",
			hookRaw: "",
			rollout: []string{`{"type":"event_msg","payload":{"type":"task_complete"}}`},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name:    "codex_total_only_ignored",
			hookRaw: "",
			rollout: []string{`{"type":"event_msg","payload":{"type":"token_count","info":{"total_token_usage":{"input_tokens":999}}}}`},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name:    "codex_partial_fields",
			hookRaw: "",
			rollout: []string{`{"type":"event_msg","payload":{"type":"token_count","info":{"last_token_usage":{"input_tokens":42}}}}`},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasInput)
				assert.Equal(t, uint64(42), got.Input)
				assert.False(t, got.HasOutput)
			},
		},
		{
			name:    "codex_multiple_token_count_picks_last",
			hookRaw: "",
			rollout: []string{
				codexTokenCountLine(10, 0, 0, 1),
				codexTokenCountLine(20, 0, 0, 2),
			},
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.Equal(t, uint64(20), got.Input)
				assert.Equal(t, uint64(2), got.Output)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hookRaw := tt.hookRaw
			if hookRaw == "" {
				var path string
				if tt.rollout != nil {
					path = writeCodexRollout(t, tt.rollout...)
				} else {
					path = filepath.Join(t.TempDir(), "empty.jsonl")
					require.NoError(t, os.WriteFile(path, nil, 0o600))
				}
				switch tt.name {
				case "codex_prompt_submit_skipped":
					hookRaw = codexPromptHook(path)
				default:
					hookRaw = codexStopHook(path)
				}
			}
			got := extract.TokensFromTranscript(agentdv1.Provider_PROVIDER_CODEX, []byte(hookRaw))
			tt.check(t, got)
		})
	}
}
