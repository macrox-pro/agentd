package cmd

import (
	"github.com/spf13/cobra"
)

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
