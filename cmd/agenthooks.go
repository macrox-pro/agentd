package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/hookedge"
)

// agenthooksCmd is the install-generated argv sentinel (Hidden).
// Generated configs invoke: agentd agenthooks run|notify|serve --provider=...
var agenthooksCmd = &cobra.Command{
	Use:    "agenthooks",
	Short:  "Internal argv sentinel used by agenthooks install",
	Hidden: true,
}

var (
	ahProvider    string
	ahArgvPayload bool
	ahTimeout     time.Duration
)

func init() {
	rootCmd.AddCommand(agenthooksCmd)
	agenthooksCmd.AddCommand(ahRunCmd, ahNotifyCmd, ahServeCmd)

	for _, c := range []*cobra.Command{ahRunCmd, ahNotifyCmd, ahServeCmd} {
		c.Flags().StringVar(&ahProvider, "provider", "", "which coding agent is calling (required)")
	}
	ahRunCmd.Flags().BoolVar(&ahArgvPayload, "argv-payload", false, "read the hook payload from argv instead of stdin")
	ahRunCmd.Flags().DurationVar(&ahTimeout, "timeout", 0, "maximum time to wait for a decision (0 uses the agent default)")
	ahNotifyCmd.Flags().DurationVar(&ahTimeout, "timeout", 0, "maximum time to wait for a decision")
	ahServeCmd.Flags().DurationVar(&ahTimeout, "timeout", 0, "per-frame deadline override")
}

var ahRunCmd = &cobra.Command{
	Use:           "run",
	Short:         "Handle a hook event (install sentinel)",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ahProvider == "" {
			return fmt.Errorf("provider is required")
		}
		var payloadArg string
		if ahArgvPayload {
			if len(args) < 1 {
				return fmt.Errorf("argv payload required")
			}
			payloadArg = args[0]
		}
		code := hookedge.Run(cmd.Context(), hookedge.Options{
			Socket:      resolveSocket(),
			Provider:    ahProvider,
			ArgvPayload: ahArgvPayload,
			Timeout:     ahTimeout,
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

var ahNotifyCmd = &cobra.Command{
	Use:           "notify",
	Short:         "Handle a Codex notify hook (install sentinel)",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ahProvider == "" {
			return fmt.Errorf("provider is required")
		}
		if len(args) < 1 {
			return fmt.Errorf("notify payload required")
		}
		code := hookedge.Notify(cmd.Context(), hookedge.Options{
			Socket:     resolveSocket(),
			Provider:   ahProvider,
			Timeout:    ahTimeout,
			PayloadArg: args[0],
			Stdout:     cmd.OutOrStdout(),
			Stderr:     cmd.ErrOrStderr(),
		})
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}

var ahServeCmd = &cobra.Command{
	Use:           "serve",
	Short:         "OpenCode NDJSON bridge (install sentinel)",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		if ahProvider == "" {
			ahProvider = "opencode"
		}
		code := hookedge.Serve(cmd.Context(), hookedge.Options{
			Socket:   resolveSocket(),
			Provider: ahProvider,
			Timeout:  ahTimeout,
			Stdin:    cmd.InOrStdin(),
			Stdout:   cmd.OutOrStdout(),
			Stderr:   cmd.ErrOrStderr(),
		})
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}
