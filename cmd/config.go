package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Validate and inspect agentd settings",
	Long: `Validate and inspect agentd configuration files.

Use validate in CI without a running service. Show prints settings; patch
updates runtime overrides managed by the service.`,
}

var (
	configMerged bool
	configLayer  string
	configCWD    string
	configPatch  string
)

func init() {
	configCmd.AddCommand(configValidateCmd, configShowCmd, configPatchCmd)

	configShowCmd.Flags().BoolVar(&configMerged, "merged", false, "show the effective merged settings")
	configShowCmd.Flags().StringVar(&configLayer, "layer", "", "show one layer: user, project, or runtime")
	configShowCmd.Flags().StringVar(&configCWD, "cwd", "", "project directory used to find project settings")

	configPatchCmd.Flags().StringVar(&configPatch, "file", "", "YAML file with runtime overrides to apply (required)")
	_ = configPatchCmd.MarkFlagRequired("file")
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check that configuration files are valid",
	Long: `Check that agentd configuration YAML is valid.

Runs without a background service, so it is safe for CI and pre-commit hooks.`,
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
	Short: "Print configuration settings",
	Long: `Print agentd configuration settings.

Use --layer to show one source file. Use --merged for the effective settings
after combining sources.`,
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
	Short: "Apply temporary runtime setting overrides",
	Long: `Apply temporary overrides to the running service without editing your
user or project config files.

Requires a running agentd service.`,
	Example: `  agentd config patch --file runtime-delta.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
