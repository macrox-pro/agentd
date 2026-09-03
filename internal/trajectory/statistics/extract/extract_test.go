package extract_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics/extract"
)

func TestExtract(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		prov  agentdv1.Provider
		raw   string
		check func(t *testing.T, got extract.TokenFields)
	}{
		{
			name: "codex_usage",
			prov: agentdv1.Provider_PROVIDER_CODEX,
			raw:  `{"usage":{"input_tokens":10,"output_tokens":5}}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasInput)
				assert.Equal(t, uint64(10), got.Input)
				assert.True(t, got.HasOutput)
				assert.Equal(t, uint64(5), got.Output)
			},
		},
		{
			name: "cursor_stop_top_level",
			prov: agentdv1.Provider_PROVIDER_CURSOR,
			raw:  `{"input_tokens":19582,"output_tokens":92,"cache_read_tokens":6272,"cache_write_tokens":0}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasInput)
				assert.Equal(t, uint64(19582), got.Input)
				assert.True(t, got.HasOutput)
				assert.Equal(t, uint64(92), got.Output)
				assert.True(t, got.HasCacheRead)
				assert.Equal(t, uint64(6272), got.CacheRead)
				assert.True(t, got.HasCacheWrite)
				assert.Equal(t, uint64(0), got.CacheWrite)
			},
		},
		{
			name: "cursor_stop_partial",
			prov: agentdv1.Provider_PROVIDER_CURSOR,
			raw:  `{"input_tokens":42}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasInput)
				assert.Equal(t, uint64(42), got.Input)
				assert.False(t, got.HasOutput)
				assert.False(t, got.HasCacheRead)
				assert.False(t, got.HasCacheWrite)
			},
		},
		{
			name: "cursor_pre_compact",
			prov: agentdv1.Provider_PROVIDER_CURSOR,
			raw:  `{"context_tokens":120000}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasContext)
				assert.Equal(t, uint64(120000), got.Context)
			},
		},
		{
			name: "codex_cached_input_tokens",
			prov: agentdv1.Provider_PROVIDER_CODEX,
			raw:  `{"usage":{"input_tokens":10,"cached_input_tokens":5,"output_tokens":2}}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasInput)
				assert.Equal(t, uint64(10), got.Input)
				assert.True(t, got.HasCacheRead)
				assert.Equal(t, uint64(5), got.CacheRead)
				assert.True(t, got.HasOutput)
				assert.Equal(t, uint64(2), got.Output)
			},
		},
		{
			name: "claude_agent_usage",
			prov: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			raw:  `{"usage":{"input_tokens":3,"output_tokens":7,"cache_read_input_tokens":1,"cache_creation_input_tokens":2}}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasInput)
				assert.True(t, got.HasOutput)
				assert.True(t, got.HasCacheRead)
				assert.True(t, got.HasCacheWrite)
			},
		},
		{
			name: "wrong_hook_no_tokens",
			prov: agentdv1.Provider_PROVIDER_GEMINI,
			raw:  `{"usage":{"input_tokens":99}}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name: "cache_write_input_tokens",
			prov: agentdv1.Provider_PROVIDER_CODEX,
			raw:  `{"usage":{"cache_write_input_tokens":9}}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.True(t, got.HasCacheWrite)
				assert.Equal(t, uint64(9), got.CacheWrite)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := extract.Tokens(tt.prov, []byte(tt.raw))
			tt.check(t, got)
		})
	}
}

func TestTokensFromTranscript(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		prov  agentdv1.Provider
		raw   string
		check func(t *testing.T, got extract.TokenFields)
	}{
		{
			name: "codex_no_scanner_provider",
			prov: agentdv1.Provider_PROVIDER_GEMINI,
			raw:  `{"hook_event_name":"Stop","transcript_path":"/tmp/x.jsonl"}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name: "codex_empty_raw",
			prov: agentdv1.Provider_PROVIDER_CODEX,
			raw:  "",
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
		{
			name: "codex_unregistered_prov_zero",
			prov: agentdv1.Provider_PROVIDER_UNSPECIFIED,
			raw:  `{"hook_event_name":"Stop","transcript_path":"/tmp/x.jsonl"}`,
			check: func(t *testing.T, got extract.TokenFields) {
				t.Helper()
				assert.False(t, got.Any())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := extract.TokensFromTranscript(tt.prov, []byte(tt.raw))
			tt.check(t, got)
		})
	}
}
