package trajectory_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/provider"
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
	tests := []struct {
		name       string
		provider   string
		sessionID  string
		project    string
		cwd        string
		wantProv   provider.ID
		wantSID    string
		weakStable bool
	}{
		{
			name:      "alias normalization",
			provider:  "Claude-Code",
			sessionID: "s1",
			project:   "/proj",
			wantProv:  provider.ClaudeCode,
			wantSID:   "s1",
		},
		{
			name:      "kimi alias",
			provider:  "kimicode",
			sessionID: "s2",
			project:   "/proj",
			wantProv:  provider.KimiCode,
			wantSID:   "s2",
		},
		{
			name:       "weak id stable",
			provider:   "kimi-code",
			sessionID:  "",
			project:    "/proj",
			cwd:        "/cwd",
			wantProv:   provider.KimiCode,
			weakStable: true,
		},
		{
			name:      "weak id hash unchanged",
			provider:  "kimi-code",
			sessionID: "",
			project:   "/proj",
			cwd:       "/cwd",
			wantProv:  provider.KimiCode,
			wantSID:   expectedWeakSessionID("kimi-code", "/proj", "/cwd"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			key := trajectory.ResolveSessionKey(tt.provider, tt.sessionID, tt.project, tt.cwd)
			require.Equal(t, tt.wantProv, key.Provider)
			if tt.wantSID != "" {
				require.Equal(t, tt.wantSID, key.SessionID)
			}
			if tt.weakStable {
				again := trajectory.ResolveSessionKey(tt.provider, tt.sessionID, tt.project, tt.cwd)
				assert.Equal(t, key.SessionID, again.SessionID)
				assert.NotEmpty(t, key.SessionID)
			}
		})
	}
}

func expectedWeakSessionID(providerName, projectRoot, cwd string) string {
	sum := sha256.Sum256([]byte(providerName + "\x00" + projectRoot + "\x00" + cwd))
	return "weak-" + hex.EncodeToString(sum[:8])
}

func TestResolveSessionKeyID(t *testing.T) {
	t.Parallel()
	key := trajectory.ResolveSessionKeyID(provider.Cursor, "s1", "/proj", "")
	assert.Equal(t, provider.Cursor, key.Provider)
	assert.Equal(t, "s1", key.SessionID)
}
