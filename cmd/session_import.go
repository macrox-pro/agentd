package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

var (
	sessionImportProvider string
	sessionImportSession  string
	sessionImportPath     string
	sessionImportOut      string
	sessionImportDryRun   bool
	sessionImportJSON     bool
)

func init() {
	sessionImportCmd.Flags().StringVar(&sessionImportProvider, "provider", "", "provider id (required)")
	sessionImportCmd.Flags().StringVar(&sessionImportSession, "session", "", "session id (required when --path omitted for some providers)")
	sessionImportCmd.Flags().StringVar(&sessionImportPath, "path", "", "explicit transcript JSONL path")
	sessionImportCmd.Flags().StringVar(&sessionImportOut, "out", "", "write parsed events to PATH or stdout (-); does not update session ledger")
	sessionImportCmd.Flags().BoolVar(&sessionImportDryRun, "dry-run", false, "print import summary without writing")
	sessionImportCmd.Flags().BoolVar(&sessionImportJSON, "json", false, "print result as JSON")
	_ = sessionImportCmd.MarkFlagRequired("provider")
}

var sessionImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import provider transcript into session ledger",
	Long: `Pull provider on-disk transcript JSONL into the local session ledger (append-only).

Offline — reads/writes local JSONL under the agentd state directory.
Supported: claude-code, codex (~/.codex/sessions rollouts). Partial (path-first): cursor. Others: explicit none.

Use --out to emit parsed events as JSONL to stdout (-) or a file without updating the ledger
or import checkpoint. Summary goes to stderr when --out is set (pipe-safe).`,
	Example: `  agentd session import --provider claude-code --session s1
  agentd session import --provider codex --session s1
  agentd session import --provider cursor --path /path/to/transcript.jsonl
  agentd session import --provider claude-code --path /path/to/session.jsonl --out -
  agentd session import --provider claude-code --path /path/to/session.jsonl --out - | jq -c .
  agentd session import --provider codex --path /path/to/rollout.jsonl --out /tmp/events.jsonl
  agentd session import --provider codex --path /path/to/rollout.jsonl --dry-run --json`,
	RunE: runSessionImport,
}

func runSessionImport(cmd *cobra.Command, _ []string) error {
	provID, err := provider.Parse(sessionImportProvider)
	if err != nil {
		return err
	}
	status := importer.ProviderImporterStatus(string(provID))
	switch status {
	case importer.ImporterNone:
		return fmt.Errorf("transcript import for provider %q is not supported", provID)
	case importer.ImporterPartial:
		if sessionImportPath == "" {
			return fmt.Errorf("cursor import requires --path")
		}
	default:
		if sessionImportPath == "" && sessionImportSession == "" {
			return fmt.Errorf("import requires --session or --path")
		}
	}

	emitOnly := cmd.Flags().Changed("out")
	if emitOnly {
		if strings.TrimSpace(sessionImportOut) == "" {
			return fmt.Errorf("import --out requires a path or -")
		}
	}

	result, err := importer.ImportSession(cmd.Context(), importer.ImportSessionOptions{
		Provider:       provID,
		SessionID:      sessionImportSession,
		TranscriptPath: sessionImportPath,
		DryRun:         sessionImportDryRun,
		EmitOnly:       emitOnly,
		ConfigPath:     resolveConfigPath(),
	})
	if err != nil {
		return err
	}

	if emitOnly && len(result.Events) > 0 {
		if sessionImportOut == "-" {
			if err := trajectory.WriteEvents(cmd.OutOrStdout(), result.Events); err != nil {
				return err
			}
		} else {
			if err := trajectory.WriteEventsToFile(sessionImportOut, result.Events); err != nil {
				return err
			}
		}
	}

	summaryOut := cmd.OutOrStdout()
	if emitOnly {
		summaryOut = cmd.ErrOrStderr()
	}
	if sessionImportJSON {
		enc := json.NewEncoder(summaryOut)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}
	fmt.Fprintf(summaryOut, "provider=%s session=%s status=%s imported=%d path=%s\n",
		result.Provider, result.SessionID, result.ImporterStatus, result.Imported, result.TranscriptPath)
	return nil
}
