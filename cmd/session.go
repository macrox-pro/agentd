package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(sessionListCmd, sessionShowCmd, sessionExportCmd, sessionSearchCmd, sessionImportCmd, sessionReplayCmd, sessionForkCmd)
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Inspect and export trajectory session ledgers",
	Long: `Inspect and export append-only hook session ledgers stored under the
agentd state directory.

These commands read local JSONL files and do not require a running daemon.`,
	Example: `  agentd session list --provider claude-code
  agentd session show s1 --provider claude-code
  agentd session search --provider claude-code --query thinking
  agentd session import --provider claude-code --session s1 --path /path/to/session.jsonl
  agentd session replay --policy --provider claude-code --session s1
  agentd session fork --provider claude-code --session s1 --new-session s1-fork
  agentd session export --provider claude-code --out sessions.jsonl`,
}
