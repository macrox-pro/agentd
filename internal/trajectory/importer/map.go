package importer

import (
	"encoding/json"
	"time"

	"github.com/speakeasy-api/agenthooks/transcript"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

type contentBlock struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Thinking string `json:"thinking"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

func baseTranscriptEvent(provider, sessionID string, ts time.Time) trajectory.Event {
	return trajectory.Event{
		TS:        ts,
		Provider:  provider,
		SessionID: sessionID,
		Source:    trajectory.SourceTranscript,
		Ignorable: true,
	}
}

func mapMessageEntry(base trajectory.Event, ent transcript.Entry, role, text, toolUseID string, cfg config.TrajectoryConfig) trajectory.Event {
	ev := base
	ev.Type = trajectory.TypeTranscriptMessage
	ev.Data = mustJSON(trajectory.TranscriptMessageData{
		Role:                role,
		Text:                trajectory.PrepareTranscriptText(text, cfg),
		ToolUseID:           toolUseID,
		TranscriptLineIndex: ent.Index,
	})
	return ev
}

func mapThinkingEntry(base trajectory.Event, ent transcript.Entry, text string, cfg config.TrajectoryConfig) (trajectory.Event, bool) {
	text = trajectory.PrepareTranscriptText(text, cfg)
	if text == "" || text == "[REDACTED]" {
		return trajectory.Event{}, false
	}
	ev := base
	ev.Type = trajectory.TypeTranscriptThinking
	ev.Data = mustJSON(trajectory.TranscriptThinkingData{
		Text:                text,
		TranscriptLineIndex: ent.Index,
	})
	return ev, true
}

// mapClaudeStyleEntry maps Claude/Codex-shaped transcript lines (message nested under "message").
func mapClaudeStyleEntry(ent transcript.Entry, provider, sessionID string, ts time.Time, cfg config.TrajectoryConfig) []trajectory.Event {
	base := baseTranscriptEvent(provider, sessionID, ts)
	if len(ent.Raw) == 0 {
		return nil
	}
	var line struct {
		Type    string `json:"type"`
		Message struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(ent.Raw, &line); err != nil {
		if ent.Text == "" {
			return nil
		}
		return []trajectory.Event{mapMessageEntry(base, ent, ent.Role, ent.Text, "", cfg)}
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
			default:
				if b.Text != "" {
					out = append(out, mapMessageEntry(base, ent, role, b.Text, b.ID, cfg))
				}
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
					if t == "" || t == "[REDACTED]" {
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
	if text == "" || text == "[REDACTED]" {
		return nil
	}
	return []trajectory.Event{mapMessageEntry(base, ent, role, text, "", cfg)}
}

func mapEntriesFrom(startIndex int, entries []transcript.Entry, mapFn func(transcript.Entry) []trajectory.Event) ([]trajectory.Event, int) {
	var events []trajectory.Event
	lastIndex := max(startIndex-1, -1)
	for _, ent := range entries {
		// startIndex is exclusive lower bound (0 = fresh import; LastLineIndex+1 to resume).
		if ent.Index < startIndex {
			continue
		}
		events = append(events, mapFn(ent)...)
		lastIndex = ent.Index
	}
	if lastIndex < startIndex-1 {
		lastIndex = startIndex - 1
	}
	return events, lastIndex
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return b
}
