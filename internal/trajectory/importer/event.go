package importer

import (
	"encoding/json"
	"time"

	"github.com/speakeasy-api/agenthooks/transcript"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

const redactedTranscriptText = "[REDACTED]"

type contentBlock struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Thinking string `json:"thinking"`
	ID       string `json:"id"`
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
	if text == "" || text == redactedTranscriptText {
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

func mapIndexed[T any](startIndex int, items []T, indexOf func(T) int, mapFn func(T) []trajectory.Event) ([]trajectory.Event, int) {
	var events []trajectory.Event
	lastIndex := max(startIndex-1, -1)
	for _, item := range items {
		idx := indexOf(item)
		// startIndex is exclusive lower bound (0 = fresh import; LastLineIndex+1 to resume).
		if idx < startIndex {
			continue
		}
		events = append(events, mapFn(item)...)
		lastIndex = idx
	}
	if lastIndex < startIndex-1 {
		lastIndex = startIndex - 1
	}
	return events, lastIndex
}

func mapEntriesFrom(startIndex int, entries []transcript.Entry, mapFn func(transcript.Entry) []trajectory.Event) ([]trajectory.Event, int) {
	return mapIndexed(startIndex, entries, func(ent transcript.Entry) int { return ent.Index }, mapFn)
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return b
}
