package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(agenthooksCmd)
	agenthooksCmd.AddCommand(ahRunCmd, ahNotifyCmd, ahServeCmd)
}

// agenthooksCmd is the install-generated argv sentinel (Hidden).
// Generated configs invoke: agentd agenthooks run|notify|serve --provider=...
var agenthooksCmd = &cobra.Command{
	Use:   "agenthooks",
	Short: "Internal argv sentinel used by agenthooks install",
	Long:  `Hidden install sentinel. Prefer "agentd hook run|notify|serve" in docs and manual setup.`,
	Example: `  agentd agenthooks run --provider=claude-code
  agentd agenthooks serve --provider=opencode`,
	Hidden: true,
}
