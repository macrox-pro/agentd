package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/install"
	"github.com/macrox-pro/agentd/internal/install/tui"
)

func runInstallWizard(cmd *cobra.Command, mode tui.Mode, yes, dryRun bool, cwdFlag string) error {
	cwd, err := resolveDoctorCWD(cwdFlag)
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home directory: %w", err)
	}
	exe, err := resolveInstallExecutable()
	if err != nil {
		return err
	}
	env := install.DiscoverEnv{Cwd: cwd, Home: home}
	findings, err := install.Discover(cmd.Context(), env)
	if err != nil {
		return err
	}
	command := []string{exe}
	deps := tui.Deps{
		Plan: func(ctx context.Context, targets []install.Target) ([]install.PlanEntry, error) {
			return install.Plan(ctx, targets, command)
		},
		Apply: func(ctx context.Context, targets []install.Target) ([]install.Result, error) {
			return install.RunAll(ctx, targets, command, false)
		},
	}
	return tui.RunWizard(cmd.Context(), findings, deps, tui.WizardOptions{
		Mode:   mode,
		Yes:    yes,
		DryRun: dryRun,
		Env:    env,
		Out:    cmd.OutOrStdout(),
	})
}
