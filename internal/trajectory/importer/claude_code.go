package importer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/speakeasy-api/agenthooks/transcript"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

// Conversation line types in Claude Code session JSONL; other top-level types are metadata.
var claudeConversationTypes = map[string]bool{
	"user":      true,
	"assistant": true,
}

// ResolveClaudeTranscriptPath finds a Claude Code session JSONL file.
func ResolveClaudeTranscriptPath(sessionID, explicitPath, projectsRoot string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err != nil {
			return "", fmt.Errorf("transcript path: %w", err)
		}
		return explicitPath, nil
	}
	root := projectsRoot
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home dir: %w", err)
		}
		root = filepath.Join(home, ".claude", "projects")
	}
	want := sessionID + ".jsonl"
	var best string
	var bestMod time.Time
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != want {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		mod := info.ModTime()
		if best == "" || mod.After(bestMod) {
			best = path
			bestMod = mod
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("scan claude projects: %w", err)
	}
	if best == "" {
		return "", fmt.Errorf("claude transcript not found for session %q under %s", sessionID, root)
	}
	return best, nil
}

// ImportClaude reads Claude transcript JSONL and maps entries to trajectory events.
func ImportClaude(opts ImportOptions) (ImportResult, error) {
	path, err := ResolveClaudeTranscriptPath(opts.SessionID, opts.TranscriptPath, opts.ProjectsRoot)
	if err != nil {
		return ImportResult{}, err
	}
	entries, err := transcript.ReadFile(agenthooks.ProviderClaudeCode, path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("read transcript: %w", err)
	}
	now := time.Now().UTC()
	sid := opts.SessionID
	if sid == "" {
		sid = SessionIDFromTranscriptPath(path)
	}
	events, lastIndex := mapEntriesFrom(opts.StartIndex, entries, func(ent transcript.Entry) []trajectory.Event {
		return mapClaudeStyleEntry(ent, "claude-code", sid, now, opts.Cfg)
	})
	return ImportResult{
		TranscriptPath: path,
		Events:         events,
		LastLineIndex:  lastIndex,
	}, nil
}

type claudeTranscriptLine struct {
	Type      string `json:"type"`
	IsMeta    bool   `json:"isMeta"`
	Timestamp string `json:"timestamp"`
	CWD       string `json:"cwd"`
	Message   struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

// mapClaudeStyleEntry maps Claude-shaped transcript lines (message nested under "message").
func mapClaudeStyleEntry(ent transcript.Entry, provider, sessionID string, fallbackTS time.Time, cfg config.TrajectoryConfig) []trajectory.Event {
	if len(ent.Raw) == 0 {
		return nil
	}
	var line claudeTranscriptLine
	if err := json.Unmarshal(ent.Raw, &line); err != nil {
		if ent.Text == "" {
			return nil
		}
		base := baseTranscriptEvent(provider, sessionID, fallbackTS)
		return []trajectory.Event{mapMessageEntry(base, ent, ent.Role, ent.Text, "", cfg)}
	}
	if line.IsMeta {
		return nil
	}
	if line.Type != "" && !claudeConversationTypes[line.Type] {
		return nil
	}
	ts := parseClaudeTimestamp(line.Timestamp, fallbackTS)
	base := baseTranscriptEvent(provider, sessionID, ts)
	if line.CWD != "" {
		base.CWD = line.CWD
	}
	role := line.Message.Role
	if role == "" {
		role = ent.Role
	}
	if len(line.Message.Content) == 0 {
		if ent.Text != "" {
			return []trajectory.Event{mapMessageEntry(base, ent, role, ent.Text, "", cfg)}
		}
		return nil
	}
	var blocks []contentBlock
	if err := json.Unmarshal(line.Message.Content, &blocks); err == nil && len(blocks) > 0 {
		var out []trajectory.Event
		for _, b := range blocks {
			switch b.Type {
			case "thinking":
				text := b.Thinking
				if text == "" {
					text = b.Text
				}
				if ev, ok := mapThinkingEntry(base, ent, text, cfg); ok {
					out = append(out, ev)
				}
			case "text":
				text := trajectory.PrepareTranscriptText(b.Text, cfg)
				if text == "" {
					continue
				}
				out = append(out, mapMessageEntry(base, ent, role, text, "", cfg))
			case "tool_use":
				out = append(out, mapMessageEntry(base, ent, role, b.Text, b.ID, cfg))
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	if ent.Text != "" {
		return []trajectory.Event{mapMessageEntry(base, ent, role, ent.Text, "", cfg)}
	}
	return nil
}

func parseClaudeTimestamp(raw string, fallback time.Time) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	if t, err := time.Parse(time.RFC3339Nano, raw); err == nil {
		return t.UTC()
	}
	if t, err := time.Parse(time.RFC3339, strings.Replace(raw, "Z", "+00:00", 1)); err == nil {
		return t.UTC()
	}
	return fallback
}
