package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	hookRunProvider    string
	hookRunArgvPayload bool
	hookRunTimeout     time.Duration
)

func init() {
	hookRunCmd.Flags().StringVar(&hookRunProvider, "provider", "", "which coding agent is calling (required)")
	hookRunCmd.Flags().BoolVar(&hookRunArgvPayload, "argv-payload", false, "read the hook payload from argv instead of stdin")
	hookRunCmd.Flags().DurationVar(&hookRunTimeout, "timeout", 0, "maximum time to wait for a decision (0 uses the agent default)")
}

var hookRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Handle a hook event and return a decision",
	Long: `Handle a hook event from a coding agent and print the response on stdout.

Point Claude Code, Cursor, Codex, Gemini CLI, or Kimi Code hook settings at
this command with --provider set for that agent. The payload is read from
stdin unless --argv-payload is set.

Requires a running agentd service ("agentd daemon start").`,
	Example: `  agentd hook run --provider=claude-code
  agentd hook run --provider=cursor --argv-payload`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var payloadArg string
		if hookRunArgvPayload {
			if len(args) < 1 {
				return fmt.Errorf("argv payload required")
			}
			payloadArg = args[0]
		}
		return runHookRun(cmd, hookCLIOpts{
			provider:    hookRunProvider,
			argvPayload: hookRunArgvPayload,
			timeout:     hookRunTimeout,
			payloadArg:  payloadArg,
		})
	},
}
