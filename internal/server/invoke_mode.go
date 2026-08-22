package server

import agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"

// normalizeInvocationMode maps wire/provider pairs when invocation_mode is missing (UNSPECIFIED).
func normalizeInvocationMode(provider agentdv1.Provider, mode agentdv1.InvocationMode) agentdv1.InvocationMode {
	if mode != agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED {
		return mode
	}
	switch provider {
	case agentdv1.Provider_PROVIDER_CURSOR:
		return agentdv1.InvocationMode_INVOCATION_MODE_ARGV
	case agentdv1.Provider_PROVIDER_CODEX:
		// Codex notify hooks pass JSON in argv; stdin run remains the default for explicit run.
		return agentdv1.InvocationMode_INVOCATION_MODE_STDIN
	default:
		return mode
	}
}
