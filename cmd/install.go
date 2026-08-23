package cmd

import (
	"errors"
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
	installCmd.Flags().BoolVar(&installGlobal, "global", false, "same as --scope=user (agent home)")
	installCmd.Flags().StringVar(&installDir, "dir", "", "install root (default: project=cwd, user=agent home; codex project=cwd/.codex)")
	_ = installCmd.MarkFlagRequired("provider")
}

var (
	installProvider string
	installScope    string
	installGlobal   bool
	installDir      string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install hook settings for a coding agent",
	Long: `Write hook settings so a coding agent calls agentd.

Install uses the agenthooks argv contract: generated configs invoke
"agentd agenthooks run|serve --provider=...". The documented human CLI
remains "agentd hook run|serve". Start the agentd service before relying
on hooks.

Without --dir: scope=project installs into the current working directory
(codex uses ./.codex); scope=user (or --global) installs into the agent
home directory (for example ~/.cursor or ~/.claude). scope=plugin and
opencode with scope=user require an explicit --dir.`,
	Example: `  agentd install --provider=claude-code --scope=project
  agentd install --provider=cursor --global
  agentd install --provider=opencode --scope=project --dir /path/to/repo`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		scope, err := resolveInstallScope(cmd)
		if err != nil {
			return err
		}
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("resolve executable: %w", err)
		}
		exe, err = filepath.Abs(exe)
		if err != nil {
			return fmt.Errorf("abs executable: %w", err)
		}
		result, err := install.Run(cmd.Context(), install.Options{
			Provider:   installProvider,
			Scope:      scope,
			Dir:        installDir,
			DirFlagSet: cmd.Flags().Changed("dir"),
			Command:    []string{exe},
		})
		if err != nil {
			return mapInstallErr(err)
		}
		return install.WriteReport(cmd.OutOrStdout(), result)
	},
}

func resolveInstallScope(cmd *cobra.Command) (string, error) {
	if !installGlobal {
		return installScope, nil
	}
	if cmd.Flags().Changed("scope") && installScope != "user" {
		return "", fmt.Errorf("--global conflicts with --scope=%s (use --scope=user or omit --scope)", installScope)
	}
	return "user", nil
}

func mapInstallErr(err error) error {
	switch {
	case errors.Is(err, install.ErrDirRequired):
		return fmt.Errorf("--dir is required for scope=plugin or provider=opencode with scope=user")
	case errors.Is(err, install.ErrHomeRequired):
		return fmt.Errorf("home directory is unavailable; set HOME or pass --dir explicitly")
	default:
		return err
	}
}
