package importer

import (
	"path/filepath"
	"strings"

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

// SessionIDFromTranscriptPath derives a session id from a transcript file path basename.
func SessionIDFromTranscriptPath(transcriptPath string) string {
	base := filepath.Base(transcriptPath)
	return strings.TrimSuffix(base, ".jsonl")
}
