package importer

import (
	"fmt"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/speakeasy-api/agenthooks/transcript"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

// ResolveCursorTranscriptPath finds a Cursor transcript JSONL.
// Prefer explicit --path; session+root scan only when configured path is set.
func ResolveCursorTranscriptPath(sessionID, explicitPath, configuredRoot string) (string, error) {
	return resolveExplicitOrSessionPath(sessionID, explicitPath, configuredRoot, "cursor")
}

// ImportCursor reads Cursor transcript JSONL into trajectory events (partial L2).
// Maps message text only; never invents thinking or tool outputs.
func ImportCursor(opts ImportOptions) (ImportResult, error) {
	path, err := ResolveCursorTranscriptPath(opts.SessionID, opts.TranscriptPath, opts.ProjectsRoot)
	if err != nil {
		return ImportResult{}, err
	}
	entries, err := transcript.ReadFile(agenthooks.ProviderCursor, path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("read transcript: %w", err)
	}
	now := time.Now().UTC()
	sid := opts.SessionID
	if sid == "" {
		sid = SessionIDFromTranscriptPath(path)
	}
	events, lastIndex := mapEntriesFrom(opts.StartIndex, entries, func(ent transcript.Entry) []trajectory.Event {
		return mapCursorEntry(ent, sid, now, opts.Cfg)
	})
	return ImportResult{
		TranscriptPath: path,
		Events:         events,
		LastLineIndex:  lastIndex,
	}, nil
}
