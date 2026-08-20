package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Inspect and validate configuration",
	Long: `Offline and online configuration operations.

Validation runs without a daemon. Show and patch may contact the daemon
for merged views and runtime overlay updates.`,
}

var (
	configMerged bool
	configLayer  string
	configCWD    string
	configPatch  string
)

func init() {
	configCmd.AddCommand(configValidateCmd, configShowCmd, configPatchCmd)

	configShowCmd.Flags().BoolVar(&configMerged, "merged", false, "show merged effective config")
	configShowCmd.Flags().StringVar(&configLayer, "layer", "", "show one layer: user, project, runtime")
	configShowCmd.Flags().StringVar(&configCWD, "cwd", "", "project root for layer resolution")

	configPatchCmd.Flags().StringVar(&configPatch, "file", "", "runtime overlay patch YAML file (required)")
	_ = configPatchCmd.MarkFlagRequired("file")
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration YAML offline",
	Long: `Parse and validate agentd YAML without a running daemon.

Suitable for CI and pre-commit hooks. Checks schema, policy, guards, and
dispatch route compilation.`,
	Example: `  agentd config validate
  agentd config validate --config ~/.agentd.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show configuration layers or merged config",
	Long: `Display user, project, runtime, or merged effective configuration.

Use --layer to inspect a single layer. Use --merged for the compiled view
(requires a running daemon unless --cwd is set for offline merge preview).`,
	Example: `  agentd config show --merged
  agentd config show --layer user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}

var configPatchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Patch runtime configuration overlay",
	Long: `Apply a YAML patch to the daemon-managed runtime overlay via gRPC.

Used for learned approvals and temporary blocks. Does not modify
~/.agentd.yaml or project .agentd.yaml files.`,
	Example: `  agentd config patch --file runtime-delta.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
