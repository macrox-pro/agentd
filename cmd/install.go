package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installProvider, "provider", "", "which coding agent to configure (required)")
	installCmd.Flags().StringVar(&installScope, "scope", "project", "where to install: user, project, or plugin")
	installCmd.Flags().StringVar(&installDir, "dir", "", "directory to install into (default: current or home)")
	_ = installCmd.MarkFlagRequired("provider")
}

var (
	installProvider string
	installScope    string
	installDir      string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install hook settings for a coding agent",
	Long: `Write hook settings so a coding agent calls agentd.

After install, the agent will run "agentd hook run --provider=..." for
configured events. Start the agentd service before relying on hooks.`,
	Example: `  agentd install --provider=claude-code --scope=project
  agentd install --provider=cursor --scope=user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
