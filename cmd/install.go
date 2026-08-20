package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/install"
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installProvider, "provider", "", "which coding agent to configure (required)")
	installCmd.Flags().StringVar(&installScope, "scope", "project", "where to install: user, project, or plugin")
	installCmd.Flags().StringVar(&installDir, "dir", "", "directory to install into (default: current working directory)")
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

Install uses the agenthooks argv contract: generated configs invoke
"agentd agenthooks run|serve --provider=...". The documented human CLI
remains "agentd hook run|serve". Start the agentd service before relying
on hooks.`,
	Example: `  agentd install --provider=claude-code --scope=project
  agentd install --provider=cursor --scope=user
  agentd install --provider=opencode --scope=project --dir /path/to/repo`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("resolve executable: %w", err)
		}
		exe, err = filepath.Abs(exe)
		if err != nil {
			return fmt.Errorf("abs executable: %w", err)
		}
		return install.Run(cmd.Context(), install.Options{
			Provider: installProvider,
			Scope:    installScope,
			Dir:      installDir,
			Command:  []string{exe},
		})
	},
}
