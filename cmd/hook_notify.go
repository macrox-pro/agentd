package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	hookNotifyProvider string
	hookNotifyTimeout  time.Duration
)

func init() {
	hookNotifyCmd.Flags().StringVar(&hookNotifyProvider, "provider", "", "which coding agent is calling (required)")
	hookNotifyCmd.Flags().DurationVar(&hookNotifyTimeout, "timeout", 0, "maximum time to wait for a decision")
}

var hookNotifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Handle a Codex notify-style hook",
	Long: `Handle Codex notify hooks that pass JSON in argv rather than on stdin.

Use this only when your Codex settings generate a notify-style command.
The agent does not wait for side effects from this path.`,
	Example:       `  agentd hook notify --provider=codex '{"type":"agent-turn-complete"}'`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("notify payload required")
		}
		return runHookNotify(cmd, hookCLIOpts{
			provider:   hookNotifyProvider,
			timeout:    hookNotifyTimeout,
			payloadArg: args[0],
		})
	},
}
