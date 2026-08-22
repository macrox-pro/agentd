package importer

import (
	"fmt"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/speakeasy-api/agenthooks/transcript"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

// ResolveCodexTranscriptPath finds a Codex transcript JSONL.
// Prefer explicit --path from the hook transcript_path field; no stable default root.
func ResolveCodexTranscriptPath(sessionID, explicitPath, configuredRoot string) (string, error) {
	return resolveExplicitOrSessionPath(sessionID, explicitPath, configuredRoot, "codex")
}

// ImportCodex reads Codex transcript JSONL into trajectory events (partial L2).
// Uses the same message shape as Claude; thinking only when explicitly present.
func ImportCodex(opts ImportOptions) (ImportResult, error) {
	path, err := ResolveCodexTranscriptPath(opts.SessionID, opts.TranscriptPath, opts.ProjectsRoot)
	if err != nil {
		return ImportResult{}, err
	}
	entries, err := transcript.ReadFile(agenthooks.ProviderCodex, path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("read transcript: %w", err)
	}
	now := time.Now().UTC()
	sid := opts.SessionID
	if sid == "" {
		sid = SessionIDFromTranscriptPath(path)
	}
	events, lastIndex := mapEntriesFrom(opts.StartIndex, entries, func(ent transcript.Entry) []trajectory.Event {
		return mapClaudeStyleEntry(ent, "codex", sid, now, opts.Cfg)
	})
	return ImportResult{
		TranscriptPath: path,
		Events:         events,
		LastLineIndex:  lastIndex,
	}, nil
}
