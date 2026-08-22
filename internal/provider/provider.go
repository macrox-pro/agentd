// Package provider canonicalizes coding-agent provider ids shared across hookedge,
// install, trajectory, and dispatch.
//
// Owns: canonical id strings, alias normalization, strict Parse, proto/agenthooks mapping.
// Must not: wire I/O, ledger persistence, install manifest, importer logic, cmd/Cobra.
//
// Use Parse for required provider id strings (empty/unknown → sentinel errors).
// Use ParseFilter for optional filters with an explicit flag-set bit.
// Use Lookup for lenient normalization of known ids only.
//
// Hot path: other.
package provider

import (
	"fmt"
	"strings"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

// ID is a canonical provider id string (e.g. "claude-code").
type ID string

const (
	ClaudeCode ID = "claude-code"
	Cursor     ID = "cursor"
	Codex      ID = "codex"
	Gemini     ID = "gemini"
	OpenCode   ID = "opencode"
	KimiCode   ID = "kimi-code"
)

func (id ID) String() string {
	return string(id)
}

// Parse validates and returns a canonical provider id (case-insensitive).
func Parse(s string) (ID, error) {
	id, ok := Lookup(s)
	if !ok {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return "", ErrProviderRequired
		}
		return "", fmt.Errorf("%w %q", ErrUnknownProvider, s)
	}
	return id, nil
}

// ParseFilter parses s when flagSet is true (optional provider filter).
// When !flagSet, returns ("", nil) so omitted flags do not reuse a prior parse value.
func ParseFilter(s string, flagSet bool) (ID, error) {
	if !flagSet {
		return "", nil
	}
	return Parse(s)
}

// Lookup returns the canonical id when s is a known provider id or alias.
func Lookup(s string) (ID, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "claude-code":
		return ClaudeCode, true
	case "cursor":
		return Cursor, true
	case "codex":
		return Codex, true
	case "gemini":
		return Gemini, true
	case "opencode":
		return OpenCode, true
	case "kimicode", "kimi-code":
		return KimiCode, true
	default:
		return "", false
	}
}

// Proto maps a canonical id to the gRPC Provider enum.
func (id ID) Proto() (agentdv1.Provider, error) {
	switch id {
	case ClaudeCode:
		return agentdv1.Provider_PROVIDER_CLAUDE_CODE, nil
	case Cursor:
		return agentdv1.Provider_PROVIDER_CURSOR, nil
	case Codex:
		return agentdv1.Provider_PROVIDER_CODEX, nil
	case Gemini:
		return agentdv1.Provider_PROVIDER_GEMINI, nil
	case OpenCode:
		return agentdv1.Provider_PROVIDER_OPENCODE, nil
	case KimiCode:
		return agentdv1.Provider_PROVIDER_KIMI_CODE, nil
	default:
		return 0, fmt.Errorf("unknown provider %q", id)
	}
}

// Agenthooks maps a canonical id to the agenthooks install provider enum.
func (id ID) Agenthooks() (agenthooks.Provider, error) {
	switch id {
	case ClaudeCode:
		return agenthooks.ProviderClaudeCode, nil
	case Cursor:
		return agenthooks.ProviderCursor, nil
	case Codex:
		return agenthooks.ProviderCodex, nil
	case Gemini:
		return agenthooks.ProviderGemini, nil
	case OpenCode:
		return agenthooks.ProviderOpenCode, nil
	case KimiCode:
		return agenthooks.ProviderKimi, nil
	default:
		return "", fmt.Errorf("unknown provider %q", id)
	}
}

// FromProto maps a gRPC Provider enum to a canonical id.
func FromProto(p agentdv1.Provider) (ID, error) {
	switch p {
	case agentdv1.Provider_PROVIDER_CLAUDE_CODE:
		return ClaudeCode, nil
	case agentdv1.Provider_PROVIDER_CURSOR:
		return Cursor, nil
	case agentdv1.Provider_PROVIDER_CODEX:
		return Codex, nil
	case agentdv1.Provider_PROVIDER_GEMINI:
		return Gemini, nil
	case agentdv1.Provider_PROVIDER_OPENCODE:
		return OpenCode, nil
	case agentdv1.Provider_PROVIDER_KIMI_CODE:
		return KimiCode, nil
	default:
		return "", fmt.Errorf("unknown provider %v", p)
	}
}
