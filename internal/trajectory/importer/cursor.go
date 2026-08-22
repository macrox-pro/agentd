package importer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/speakeasy-api/agenthooks/transcript"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

// ResolveCursorTranscriptPath finds a Cursor transcript JSONL.
// Prefer explicit path; session+root scan only when configured path is set.
func ResolveCursorTranscriptPath(sessionID, explicitPath, configuredRoot string) (string, error) {
	return resolveExplicitOrSessionPath(sessionID, explicitPath, configuredRoot, "cursor")
}

func resolveExplicitOrSessionPath(sessionID, explicitPath, configuredRoot, label string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err != nil {
			return "", fmt.Errorf("transcript path: %w", err)
		}
		return explicitPath, nil
	}
	if configuredRoot == "" {
		return "", ErrTranscriptRootRequired
	}
	want := sessionID + ".jsonl"
	candidate := filepath.Join(configuredRoot, want)
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}
	var found string
	err := filepath.WalkDir(configuredRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == want {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil && found == "" {
		return "", fmt.Errorf("scan %s transcripts: %w", label, err)
	}
	if found == "" {
		return "", fmt.Errorf("%s transcript not found for session %q under %s", label, sessionID, configuredRoot)
	}
	return found, nil
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

// mapCursorEntry maps Cursor-shaped transcript lines (role at top level).
// Does not invent thinking or tool outputs; skips empty/[REDACTED]-only text.
func mapCursorEntry(ent transcript.Entry, sessionID string, ts time.Time, cfg config.TrajectoryConfig) []trajectory.Event {
	base := baseTranscriptEvent("cursor", sessionID, ts)
	if len(ent.Raw) == 0 {
		return nil
	}
	role := ent.Role
	text := ent.Text
	if text == "" {
		var line struct {
			Role    string `json:"role"`
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal(ent.Raw, &line); err == nil {
			if role == "" {
				role = line.Role
			}
			var blocks []contentBlock
			if json.Unmarshal(line.Message.Content, &blocks) == nil {
				var out []trajectory.Event
				for _, b := range blocks {
					if b.Type != "text" && b.Type != "" {
						continue
					}
					t := trajectory.PrepareTranscriptText(b.Text, cfg)
					if t == "" || t == redactedTranscriptText {
						continue
					}
					out = append(out, mapMessageEntry(base, ent, role, t, "", cfg))
				}
				return out
			}
		}
		return nil
	}
	text = trajectory.PrepareTranscriptText(text, cfg)
	if text == "" || text == redactedTranscriptText {
		return nil
	}
	return []trajectory.Event{mapMessageEntry(base, ent, role, text, "", cfg)}
}
