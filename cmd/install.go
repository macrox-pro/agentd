package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installProvider, "provider", "", "agent provider (required)")
	installCmd.Flags().StringVar(&installScope, "scope", "project", "install scope: user, project, plugin")
	installCmd.Flags().StringVar(&installDir, "dir", "", "target directory (default: cwd or home)")
	_ = installCmd.MarkFlagRequired("provider")
}

var (
	installProvider string
	installScope    string
	installDir      string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install hook configs for a coding agent",
	Long: `Write hooks.json / settings.json for the given provider using agenthooks/install.

Generated configs invoke "agentd hook run --provider=..." so hooks connect
to this daemon. Installation is explicit; the daemon does not auto-install.`,
	Example: `  agentd install --provider=claude-code --scope=project
  agentd install --provider=cursor --scope=user --dir ~`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd
		_ = args
		return errNotImplemented
	},
}
