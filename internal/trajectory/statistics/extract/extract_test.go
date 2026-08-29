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
		name string
		prov agentdv1.Provider
		raw  string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := extract.Tokens(tt.prov, []byte(tt.raw))
			tt.check(t, got)
		})
	}
}
