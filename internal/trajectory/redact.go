package trajectory

import (
	"encoding/json"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/guard"
)

// PrepareRaw returns raw payload bytes for storage per trajectory config.
func PrepareRaw(raw []byte, cfg config.TrajectoryConfig) json.RawMessage {
	if !cfg.Enabled || !cfg.IncludeRaw || len(raw) == 0 {
		return nil
	}
	out := append([]byte(nil), raw...)
	if cfg.RedactSecretRules {
		out = redactRaw(out)
	}
	if cfg.MaxEventBytes > 0 && len(out) > cfg.MaxEventBytes {
		out = out[:cfg.MaxEventBytes]
	}
	return json.RawMessage(out)
}

func redactRaw(raw []byte) []byte {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		b, _ := json.Marshal(string(raw))
		if len(guard.Scan(b, nil)) > 0 {
			return []byte(`"[REDACTED]"`)
		}
		return raw
	}
	redactValue(&v)
	out, err := json.Marshal(v)
	if err != nil {
		return raw
	}
	return out
}

func redactValue(v *any) {
	switch t := (*v).(type) {
	case map[string]any:
		for k, val := range t {
			redactValue(&val)
			t[k] = val
		}
	case []any:
		for i, val := range t {
			redactValue(&val)
			t[i] = val
		}
	case string:
		b, err := json.Marshal(t)
		if err != nil {
			return
		}
		if len(guard.Scan(b, nil)) > 0 {
			*v = "[REDACTED]"
		}
	}
}
