package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

func TestNormalizeInvocationMode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		provider agentdv1.Provider
		in       agentdv1.InvocationMode
		want     agentdv1.InvocationMode
	}{
		{
			name:     "cursor unspecified to argv",
			provider: agentdv1.Provider_PROVIDER_CURSOR,
			in:       agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
			want:     agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
		},
		{
			name:     "cursor argv preserved",
			provider: agentdv1.Provider_PROVIDER_CURSOR,
			in:       agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
			want:     agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
		},
		{
			name:     "claude unspecified stays unspecified",
			provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			in:       agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
			want:     agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, normalizeInvocationMode(tt.provider, tt.in))
		})
	}
}
