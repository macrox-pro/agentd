package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

const (
	codexSessionsDirName = "sessions"
	codexHomeEnv         = "CODEX_HOME"
	codexJSONLExt        = ".jsonl"
	codexRolloutPrefix   = "rollout-"
)

// ResolveCodexTranscriptPath finds a Codex rollout JSONL.
// Prefer explicit path; else scan under configuredRoot or default
// $CODEX_HOME/sessions (or ~/.codex/sessions) for files ending in -{sessionID}.jsonl
// (real layout: sessions/YYYY/MM/DD/rollout-<ts>-{sessionID}.jsonl). Newest ModTime wins.
func ResolveCodexTranscriptPath(sessionID, explicitPath, configuredRoot string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err != nil {
			return "", fmt.Errorf("transcript path: %w", err)
		}
		return explicitPath, nil
	}
	root := configuredRoot
	if root == "" {
		var err error
		root, err = defaultCodexSessionsRoot()
		if err != nil {
			return "", err
		}
	}
	suffix := "-" + sessionID + codexJSONLExt
	var best string
	var bestMod time.Time
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, suffix) {
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
		return "", fmt.Errorf("scan codex sessions: %w", err)
	}
	if best == "" {
		return "", fmt.Errorf("codex transcript not found for session %q under %s", sessionID, root)
	}
	return best, nil
}

func defaultCodexSessionsRoot() (string, error) {
	home := os.Getenv(codexHomeEnv)
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home dir: %w", err)
		}
		home = filepath.Join(home, ".codex")
	}
	return filepath.Join(home, codexSessionsDirName), nil
}

// ImportCodex reads Codex rollout JSONL into trajectory events.
func ImportCodex(opts ImportOptions) (ImportResult, error) {
	path, err := ResolveCodexTranscriptPath(opts.SessionID, opts.TranscriptPath, opts.ProjectsRoot)
	if err != nil {
		return ImportResult{}, err
	}
	lines, err := readCodexJSONL(path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("read transcript: %w", err)
	}
	now := time.Now().UTC()
	sid := opts.SessionID
	if sid == "" {
		sid = sessionIDFromCodexRollout(path, lines)
	}
	events, lastIndex := mapIndexed(opts.StartIndex, lines, func(line rawLine) int { return line.Index }, func(line rawLine) []trajectory.Event {
		return mapCodexRolloutLine(line, sid, now, opts.Cfg)
	})
	return ImportResult{
		TranscriptPath: path,
		Events:         events,
		LastLineIndex:  lastIndex,
	}, nil
}

func sessionIDFromCodexRollout(path string, lines []rawLine) string {
	for _, line := range lines {
		if sid := sessionIDFromSessionMeta(line.Raw); sid != "" {
			return sid
		}
	}
	if sid := sessionIDFromRolloutFilename(path); sid != "" {
		return sid
	}
	return SessionIDFromTranscriptPath(path)
}

func sessionIDFromRolloutFilename(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), codexJSONLExt)
	if !strings.HasPrefix(base, codexRolloutPrefix) {
		return ""
	}
	// rollout-<local-ts>-<session_id>
	rest := strings.TrimPrefix(base, codexRolloutPrefix)
	// session id is UUID-like suffix after the last timestamp segment;
	// real names: rollout-2026-07-26T17-33-55-019f9ed8-c891-7dd0-9808-e31c3b38ce48
	// Take everything after the timestamp (date T time) — five hyphen groups after T.
	_, after, ok := strings.Cut(rest, "T")
	if !ok {
		return ""
	}
	afterT := after
	// afterT like 17-33-55-019f9ed8-c891-7dd0-9808-e31c3b38ce48
	parts := strings.Split(afterT, "-")
	if len(parts) < 8 {
		// need HH-MM-SS + 5 uuid segments (019f9ed8-c891-7dd0-9808-e31c3b38ce48)
		return ""
	}
	return strings.Join(parts[3:], "-")
}
