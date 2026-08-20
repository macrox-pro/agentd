package hookedge

import (
	"fmt"
	"strings"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

func parseProvider(s string) (agentdv1.Provider, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "claude-code":
		return agentdv1.Provider_PROVIDER_CLAUDE_CODE, nil
	case "cursor":
		return agentdv1.Provider_PROVIDER_CURSOR, nil
	case "codex":
		return agentdv1.Provider_PROVIDER_CODEX, nil
	case "gemini":
		return agentdv1.Provider_PROVIDER_GEMINI, nil
	case "opencode":
		return agentdv1.Provider_PROVIDER_OPENCODE, nil
	case "kimicode", "kimi-code":
		return agentdv1.Provider_PROVIDER_KIMI_CODE, nil
	case "":
		return 0, fmt.Errorf("provider is required")
	default:
		return 0, fmt.Errorf("unknown provider %q", s)
	}
}
