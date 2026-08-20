package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/config"
)

var configValidateCWD string

func init() {
	configValidateCmd.Flags().StringVar(&configValidateCWD, "cwd", "", "project directory used to find project settings")
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check that configuration files are valid",
	Long: `Check that agentd configuration YAML is valid.

Runs without a background service, so it is safe for CI and pre-commit hooks.
With --cwd, also merges the nearest project .agentd.yaml.`,
	Example: `  agentd config validate
  agentd config validate --config ~/.agentd.yaml
  agentd config validate --cwd /path/to/repo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		store, err := config.Load(cmd.Context(), resolveConfigPath())
		if err != nil {
			return err
		}
		if configValidateCWD != "" {
			if _, err := store.EnsureProject(configValidateCWD, ""); err != nil {
				return err
			}
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "ok")
		return nil
	},
}
