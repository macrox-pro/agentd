package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

const defaultOpenCodeProvider = "opencode"

var (
	ahServeProvider string
	ahServeTimeout  time.Duration
)

func init() {
	ahServeCmd.Flags().StringVar(&ahServeProvider, "provider", "", "which coding agent is calling (defaults to opencode)")
	ahServeCmd.Flags().DurationVar(&ahServeTimeout, "timeout", 0, "per-frame deadline override")
}

var ahServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "OpenCode NDJSON bridge (install sentinel)",
	Long: `Install-generated alias for "agentd hook serve".

When --provider is omitted, defaults to opencode for generated OpenCode shims.`,
	Example:       `  agentd agenthooks serve --provider=opencode`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		return runHookServe(cmd, hookCLIOpts{
			provider:        ahServeProvider,
			timeout:         ahServeTimeout,
			defaultProvider: defaultOpenCodeProvider,
		})
	},
}
