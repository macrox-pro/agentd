package dispatch

import (
	"strings"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

// MatchRoute returns the first compiled route that matches the event, or a
// default-kind fallback. User routes (Default=false) are listed first.
func MatchRoute(routes []config.CompiledRoute, typed any) *config.CompiledRoute {
	base := agenthooks.EventOf(typed)
	kind := string(agenthooks.KindOther)
	provider := ""
	if base != nil {
		if base.Kind != "" {
			kind = string(base.Kind)
		}
		provider = string(base.Provider)
	}
	toolName, toolCanon := toolIdentity(typed)

	var fallback *config.CompiledRoute
	for i := range routes {
		r := &routes[i]
		if !routeMatches(r, kind, provider, toolName, toolCanon) {
			continue
		}
		if !r.Default {
			return r
		}
		if fallback == nil {
			fallback = r
		}
	}
	return fallback
}

func routeMatches(r *config.CompiledRoute, kind, provider, toolName string, toolCanon agenthooks.CanonicalTool) bool {
	// Legacy single-kind default routes.
	if r.Kind != "" && len(r.Match.Kinds) == 0 {
		return r.Kind == kind || r.Kind == string(agenthooks.KindOther)
	}
	m := r.Match
	if len(m.Kinds) > 0 && !stringIn(m.Kinds, kind) {
		return false
	}
	if !providerMatches(m.Providers, provider) {
		return false
	}
	if len(m.Tools) > 0 && !toolMatches(m.Tools, toolName, toolCanon) {
		return false
	}
	return true
}

func providerMatches(providers []string, provider string) bool {
	if len(providers) == 0 {
		return true
	}
	for _, p := range providers {
		if p == "*" || strings.EqualFold(p, provider) {
			return true
		}
	}
	return false
}

func toolMatches(tools []string, name string, canon agenthooks.CanonicalTool) bool {
	if name == "" && canon == "" {
		return false
	}
	for _, t := range tools {
		if strings.EqualFold(t, name) {
			return true
		}
		if canon != "" && strings.EqualFold(t, string(canon)) {
			return true
		}
		// DESIGN examples use Shell / MCP as class labels.
		if canon == agenthooks.ToolShell && strings.EqualFold(t, "Shell") {
			return true
		}
		if canon == agenthooks.ToolMCP && strings.EqualFold(t, "MCP") {
			return true
		}
	}
	return false
}

func toolIdentity(typed any) (name string, canon agenthooks.CanonicalTool) {
	switch ev := typed.(type) {
	case *agenthooks.ToolPreEvent:
		return ev.Tool.Name, ev.Tool.Canonical
	case *agenthooks.ToolPostEvent:
		return ev.Tool.Name, ev.Tool.Canonical
	case *agenthooks.PermissionEvent:
		return ev.Tool.Name, ev.Tool.Canonical
	default:
		return "", ""
	}
}

func stringIn(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}
