package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	hookServeProvider string
	hookServeTimeout  time.Duration
)

func init() {
	hookServeCmd.Flags().StringVar(&hookServeProvider, "provider", "", "which coding agent is calling (required)")
	hookServeCmd.Flags().DurationVar(&hookServeTimeout, "timeout", 0, "per-frame deadline override")
}

var hookServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Stay running and handle OpenCode hook frames",
	Long: `Stay running and exchange hook frames with OpenCode over stdin/stdout.

Use this when OpenCode's plugin starts agentd as a long-lived process.
Requires --provider=opencode. Install-generated shims may also invoke the
hidden "agentd agenthooks serve" sentinel with the same behavior.`,
	Example:       `  agentd hook serve --provider=opencode`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		return runHookServe(cmd, hookCLIOpts{
			provider: hookServeProvider,
			timeout:  hookServeTimeout,
		})
	},
}
