package extract

import (
	"encoding/json"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

// TokenFields holds optional token counters from hook raw JSON.
type TokenFields struct {
	Input, Output, CacheRead, CacheWrite, Context uint64
	HasInput, HasOutput, HasCacheRead, HasCacheWrite, HasContext bool
}

// Any reports whether any token field was extracted.
func (t TokenFields) Any() bool {
	return t.HasInput || t.HasOutput || t.HasCacheRead || t.HasCacheWrite || t.HasContext
}

// Tokens extracts provider-specific token fields from raw hook JSON.
func Tokens(prov agentdv1.Provider, raw []byte) TokenFields {
	if len(raw) == 0 {
		return TokenFields{}
	}
	switch prov {
	case agentdv1.Provider_PROVIDER_CODEX:
		return usageFromCodex(raw)
	case agentdv1.Provider_PROVIDER_CURSOR:
		return contextFromCursor(raw)
	case agentdv1.Provider_PROVIDER_CLAUDE_CODE:
		return usageFromClaude(raw)
	default:
		return TokenFields{}
	}
}

func usageFromCodex(raw []byte) TokenFields {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return TokenFields{}
	}
	for _, key := range []string{"usage", "token_usage", "metrics"} {
		if u, ok := top[key]; ok {
			if t := parseUsageObject(u); t.Any() {
				return t
			}
		}
	}
	return TokenFields{}
}

func contextFromCursor(raw []byte) TokenFields {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return TokenFields{}
	}
	var out TokenFields
	if v, ok := uint64Field(m, "context_tokens"); ok {
		out.Context = v
		out.HasContext = true
	}
	return out
}

func usageFromClaude(raw []byte) TokenFields {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return TokenFields{}
	}
	for _, key := range []string{"usage", "message", "result"} {
		if u, ok := top[key]; ok {
			if t := parseUsageObject(u); t.Any() {
				return t
			}
		}
	}
	return TokenFields{}
}

func parseUsageObject(raw json.RawMessage) TokenFields {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return TokenFields{}
	}
	var out TokenFields
	if v, ok := uint64Field(m, "input_tokens"); ok {
		out.Input = v
		out.HasInput = true
	}
	if v, ok := uint64Field(m, "output_tokens"); ok {
		out.Output = v
		out.HasOutput = true
	}
	if v, ok := uint64Field(m, "cache_read_input_tokens"); ok {
		out.CacheRead = v
		out.HasCacheRead = true
	}
	if v, ok := uint64Field(m, "cache_creation_input_tokens"); ok {
		out.CacheWrite = v
		out.HasCacheWrite = true
	}
	return out
}

func uint64Field(m map[string]any, key string) (uint64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		if n < 0 {
			return 0, false
		}
		return uint64(n), true
	case json.Number:
		i, err := n.Int64()
		if err != nil || i < 0 {
			return 0, false
		}
		return uint64(i), true
	default:
		return 0, false
	}
}
