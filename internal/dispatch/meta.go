package dispatch

import (
	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/dispatch/targets"
)

// InvokeMeta carries trajectory fields extracted once during Invoke.
type InvokeMeta struct {
	Provider  string
	SessionID string
	EventKind string
	ToolName  string
	ToolUseID string
	HasRoute  bool
}

// MetaFromTyped builds trajectory metadata from a decoded event.
func MetaFromTyped(providerName string, typed any, hasRoute bool) InvokeMeta {
	meta := InvokeMeta{
		Provider:  providerName,
		EventKind: targets.EventKindOf(typed),
		HasRoute:  hasRoute,
	}
	if base := agenthooks.EventOf(typed); base != nil {
		meta.SessionID = base.Session.ID
	}
	switch e := typed.(type) {
	case *agenthooks.ToolPreEvent:
		meta.ToolName = e.Tool.Name
		meta.ToolUseID = e.Tool.ID
	case *agenthooks.ToolPostEvent:
		meta.ToolName = e.Tool.Name
		meta.ToolUseID = e.Tool.ID
	}
	return meta
}
