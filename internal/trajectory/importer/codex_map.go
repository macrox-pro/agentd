package importer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

const (
	codexJSONLMaxToken = 16 << 20

	lineTypeSessionMeta  = "session_meta"
	lineTypeEventMsg     = "event_msg"
	lineTypeResponseItem = "response_item"

	eventUserMessage    = "user_message"
	eventAgentMessage   = "agent_message"
	eventAgentReasoning = "agent_reasoning"

	payloadTypeMessage              = "message"
	payloadTypeReasoning            = "reasoning"
	payloadTypeFunctionCall         = "function_call"
	payloadTypeFunctionCallOutput   = "function_call_output"
	payloadTypeCustomToolCall       = "custom_tool_call"
	payloadTypeCustomToolCallOutput = "custom_tool_call_output"

	roleUser      = "user"
	roleAssistant = "assistant"
)

// rawLine is one JSONL line with its zero-based index.
type rawLine struct {
	Index int
	Raw   json.RawMessage
}

func readCodexJSONL(path string) ([]rawLine, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64<<10), codexJSONLMaxToken)
	var out []rawLine
	i := 0
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		raw := make([]byte, len(line))
		copy(raw, line)
		out = append(out, rawLine{Index: i, Raw: raw})
		i++
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan jsonl: %w", err)
	}
	return out, nil
}

func mapCodexRolloutLine(line rawLine, sessionID string, ts time.Time, cfg config.TrajectoryConfig) []trajectory.Event {
	if len(line.Raw) == 0 {
		return nil
	}
	var env struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(line.Raw, &env); err != nil {
		return nil
	}
	base := baseTranscriptEvent("codex", sessionID, ts)
	switch env.Type {
	case lineTypeEventMsg:
		return mapCodexEventMsg(base, line, env.Payload, cfg)
	case lineTypeResponseItem:
		return mapCodexResponseItem(base, line, env.Payload, cfg)
	default:
		return nil
	}
}

func mapCodexEventMsg(base trajectory.Event, line rawLine, payload json.RawMessage, cfg config.TrajectoryConfig) []trajectory.Event {
	var p struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Text    string `json:"text"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil
	}
	switch p.Type {
	case eventUserMessage:
		text := trajectory.PrepareTranscriptText(p.Message, cfg)
		if text == "" {
			return nil
		}
		return []trajectory.Event{mapCodexMessage(base, line.Index, roleUser, text, "")}
	case eventAgentMessage:
		text := trajectory.PrepareTranscriptText(p.Message, cfg)
		if text == "" {
			return nil
		}
		return []trajectory.Event{mapCodexMessage(base, line.Index, roleAssistant, text, "")}
	case eventAgentReasoning:
		text := strings.TrimSpace(p.Text)
		text = trajectory.PrepareTranscriptText(text, cfg)
		if text == "" || text == redactedTranscriptText {
			return nil
		}
		ev := base
		ev.Type = trajectory.TypeTranscriptThinking
		ev.Data = mustJSON(trajectory.TranscriptThinkingData{
			Text:                text,
			TranscriptLineIndex: line.Index,
		})
		return []trajectory.Event{ev}
	default:
		return nil
	}
}

func mapCodexResponseItem(base trajectory.Event, line rawLine, payload json.RawMessage, cfg config.TrajectoryConfig) []trajectory.Event {
	var p struct {
		Type      string `json:"type"`
		CallID    string `json:"call_id"`
		Arguments string `json:"arguments"`
		Output    string `json:"output"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil
	}
	switch p.Type {
	case payloadTypeFunctionCall, payloadTypeCustomToolCall:
		if p.CallID == "" {
			return nil
		}
		text := trajectory.PrepareTranscriptText(p.Arguments, cfg)
		return []trajectory.Event{mapCodexMessage(base, line.Index, roleAssistant, text, p.CallID)}
	case payloadTypeFunctionCallOutput, payloadTypeCustomToolCallOutput:
		if p.CallID == "" {
			return nil
		}
		text := trajectory.PrepareTranscriptText(p.Output, cfg)
		return []trajectory.Event{mapCodexMessage(base, line.Index, roleAssistant, text, p.CallID)}
	case payloadTypeReasoning, payloadTypeMessage:
		// Encrypted reasoning and response_item messages (incl. developer/injections) are skipped;
		// conversational text comes from event_msg only.
		return nil
	default:
		return nil
	}
}

func mapCodexMessage(base trajectory.Event, lineIndex int, role, text, toolUseID string) trajectory.Event {
	ev := base
	ev.Type = trajectory.TypeTranscriptMessage
	ev.Data = mustJSON(trajectory.TranscriptMessageData{
		Role:                role,
		Text:                text,
		ToolUseID:           toolUseID,
		TranscriptLineIndex: lineIndex,
	})
	return ev
}

func sessionIDFromSessionMeta(raw json.RawMessage) string {
	var env struct {
		Type    string `json:"type"`
		Payload struct {
			SessionID string `json:"session_id"`
			ID        string `json:"id"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return ""
	}
	if env.Type != lineTypeSessionMeta {
		return ""
	}
	if env.Payload.SessionID != "" {
		return env.Payload.SessionID
	}
	return env.Payload.ID
}
