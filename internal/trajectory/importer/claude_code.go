package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"github.com/speakeasy-api/agenthooks/transcript"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

// ResolveClaudeTranscriptPath finds a Claude Code session JSONL file.
func ResolveClaudeTranscriptPath(sessionID, explicitPath, projectsRoot string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err != nil {
			return "", fmt.Errorf("transcript path: %w", err)
		}
		return explicitPath, nil
	}
	if sessionID == "" {
		return "", fmt.Errorf("session id or --path required")
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
	var found string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == want || strings.HasSuffix(d.Name(), want) {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil && found == "" {
		return "", fmt.Errorf("scan claude projects: %w", err)
	}
	if found == "" {
		return "", fmt.Errorf("claude transcript not found for session %q under %s", sessionID, root)
	}
	return found, nil
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
