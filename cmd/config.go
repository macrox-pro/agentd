package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var errNotImplemented = errors.New("not implemented")

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configValidateCmd, configShowCmd, configPatchCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Validate and inspect agentd settings",
	Long: `Validate and inspect agentd configuration files.

Use validate in CI without a running service. Show prints settings; patch
updates runtime overrides managed by the service.`,
	Example: `  agentd config validate
  agentd config show --merged`,
}
