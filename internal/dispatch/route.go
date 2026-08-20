package dispatch

import (
	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

// MatchRoute returns the first compiled route for the event kind, or the "other" default.
func MatchRoute(routes []config.CompiledRoute, typed any) *config.CompiledRoute {
	kind := string(agenthooks.KindOther)
	if base := agenthooks.EventOf(typed); base != nil && base.Kind != "" {
		kind = string(base.Kind)
	}
	var other *config.CompiledRoute
	for i := range routes {
		r := &routes[i]
		if r.Kind == kind {
			return r
		}
		if r.Kind == string(agenthooks.KindOther) {
			other = r
		}
	}
	return other
}
