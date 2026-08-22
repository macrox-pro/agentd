package importer

import (
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

// ImportOptions configures a provider transcript import run.
type ImportOptions struct {
	SessionID      string
	TranscriptPath string
	ProjectsRoot   string
	StartIndex     int
	Cfg            config.TrajectoryConfig
}

// ImportResult summarizes one import pass.
type ImportResult struct {
	TranscriptPath string
	Events         []trajectory.Event
	LastLineIndex  int
}

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
		sid = strings.TrimSuffix(filepath.Base(path), ".jsonl")
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

func resolveExplicitOrSessionPath(sessionID, explicitPath, configuredRoot, label string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err != nil {
			return "", fmt.Errorf("transcript path: %w", err)
		}
		return explicitPath, nil
	}
	if sessionID == "" {
		return "", fmt.Errorf("%s import requires --path (or --session with configured import path)", label)
	}
	if configuredRoot == "" {
		return "", fmt.Errorf("%s import requires --path (no stable default transcript root)", label)
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
