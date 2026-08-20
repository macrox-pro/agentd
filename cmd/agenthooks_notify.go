package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	ahNotifyProvider string
	ahNotifyTimeout  time.Duration
)

func init() {
	ahNotifyCmd.Flags().StringVar(&ahNotifyProvider, "provider", "", "which coding agent is calling (required)")
	ahNotifyCmd.Flags().DurationVar(&ahNotifyTimeout, "timeout", 0, "maximum time to wait for a decision")
}

var ahNotifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Handle a Codex notify hook (install sentinel)",
	Long: `Install-generated alias for "agentd hook notify".

Same wire path as the public notify command.`,
	Example:       `  agentd agenthooks notify --provider=codex '{"type":"agent-turn-complete"}'`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("notify payload required")
		}
		return runHookNotify(cmd, hookCLIOpts{
			provider:   ahNotifyProvider,
			timeout:    ahNotifyTimeout,
			payloadArg: args[0],
		})
	},
}
