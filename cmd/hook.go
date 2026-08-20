package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(hookCmd)
}

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Agent hook entrypoint (wire proxy)",
	Long: `Commands invoked by coding agents from hooks.json / settings.json.

These commands decode provider wire format, call the daemon via gRPC Invoke,
and encode the response. They contain no guard or routing logic.`,
}

var (
	hookProvider   string
	hookArgvPayload bool
	hookTimeout    time.Duration
)

func init() {
	hookCmd.AddCommand(hookRunCmd, hookNotifyCmd, hookServeCmd)

	for _, c := range []*cobra.Command{hookRunCmd, hookNotifyCmd, hookServeCmd} {
		c.Flags().StringVar(&hookProvider, "provider", "", "agent provider (required)")
	}

	hookRunCmd.Flags().BoolVar(&hookArgvPayload, "argv-payload", false, "read payload from argv (legacy Cursor CLI)")
	hookRunCmd.Flags().DurationVar(&hookTimeout, "timeout", 0, "override hook timeout (0 = from provider config)")
}

var hookRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a blocking hook (stdin mode)",
	Long: `Primary entrypoint for Claude Code, Cursor, Codex, Gemini CLI, and Kimi Code.

Reads hook JSON from stdin (or argv with --argv-payload), forwards to the
daemon, writes provider-correct JSON to stdout and exit code.

The --provider flag is required (flag-first detection per agenthooks).`,
	Example: `  agentd hook run --provider=claude-code
  agentd hook run --provider=cursor --argv-payload`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}

var hookNotifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Codex notify hook (argv JSON)",
	Long: `Handle Codex legacy notify hooks where JSON is passed in argv.

Always mapped to async-only dispatch semantics; the agent does not wait
for side effects.`,
	Example: `  agentd hook notify --provider=codex`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}

var hookServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "OpenCode NDJSON stdio bridge",
	Long: `Long-lived process for the OpenCode plugin shim.

Reads NDJSON frames from stdin and forwards each to the daemon via gRPC.
Only --provider=opencode is supported in this mode.`,
	Example: `  agentd hook serve --provider=opencode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
