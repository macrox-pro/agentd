package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/hookclient"
)

var (
	configToggleScope string
	configToggleCWD   string
)

func init() {
	configEnableCmd.Flags().StringVar(&configToggleScope, "scope", "", "config layer to write: user or project (default from feature catalog)")
	configEnableCmd.Flags().StringVar(&configToggleCWD, "cwd", "", "project directory for project scope (default: current directory)")
}

var configEnableCmd = &cobra.Command{
	Use:   "enable FEATURE",
	Short: "Enable a curated agentd feature in user or project config",
	Long: `Enable a curated feature by writing to your user or project agentd config file.

Features are persisted in ~/.agentd.yaml (user scope) or .agentd.yaml in the
project directory (project scope). This does not modify the runtime overlay;
use "agentd config patch" for temporary runtime overrides.

Does not require a running daemon. If the daemon is running, it reloads
automatically when the config file changes.`,
	Example: `  agentd config enable trajectory
  agentd config enable guard-shell --scope project
  agentd config enable trajectory-raw
  agentd config enable trajectory-statistics`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigSetToggle(true),
}

func runConfigSetToggle(enabled bool) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		scope, err := parseToggleScope(configToggleScope)
		if err != nil {
			return err
		}
		projectDir, err := resolveToggleProjectDir(configToggleCWD)
		if err != nil {
			return err
		}
		userPath := resolveConfigPath()
		if scope == config.ToggleScopeUser && userPath == "" {
			return fmt.Errorf("--config path or home directory required for user scope")
		}

		result, err := config.SetToggle(config.SetToggleOptions{
			Name:       args[0],
			Scope:      scope,
			Enabled:    enabled,
			UserPath:   userPath,
			ProjectDir: projectDir,
		})
		if errors.Is(err, config.ErrUnknownToggle) {
			return fmt.Errorf("unknown feature %q (valid: %s)", args[0], strings.Join(config.ListToggleNames(), ", "))
		}
		if errors.Is(err, config.ErrToggleAlreadySet) {
			_, printErr := fmt.Fprintf(cmd.OutOrStdout(), "%s: already %s (%s %s)\n",
				result.Name, toggleStateWord(result.Enabled), result.Scope, result.ConfigPath)
			return printErr
		}
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s (%s %s)\n",
			result.Name, toggleActionWord(result.Enabled), result.Scope, result.ConfigPath)
		if err != nil {
			return err
		}
		maybePrintDaemonReloadHint(cmd)
		return nil
	}
}

func parseToggleScope(raw string) (config.ToggleScope, error) {
	if raw == "" {
		return "", nil
	}
	switch raw {
	case string(config.ToggleScopeUser):
		return config.ToggleScopeUser, nil
	case string(config.ToggleScopeProject):
		return config.ToggleScopeProject, nil
	default:
		return "", fmt.Errorf("unknown scope %q (want user or project)", raw)
	}
}

func resolveToggleProjectDir(cwd string) (string, error) {
	dir := cwd
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("project directory: %w", err)
		}
	}
	return dir, nil
}

func toggleActionWord(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func toggleStateWord(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func maybePrintDaemonReloadHint(cmd *cobra.Command) {
	cli, err := hookclient.DialReady(cmd.Context(), resolveSocket())
	if err != nil {
		return
	}
	_ = cli.Close()
	_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "daemon will reload config automatically")
}
