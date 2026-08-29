package statistics

import (
	"context"
	"fmt"
	"time"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

// Session is on-demand stats computed from a session JSONL ledger.
type Session struct {
	Provider              string            `json:"provider"`
	SessionID             string            `json:"session_id"`
	EventCount            uint64            `json:"event_count"`
	FirstSeq              uint64            `json:"first_seq"`
	LastSeq               uint64            `json:"last_seq"`
	FirstTS               time.Time         `json:"first_ts"`
	LastTS                time.Time         `json:"last_ts"`
	HooksByKind           map[string]uint64 `json:"hooks_by_kind"`
	DecisionsByKind       map[string]uint64 `json:"decisions_by_kind"`
	AsyncDispatched       uint64            `json:"async_dispatched"`
	AsyncDropped          uint64            `json:"async_dropped"`
	EventsByType          map[string]uint64 `json:"events_by_type"`
	EventsBySource        map[string]uint64 `json:"events_by_source"`
	TranscriptMessages    uint64            `json:"transcript_messages"`
	TranscriptThinking    uint64            `json:"transcript_thinking"`
	InputTokensTotal      uint64            `json:"input_tokens_total"`
	OutputTokensTotal     uint64            `json:"output_tokens_total"`
	CacheReadTokensTotal  uint64            `json:"cache_read_tokens_total"`
	CacheWriteTokensTotal uint64            `json:"cache_write_tokens_total"`
	ContextTokensLast     uint64            `json:"context_tokens_last"`
}

func newSession() Session {
	return Session{
		HooksByKind:     map[string]uint64{},
		DecisionsByKind: map[string]uint64{},
		EventsByType:    map[string]uint64{},
		EventsBySource:  map[string]uint64{},
	}
}

// SessionOptions configures offline session stats loading.
type SessionOptions struct {
	ConfigPath   string
	SessionsRoot string
	Provider     string
	SessionID    string
}

// Load reads config, gates statistics, and scans the session ledger.
func Load(ctx context.Context, opts SessionOptions) (Session, error) {
	store, err := config.Load(ctx, opts.ConfigPath)
	if err != nil {
		return Session{}, fmt.Errorf("load config: %w", err)
	}
	if err := Gate(store.Current().Trajectory); err != nil {
		return Session{}, err
	}
	root := opts.SessionsRoot
	if root == "" {
		root = trajectory.DefaultSessionsDir()
	}
	path, err := trajectory.FindSessionPath(root, opts.Provider, opts.SessionID)
	if err != nil {
		return Session{}, err
	}
	events, err := trajectory.ReadEvents(path)
	if err != nil {
		return Session{}, err
	}
	stats := FromEvents(events)
	if stats.Provider == "" {
		stats.Provider = opts.Provider
	}
	if stats.SessionID == "" {
		stats.SessionID = opts.SessionID
	}
	return stats, nil
}
