package cmd

import (
	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

var (
	sessionExportProvider string
	sessionExportSession  string
	sessionExportOut      string
)

func init() {
	sessionExportCmd.Flags().StringVar(&sessionExportProvider, "provider", "", "filter by provider id")
	sessionExportCmd.Flags().StringVar(&sessionExportSession, "session", "", "export one session id")
	sessionExportCmd.Flags().StringVar(&sessionExportOut, "out", "", "write to file instead of stdout")
}

var sessionExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export session JSONL",
	Long: `Export trajectory session ledger JSONL for external viewers.

Without --session, exports all sessions for --provider (or all providers).`,
	Example: `  agentd session export --provider=claude-code --out ledger.jsonl
  agentd session export --provider=cursor --session=s1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		root := trajectory.DefaultSessionsDir()
		if sessionExportOut != "" {
			return trajectory.ExportToFile(sessionExportOut, root, sessionExportProvider, sessionExportSession)
		}
		return trajectory.Export(cmd.OutOrStdout(), root, sessionExportProvider, sessionExportSession)
	},
}
