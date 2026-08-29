package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configDisableCmd.Flags().StringVar(&configToggleScope, "scope", "", "config layer to write: user or project (default from feature catalog)")
	configDisableCmd.Flags().StringVar(&configToggleCWD, "cwd", "", "project directory for project scope (default: current directory)")
}

var configDisableCmd = &cobra.Command{
	Use:   "disable FEATURE",
	Short: "Disable a curated agentd feature in user or project config",
	Long: `Disable a curated feature by writing enabled: false to your user or project
agentd config file.

Does not require a running daemon. If the daemon is running, it reloads
automatically when the config file changes.`,
	Example: `  agentd config disable trajectory
  agentd config disable guard-shell --scope project`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigSetToggle(false),
}
