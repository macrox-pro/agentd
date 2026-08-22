package trajectory

import (
	"encoding/json"
	"time"
)

const (
	TypeSessionOpen     = "session/open"
	TypeHookInvoked     = "hook/invoked"
	TypeHookDecided     = "hook/decided"
	TypeAsyncDispatched = "async/dispatched"
	TypeAsyncDropped    = "async/dropped"
)

const (
	SourceSystem   = "system"
	SourceHook     = "hook"
	SourceDecision = "decision"
)

// Event is one append-only ledger record (DESIGN §14.3).
type Event struct {
	Seq            uint64          `json:"seq"`
	Type           string          `json:"type"`
	Source         string          `json:"source"`
	TS             time.Time       `json:"ts"`
	Provider       string          `json:"provider"`
	InvocationMode string          `json:"invocation_mode,omitempty"`
	SessionID      string          `json:"session_id"`
	ProjectRoot    string          `json:"project_root,omitempty"`
	CWD            string          `json:"cwd,omitempty"`
	Data           json.RawMessage `json:"data,omitempty"`
	Raw            json.RawMessage `json:"raw,omitempty"`
	Ignorable      bool            `json:"ignorable,omitempty"`
}

// SessionOpenData is the payload for session/open.
type SessionOpenData struct {
	Provider    string `json:"provider"`
	CWD         string `json:"cwd,omitempty"`
	ProjectRoot string `json:"project_root,omitempty"`
}

// HookInvokedData is the payload for hook/invoked.
type HookInvokedData struct {
	Kind      string `json:"kind"`
	ToolName  string `json:"tool_name,omitempty"`
	ToolUseID string `json:"tool_use_id,omitempty"`
	HasRoute  bool   `json:"has_route"`
}

// HookDecidedData is the payload for hook/decided.
type HookDecidedData struct {
	Kind                 string `json:"kind"`
	Decision             string `json:"decision"`
	Reason               string `json:"reason,omitempty"`
	ConfigGeneration     uint64 `json:"config_generation"`
	ConfigFingerprint    string `json:"config_fingerprint"`
	AsyncDispatchedCount uint32 `json:"async_dispatched_count,omitempty"`
}

// AsyncDispatchedData is the payload for async/dispatched.
type AsyncDispatchedData struct {
	Count uint32 `json:"count"`
}

// AsyncDroppedData is the payload for async/dropped (trajectory queue overflow).
type AsyncDroppedData struct {
	Reason string `json:"reason"`
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return b
}
