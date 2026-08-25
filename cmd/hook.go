package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/hookedge"
)

func init() {
	rootCmd.AddCommand(hookCmd)
	hookCmd.AddCommand(hookRunCmd, hookNotifyCmd, hookServeCmd)
}

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Commands that coding agents call from hook settings",
	Long: `Commands intended for coding-agent hook configuration.

Install "agentd hook run" (or notify/serve where needed) in your agent's
hook settings. These commands forward each event to the running agentd
service and print the response the agent expects.

When the daemon is unreachable, policy.offline applies (default fail_open).`,
	Example: `  agentd hook run --provider=claude-code
  agentd hook serve --provider=opencode`,
}

type hookCLIOpts struct {
	provider        string
	payloadArg      string
	argvPayload     bool
	timeout         time.Duration
	defaultProvider string
}

func runHookRun(cmd *cobra.Command, o hookCLIOpts) error {
	code := hookedge.Run(cmd.Context(), hookedge.Options{
		Socket:      resolveSocket(),
		ConfigPath:  resolveConfigPath(),
		Provider:    o.provider,
		ArgvPayload: o.argvPayload,
		Timeout:     o.timeout,
		PayloadArg:  o.payloadArg,
		Stdin:       cmd.InOrStdin(),
		Stdout:      cmd.OutOrStdout(),
		Stderr:      cmd.ErrOrStderr(),
	})
	if code != 0 {
		os.Exit(code)
	}
	return nil
}

func runHookNotify(cmd *cobra.Command, o hookCLIOpts) error {
	if o.payloadArg == "" {
		return fmt.Errorf("notify payload required")
	}
	code := hookedge.Notify(cmd.Context(), hookedge.Options{
		Socket:     resolveSocket(),
		ConfigPath: resolveConfigPath(),
		Provider:   o.provider,
		Timeout:    o.timeout,
		PayloadArg: o.payloadArg,
		Stdout:     cmd.OutOrStdout(),
		Stderr:     cmd.ErrOrStderr(),
	})
	if code != 0 {
		os.Exit(code)
	}
	return nil
}

func runHookServe(cmd *cobra.Command, o hookCLIOpts) error {
	provider := o.provider
	if provider == "" {
		provider = o.defaultProvider
	}
	code := hookedge.Serve(cmd.Context(), hookedge.Options{
		Socket:     resolveSocket(),
		ConfigPath: resolveConfigPath(),
		Provider:   provider,
		Timeout:    o.timeout,
		Stdin:      cmd.InOrStdin(),
		Stdout:     cmd.OutOrStdout(),
		Stderr:     cmd.ErrOrStderr(),
	})
	if code != 0 {
		os.Exit(code)
	}
	return nil
}
