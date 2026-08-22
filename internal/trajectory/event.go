package trajectory

import (
	"encoding/json"
	"time"
)

const (
	TypeSessionOpen        = "session/open"
	TypeHookInvoked        = "hook/invoked"
	TypeHookDecided        = "hook/decided"
	TypeAsyncDispatched    = "async/dispatched"
	TypeAsyncDropped       = "async/dropped"
	TypeTranscriptMessage  = "transcript/message"
	TypeTranscriptThinking = "transcript/thinking"
	TypeSessionFork        = "session/fork"
	TypeSessionEndSeed     = "session/end-seed"
)

const (
	SourceSystem    = "system"
	SourceHook      = "hook"
	SourceDecision  = "decision"
	SourceTranscript = "transcript"
)

// SchemaVersion is the frozen trajectory event contract version (v1.1).
const SchemaVersion uint32 = 1

// Event is one append-only ledger record (DESIGN §14.3).
type Event struct {
	SchemaVersion  uint32          `json:"schema_version,omitempty"`
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

// TranscriptMessageData is the payload for transcript/message.
type TranscriptMessageData struct {
	Role                 string `json:"role,omitempty"`
	Text                 string `json:"text,omitempty"`
	ToolUseID            string `json:"tool_use_id,omitempty"`
	TranscriptLineIndex  int    `json:"transcript_line_index"`
}

// TranscriptThinkingData is the payload for transcript/thinking.
type TranscriptThinkingData struct {
	Text                string `json:"text,omitempty"`
	TranscriptLineIndex int    `json:"transcript_line_index"`
}

// SessionForkData is the payload for session/fork (audit lineage).
type SessionForkData struct {
	ParentProvider string `json:"parent_provider"`
	ParentSession  string `json:"parent_session"`
	BoundarySeq    uint64 `json:"boundary_seq"`
}

// SessionEndSeedData marks the lineage boundary after a fork seed copy.
type SessionEndSeedData struct {
	ParentProvider string `json:"parent_provider"`
	ParentSession  string `json:"parent_session"`
	BoundarySeq    uint64 `json:"boundary_seq"`
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return b
}

func stampSchemaVersion(e *Event) {
	if e != nil && e.SchemaVersion == 0 {
		e.SchemaVersion = SchemaVersion
	}
}
