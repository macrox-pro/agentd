package extract

import (
	"bytes"
	"encoding/json"
	"os"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

const (
	codexHookEventStop    = "Stop"
	codexTokenCountMarker = `"token_count"`
	tailScanWindowBytes   = 1 << 20
)

func init() {
	transcriptScanners[agentdv1.Provider_PROVIDER_CODEX] = scanCodexTranscript
}

func scanCodexTranscript(hookRaw []byte) TokenFields {
	var hook struct {
		HookEventName  string `json:"hook_event_name"`
		TranscriptPath string `json:"transcript_path"`
	}
	if err := json.Unmarshal(hookRaw, &hook); err != nil {
		return TokenFields{}
	}
	if hook.HookEventName != codexHookEventStop || hook.TranscriptPath == "" {
		return TokenFields{}
	}
	line := tailScanJSONL(hook.TranscriptPath, func(line []byte) bool {
		return bytes.Contains(line, []byte(codexTokenCountMarker))
	})
	if len(line) == 0 {
		return TokenFields{}
	}
	var env struct {
		Type    string `json:"type"`
		Payload struct {
			Type string `json:"type"`
			Info struct {
				LastTokenUsage json.RawMessage `json:"last_token_usage"`
			} `json:"info"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(line, &env); err != nil {
		return TokenFields{}
	}
	if env.Type != "event_msg" || env.Payload.Type != "token_count" || len(env.Payload.Info.LastTokenUsage) == 0 {
		return TokenFields{}
	}
	return parseUsageObject(env.Payload.Info.LastTokenUsage)
}

func tailScanJSONL(path string, match func(line []byte) bool) []byte {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.Size() == 0 {
		return nil
	}
	size := info.Size()
	window := min(int64(tailScanWindowBytes), size)
	buf := make([]byte, window)
	if _, err := f.ReadAt(buf, size-window); err != nil {
		return nil
	}
	if window < size {
		if i := bytes.IndexByte(buf, '\n'); i >= 0 {
			buf = buf[i+1:]
		}
	}
	for len(buf) > 0 {
		i := bytes.LastIndexByte(buf, '\n')
		var line []byte
		if i < 0 {
			line = bytes.TrimSpace(buf)
			buf = nil
		} else {
			line = bytes.TrimSpace(buf[i+1:])
			buf = buf[:i]
		}
		if len(line) == 0 {
			continue
		}
		if match(line) {
			out := make([]byte, len(line))
			copy(out, line)
			return out
		}
	}
	return nil
}
