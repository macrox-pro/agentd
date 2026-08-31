package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/install/tui"
)

var (
	setupCWD    string
	setupYes    bool
	setupDryRun bool
)

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVar(&setupCWD, "cwd", "", "working directory for discovery (default: current directory)")
	setupCmd.Flags().BoolVar(&setupYes, "yes", false, "apply without confirmation")
	setupCmd.Flags().BoolVar(&setupDryRun, "dry-run", false, "preview only; do not write")
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive wizard to discover and install agent hooks",
	Long: `Discover coding agents and install agentd hooks interactively.

Requires a TTY. Set AGENTD_NO_TUI=1 or CI=true to disable the wizard.`,
	Example: `  agentd setup
  agentd setup --yes
  agentd setup --dry-run`,
	SilenceUsage: true,
	RunE:         runSetup,
}

func runSetup(cmd *cobra.Command, args []string) error {
	_ = args
	if !tui.Interactive(os.Getenv, os.Stdout) {
		return fmt.Errorf("%w: requires a TTY; use --provider or --all-detected", tui.ErrNonInteractive)
	}
	if setupYes && setupDryRun {
		return fmt.Errorf("--dry-run and --yes are mutually exclusive")
	}
	return runInstallWizard(cmd, tui.ModeFull, setupYes, setupDryRun, setupCWD)
}
