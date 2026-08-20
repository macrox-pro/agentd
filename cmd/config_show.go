package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configMerged bool
	configLayer  string
	configCWD    string
)

func init() {
	configShowCmd.Flags().BoolVar(&configMerged, "merged", false, "show the effective merged settings")
	configShowCmd.Flags().StringVar(&configLayer, "layer", "", "show one layer: user, project, or runtime")
	configShowCmd.Flags().StringVar(&configCWD, "cwd", "", "project directory used to find project settings")
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
