package cmd

import (
	"github.com/spf13/cobra"
)

var configPatch string

func init() {
	configPatchCmd.Flags().StringVar(&configPatch, "file", "", "YAML file with runtime overrides to apply (required)")
	_ = configPatchCmd.MarkFlagRequired("file")
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
