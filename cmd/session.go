package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(sessionListCmd, sessionShowCmd, sessionExportCmd)
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Inspect and export trajectory session ledgers",
	Long: `Inspect and export append-only hook session ledgers stored under the
agentd state directory.

These commands read local JSONL files and do not require a running daemon.`,
	Example: `  agentd session list --provider claude-code
  agentd session show s1 --provider claude-code
  agentd session export --provider claude-code --out sessions.jsonl`,
}
