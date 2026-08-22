package trajectory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestCanonicalProvider(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "claude lower", in: "claude-code", want: "claude-code"},
		{name: "claude mixed", in: "Claude-Code", want: "claude-code"},
		{name: "cursor upper", in: "CURSOR", want: "cursor"},
		{name: "codex mixed", in: "Codex", want: "codex"},
		{name: "gemini", in: "Gemini", want: "gemini"},
		{name: "opencode", in: "OpenCode", want: "opencode"},
		{name: "kimi alias", in: "Kimicode", want: "kimi-code"},
		{name: "kimi canonical", in: "kimi-code", want: "kimi-code"},
		{name: "unknown trimmed", in: "  CustomAgent  ", want: "CustomAgent"},
		{name: "empty", in: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := trajectory.CanonicalProvider(tt.in)
			assert.Equal(t, tt.want, got, "CanonicalProvider(%q)", tt.in)
		})
	}
}

func TestResolveSessionKeyUsesCanonicalProvider(t *testing.T) {
	t.Parallel()
	key := trajectory.ResolveSessionKey("Claude-Code", "s1", "/proj", "")
	require.Equal(t, "claude-code", key.Provider)
	require.Equal(t, "s1", key.SessionID)
}
