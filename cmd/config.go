package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configValidateCmd, configShowCmd, configPatchCmd, configRecordDecisionCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Validate and inspect agentd settings",
	Long: `Validate and inspect agentd configuration files.

Use validate in CI without a running service. Show prints settings; patch and
record-decision update runtime overrides managed by the service.`,
	Example: `  agentd config validate
  agentd config show --merged
  agentd config patch --file runtime-delta.yaml
  agentd config record-decision --fingerprint sha256:secrets/abc --scope project --project-root .`,
}
