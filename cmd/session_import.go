package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

var (
	sessionImportProvider string
	sessionImportSession  string
	sessionImportPath     string
	sessionImportDryRun   bool
	sessionImportJSON     bool
)

func init() {
	sessionImportCmd.Flags().StringVar(&sessionImportProvider, "provider", "", "provider id (required)")
	sessionImportCmd.Flags().StringVar(&sessionImportSession, "session", "", "session id (required when --path omitted for some providers)")
	sessionImportCmd.Flags().StringVar(&sessionImportPath, "path", "", "explicit transcript JSONL path")
	sessionImportCmd.Flags().BoolVar(&sessionImportDryRun, "dry-run", false, "print import summary without writing")
	sessionImportCmd.Flags().BoolVar(&sessionImportJSON, "json", false, "print result as JSON")
	_ = sessionImportCmd.MarkFlagRequired("provider")
}

var sessionImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import provider transcript into session ledger",
	Long: `Pull provider on-disk transcript JSONL into the local session ledger (append-only).

Offline — reads/writes local JSONL under the agentd state directory.
Supported: claude-code. Partial (path-first): cursor, codex. Others: explicit none.`,
	Example: `  agentd session import --provider claude-code --session s1
  agentd session import --provider cursor --path /path/to/transcript.jsonl
  agentd session import --provider codex --path /path/to/transcript.jsonl --dry-run --json`,
	RunE: runSessionImport,
}

type importOutput struct {
	Provider       string `json:"provider"`
	SessionID      string `json:"session_id"`
	ImporterStatus string `json:"importer_status"`
	Imported       int    `json:"imported"`
	TranscriptPath string `json:"transcript_path,omitempty"`
	LastLineIndex  int    `json:"last_line_index,omitempty"`
	DryRun         bool   `json:"dry_run,omitempty"`
}

func runSessionImport(cmd *cobra.Command, _ []string) error {
	prov := trajectory.CanonicalProvider(sessionImportProvider)
	status := trajectory.ProviderImporterStatus(prov)
	if status == trajectory.ImporterNone {
		return fmt.Errorf("transcript import for provider %q is not supported (importer status: none)", prov)
	}

	cfg := loadTrajectoryConfigForImport(cmd)
	sessionsRoot := trajectory.DefaultSessionsDir()
	if sessionsRoot == "" {
		return fmt.Errorf("sessions dir unavailable")
	}

	sid := sessionImportSession
	if sid == "" && sessionImportPath != "" {
		sid = sessionIDFromTranscriptPath(sessionImportPath)
	}

	startIndex := 0
	if sid != "" {
		cp, err := trajectory.LoadImportCheckpoint(trajectory.ImportSidecarPath(sessionsRoot, prov, sid))
		if err != nil {
			return fmt.Errorf("load import checkpoint: %w", err)
		}
		if cp.SourcePath != "" {
			startIndex = cp.LastLineIndex + 1
		}
	}

	opts := importer.ImportOptions{
		SessionID:      sid,
		TranscriptPath: sessionImportPath,
		StartIndex:     startIndex,
		Cfg:            cfg,
	}
	var result importer.ImportResult
	var err error
	switch prov {
	case "claude-code":
		opts.ProjectsRoot = cfg.ClaudeImport().Path
		result, err = importer.ImportClaude(opts)
	case "cursor":
		opts.ProjectsRoot = cfg.CursorImport().Path
		result, err = importer.ImportCursor(opts)
	case "codex":
		opts.ProjectsRoot = cfg.CodexImport().Path
		result, err = importer.ImportCodex(opts)
	default:
		return fmt.Errorf("transcript import for provider %q is not supported (importer status: %s)", prov, status)
	}
	if err != nil {
		return err
	}
	if sid == "" {
		sid = sessionIDFromTranscriptPath(result.TranscriptPath)
	}

	out := importOutput{
		Provider:       prov,
		SessionID:      sid,
		ImporterStatus: string(status),
		Imported:       len(result.Events),
		TranscriptPath: result.TranscriptPath,
		LastLineIndex:  result.LastLineIndex,
		DryRun:         sessionImportDryRun,
	}

	if sessionImportDryRun {
		return emitImportOutput(cmd, out)
	}

	key := trajectory.ResolveSessionKey(prov, sid, "", "")
	if err := trajectory.AppendImported(sessionsRoot, key, result.Events); err != nil {
		return fmt.Errorf("append imported events: %w", err)
	}

	st, err := os.Stat(result.TranscriptPath)
	if err != nil {
		return fmt.Errorf("stat transcript: %w", err)
	}
	cp := trajectory.ImportCheckpoint{
		LastLineIndex: result.LastLineIndex,
		SourcePath:    result.TranscriptPath,
		SourceModTime: st.ModTime().UTC(),
	}
	if err := trajectory.SaveImportCheckpoint(trajectory.ImportSidecarPath(sessionsRoot, prov, sid), cp); err != nil {
		return fmt.Errorf("save import checkpoint: %w", err)
	}

	return emitImportOutput(cmd, out)
}

func sessionIDFromTranscriptPath(transcriptPath string) string {
	base := transcriptPath
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '/' || base[i] == '\\' {
			base = base[i+1:]
			break
		}
	}
	if len(base) > 6 && base[len(base)-6:] == ".jsonl" {
		return base[:len(base)-6]
	}
	return base
}

func loadTrajectoryConfigForImport(cmd *cobra.Command) config.TrajectoryConfig {
	path := resolveConfigPath()
	if path != "" {
		store, err := config.Load(cmd.Context(), path)
		if err == nil {
			return store.Current().Trajectory
		}
	}
	return trajectory.DefaultImportConfig()
}

func emitImportOutput(cmd *cobra.Command, out importOutput) error {
	if sessionImportJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "provider=%s session=%s status=%s imported=%d path=%s\n",
		out.Provider, out.SessionID, out.ImporterStatus, out.Imported, out.TranscriptPath)
	return nil
}
