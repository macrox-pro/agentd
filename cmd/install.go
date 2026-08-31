package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/install"
	"github.com/macrox-pro/agentd/internal/install/tui"
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installProvider, "provider", "", "which coding agent to configure")
	installCmd.Flags().StringVar(&installScope, "scope", "project", "where to install: user, project, or plugin")
	installCmd.Flags().BoolVar(&installGlobal, "global", false, "same as --scope=user (agent home)")
	installCmd.Flags().StringVar(&installDir, "dir", "", "install root (default: project=cwd, user=agent home; codex project=cwd/.codex)")
	installCmd.Flags().BoolVar(&installAllDetected, "all-detected", false, "plan or install all high-confidence detected agents")
	installCmd.Flags().BoolVar(&installYes, "yes", false, "apply changes (required with --all-detected to write)")
	installCmd.Flags().BoolVar(&installDryRun, "dry-run", false, "show planned changes without writing")
}

var (
	installProvider    string
	installScope       string
	installGlobal      bool
	installDir         string
	installAllDetected bool
	installYes         bool
	installDryRun      bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install hook settings for a coding agent",
	Long: `Write hook settings so a coding agent calls agentd.

Install uses the agenthooks argv contract: generated configs invoke
"agentd agenthooks run|serve --provider=...". The documented human CLI
remains "agentd hook run|serve". Start the agentd service before relying
on hooks.

Use --all-detected to discover high-confidence targets (config dirs present).
Without --yes, --all-detected prints a plan only. Use --yes to apply.

Without --dir: scope=project installs into the current working directory
(codex uses ./.codex); scope=user (or --global) installs into the agent
home directory (for example ~/.cursor or ~/.claude). scope=plugin and
opencode with scope=user require an explicit --dir.`,
	Example: `  agentd install --provider=claude-code --scope=project
  agentd install --provider=cursor --global
  agentd install --all-detected
  agentd install --all-detected --yes
  agentd install --provider=cursor --dry-run`,
	SilenceUsage: true,
	RunE: runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	_ = args
	if installProvider == "" && !installAllDetected && tui.Interactive(os.Getenv, os.Stdout) {
		return runInstallWizard(cmd, tui.ModeShort, false, false, "")
	}
	if err := validateInstallFlags(cmd); err != nil {
		return err
	}
	exe, err := resolveInstallExecutable()
	if err != nil {
		return err
	}
	if installAllDetected {
		return runInstallAllDetected(cmd, exe)
	}
	scope, err := resolveInstallScope(cmd)
	if err != nil {
		return err
	}
	result, err := install.Run(cmd.Context(), install.Options{
		Provider:   installProvider,
		Scope:      scope,
		Dir:        installDir,
		DirFlagSet: cmd.Flags().Changed("dir"),
		Command:    []string{exe},
		DryRun:     installDryRun,
	})
	if err != nil {
		return mapInstallErr(err)
	}
	return install.WriteReport(cmd.OutOrStdout(), result)
}

func validateInstallFlags(cmd *cobra.Command) error {
	if installProvider != "" && installAllDetected {
		return fmt.Errorf("--provider and --all-detected are mutually exclusive")
	}
	if installYes && !installAllDetected {
		return fmt.Errorf("--yes requires --all-detected")
	}
	if installDryRun && installYes {
		return fmt.Errorf("--dry-run and --yes are mutually exclusive")
	}
	if installAllDetected {
		return nil
	}
	if installProvider == "" {
		return fmt.Errorf("--provider or --all-detected is required")
	}
	return nil
}

func runInstallAllDetected(cmd *cobra.Command, exe string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home directory: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}
	env := install.DiscoverEnv{Cwd: cwd, Home: home}
	findings, err := install.Discover(cmd.Context(), env)
	if err != nil {
		return err
	}
	targets, err := install.TargetsFromHighConfidence(findings, env)
	if err != nil {
		return mapInstallErr(err)
	}
	if len(targets) == 0 {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), "no high-confidence agents detected; run agentd doctor")
		return err
	}
	apply := installYes
	if !apply {
		entries, err := install.Plan(cmd.Context(), targets, []string{exe})
		if err != nil {
			return err
		}
		for _, e := range entries {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "provider=%s scope=%s dir=%s hooks=%s\n",
				e.Target.Provider, e.Target.Scope, e.Target.Dir, e.Status); err != nil {
				return err
			}
			for _, c := range e.Changes {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  %-9s %s\n", c.State, c.Path); err != nil {
					return err
				}
			}
		}
		return nil
	}
	results, err := install.RunAll(cmd.Context(), targets, []string{exe}, installDryRun)
	if err != nil {
		return mapInstallErr(err)
	}
	return install.WriteReports(cmd.OutOrStdout(), results)
}

func resolveInstallExecutable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable: %w", err)
	}
	return filepath.Abs(exe)
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
	case errors.Is(err, install.ErrCommandRequired):
		return fmt.Errorf("agentd executable path is required")
	default:
		if strings.Contains(err.Error(), "unknown provider") {
			return err
		}
		return err
	}
}
