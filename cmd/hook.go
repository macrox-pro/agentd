package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/hookedge"
)

func init() {
	rootCmd.AddCommand(hookCmd)
}

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Commands that coding agents call from hook settings",
	Long: `Commands intended for coding-agent hook configuration.

Install "agentd hook run" (or notify/serve where needed) in your agent's
hook settings. These commands forward each event to the running agentd
service and print the response the agent expects.`,
}

var (
	hookProvider    string
	hookArgvPayload bool
	hookTimeout     time.Duration
)

func init() {
	hookCmd.AddCommand(hookRunCmd, hookNotifyCmd, hookServeCmd)

	for _, c := range []*cobra.Command{hookRunCmd, hookNotifyCmd, hookServeCmd} {
		c.Flags().StringVar(&hookProvider, "provider", "", "which coding agent is calling (required)")
	}

	hookRunCmd.Flags().BoolVar(&hookArgvPayload, "argv-payload", false, "read the hook payload from argv instead of stdin")
	hookRunCmd.Flags().DurationVar(&hookTimeout, "timeout", 0, "maximum time to wait for a decision (0 uses the agent default)")
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
		if hookProvider == "" {
			return fmt.Errorf("provider is required")
		}
		var payloadArg string
		if hookArgvPayload {
			if len(args) < 1 {
				return fmt.Errorf("argv payload required")
			}
			payloadArg = args[0]
		}
		code := hookedge.Run(cmd.Context(), hookedge.Options{
			Socket:      resolveSocket(),
			Provider:    hookProvider,
			ArgvPayload: hookArgvPayload,
			Timeout:     hookTimeout,
			PayloadArg:  payloadArg,
			Stdin:       cmd.InOrStdin(),
			Stdout:      cmd.OutOrStdout(),
			Stderr:      cmd.ErrOrStderr(),
		})
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}

var hookNotifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Handle a Codex notify-style hook",
	Long: `Handle Codex notify hooks that pass JSON in argv rather than on stdin.

Use this only when your Codex settings generate a notify-style command.
The agent does not wait for side effects from this path.`,
	Example: `  agentd hook notify --provider=codex`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errors.New("hook notify is not implemented")
	},
}

var hookServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Stay running and handle OpenCode hook frames",
	Long: `Stay running and exchange hook frames with OpenCode over stdin/stdout.

Use this when OpenCode's plugin starts agentd as a long-lived process.
Requires --provider=opencode.`,
	Example: `  agentd hook serve --provider=opencode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errors.New("hook serve is not implemented")
	},
}
