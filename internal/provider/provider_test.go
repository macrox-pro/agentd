package provider_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/provider"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      string
		want    provider.ID
		wantErr bool
	}{
		{name: "claude-code", in: "claude-code", want: provider.ClaudeCode},
		{name: "claude mixed case", in: "Claude-Code", want: provider.ClaudeCode},
		{name: "cursor upper", in: "CURSOR", want: provider.Cursor},
		{name: "codex", in: "codex", want: provider.Codex},
		{name: "gemini", in: "Gemini", want: provider.Gemini},
		{name: "opencode", in: "OpenCode", want: provider.OpenCode},
		{name: "kimi canonical", in: "kimi-code", want: provider.KimiCode},
		{name: "kimi alias", in: "kimicode", want: provider.KimiCode},
		{name: "trim", in: "  cursor  ", want: provider.Cursor},
		{name: "empty", in: "", wantErr: true},
		{name: "unknown", in: "nope", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := provider.Parse(tt.in)
			if tt.wantErr {
				require.Error(t, err, "Parse(%q)", tt.in)
				if tt.in == "" {
					assert.ErrorIs(t, err, provider.ErrProviderRequired)
				}
				if tt.in == "nope" {
					assert.ErrorIs(t, err, provider.ErrUnknownProvider)
				}
				return
			}
			require.NoError(t, err, "Parse(%q)", tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      string
		flagSet bool
		want    provider.ID
		wantErr bool
	}{
		{name: "unset", in: "nope", flagSet: false},
		{name: "set valid", in: "cursor", flagSet: true, want: provider.Cursor},
		{name: "set unknown", in: "nope", flagSet: true, wantErr: true},
		{name: "set empty", in: "", flagSet: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := provider.ParseFilter(tt.in, tt.flagSet)
			if tt.wantErr {
				require.Error(t, err, "ParseFilter(%q, %v)", tt.in, tt.flagSet)
				return
			}
			require.NoError(t, err, "ParseFilter(%q, %v)", tt.in, tt.flagSet)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		in    string
		want  provider.ID
		found bool
	}{
		{name: "known", in: "claude-code", want: provider.ClaudeCode, found: true},
		{name: "alias", in: "Kimicode", want: provider.KimiCode, found: true},
		{name: "unknown", in: "CustomAgent", found: false},
		{name: "empty", in: "", found: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := provider.Lookup(tt.in)
			assert.Equal(t, tt.found, ok, "Lookup(%q)", tt.in)
			if tt.found {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProtoRoundTrip(t *testing.T) {
	t.Parallel()
	ids := []provider.ID{
		provider.ClaudeCode,
		provider.Cursor,
		provider.Codex,
		provider.Gemini,
		provider.OpenCode,
		provider.KimiCode,
	}
	for _, id := range ids {
		t.Run(string(id), func(t *testing.T) {
			t.Parallel()
			p, err := id.Proto()
			require.NoError(t, err, "Proto(%q)", id)
			back, err := provider.FromProto(p)
			require.NoError(t, err, "FromProto(%v)", p)
			assert.Equal(t, id, back)
		})
	}
}

func TestFromProtoUnknown(t *testing.T) {
	t.Parallel()
	_, err := provider.FromProto(agentdv1.Provider_PROVIDER_UNSPECIFIED)
	require.Error(t, err)
}

func TestAgenthooksMapping(t *testing.T) {
	t.Parallel()
	ids := []provider.ID{
		provider.ClaudeCode,
		provider.Cursor,
		provider.Codex,
		provider.Gemini,
		provider.OpenCode,
		provider.KimiCode,
	}
	for _, id := range ids {
		t.Run(string(id), func(t *testing.T) {
			t.Parallel()
			ah, err := id.Agenthooks()
			require.NoError(t, err, "Agenthooks(%q)", id)
			assert.NotEmpty(t, string(ah), "Agenthooks(%q)", id)
		})
	}
}
