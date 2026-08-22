package importer

import (
	"fmt"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/provider"
)

type importerEntry struct {
	status       ImporterStatus
	importFn     func(ImportOptions) (ImportResult, error)
	projectsRoot func(config.TrajectoryConfig) string
}

var registry = map[provider.ID]importerEntry{
	provider.ClaudeCode: {
		status:       ImporterSupported,
		importFn:     ImportClaude,
		projectsRoot: func(c config.TrajectoryConfig) string { return c.ClaudeImport().Path },
	},
	provider.Cursor: {
		status:       ImporterPartial,
		importFn:     ImportCursor,
		projectsRoot: func(c config.TrajectoryConfig) string { return c.CursorImport().Path },
	},
	provider.Codex: {
		status:       ImporterSupported,
		importFn:     ImportCodex,
		projectsRoot: func(c config.TrajectoryConfig) string { return c.CodexImport().Path },
	},
}

// Import dispatches to the provider importer and sets ProjectsRoot from cfg when empty.
func Import(id provider.ID, opts ImportOptions) (ImportResult, error) {
	entry, ok := registry[id]
	if !ok || entry.importFn == nil {
		return ImportResult{}, fmt.Errorf("%w: %q", ErrImportNotSupported, id)
	}
	if entry.status == ImporterNone {
		return ImportResult{}, fmt.Errorf("%w: %q", ErrImportNotSupported, id)
	}
	if opts.ProjectsRoot == "" && entry.projectsRoot != nil {
		opts.ProjectsRoot = entry.projectsRoot(opts.Cfg)
	}
	return entry.importFn(opts)
}
