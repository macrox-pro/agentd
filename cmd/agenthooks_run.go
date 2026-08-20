package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	ahRunProvider    string
	ahRunArgvPayload bool
	ahRunTimeout     time.Duration
)

func init() {
	ahRunCmd.Flags().StringVar(&ahRunProvider, "provider", "", "which coding agent is calling (required)")
	ahRunCmd.Flags().BoolVar(&ahRunArgvPayload, "argv-payload", false, "read the hook payload from argv instead of stdin")
	ahRunCmd.Flags().DurationVar(&ahRunTimeout, "timeout", 0, "maximum time to wait for a decision (0 uses the agent default)")
}

var ahRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Handle a hook event (install sentinel)",
	Long: `Install-generated alias for "agentd hook run".

Same wire path as the public hook command; used when agenthooks writes
argv that invokes the agenthooks sentinel.`,
	Example:       `  agentd agenthooks run --provider=claude-code`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var payloadArg string
		if ahRunArgvPayload {
			if len(args) < 1 {
				return fmt.Errorf("argv payload required")
			}
			payloadArg = args[0]
		}
		return runHookRun(cmd, hookCLIOpts{
			provider:    ahRunProvider,
			argvPayload: ahRunArgvPayload,
			timeout:     ahRunTimeout,
			payloadArg:  payloadArg,
		})
	},
}
