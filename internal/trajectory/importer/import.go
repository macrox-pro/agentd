package importer

import (
	"fmt"
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

// Import dispatches to the provider importer and sets ProjectsRoot from cfg when empty.
func Import(provider string, opts ImportOptions) (ImportResult, error) {
	prov := trajectory.CanonicalProvider(provider)
	status := trajectory.ProviderImporterStatus(prov)
	if status == trajectory.ImporterNone {
		return ImportResult{}, fmt.Errorf("transcript import for provider %q is not supported (importer status: none)", prov)
	}
	if opts.ProjectsRoot == "" {
		switch prov {
		case "claude-code":
			opts.ProjectsRoot = opts.Cfg.ClaudeImport().Path
		case "cursor":
			opts.ProjectsRoot = opts.Cfg.CursorImport().Path
		case "codex":
			opts.ProjectsRoot = opts.Cfg.CodexImport().Path
		}
	}
	switch prov {
	case "claude-code":
		return ImportClaude(opts)
	case "cursor":
		return ImportCursor(opts)
	case "codex":
		return ImportCodex(opts)
	default:
		return ImportResult{}, fmt.Errorf("transcript import for provider %q is not supported (importer status: %s)", prov, status)
	}
}

// SessionIDFromTranscriptPath derives a session id from a transcript file path basename.
func SessionIDFromTranscriptPath(transcriptPath string) string {
	base := filepath.Base(transcriptPath)
	return strings.TrimSuffix(base, ".jsonl")
}
