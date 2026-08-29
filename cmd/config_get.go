package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/config"
)

func init() {
	configGetCmd.Flags().StringVar(&configToggleCWD, "cwd", "", "project directory used to find project config")
}

var configGetCmd = &cobra.Command{
	Use:   "get FEATURE",
	Short: "Show effective on/off state for a curated feature",
	Long: `Print whether a curated feature is on or off after merging defaults,
user, and project config. Runtime overlay is not included.

Does not require a running daemon.`,
	Example: `  agentd config get trajectory
  agentd config get guard-shell --cwd /path/to/repo`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir := configToggleCWD
		if projectDir == "" {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("project directory: %w", err)
			}
			projectDir = wd
		}

		result, err := config.GetToggle(config.GetToggleOptions{
			Name:       args[0],
			UserPath:   resolveConfigPath(),
			ProjectDir: projectDir,
		})
		if errors.Is(err, config.ErrUnknownToggle) {
			return fmt.Errorf("unknown feature %q (valid: %s)", args[0], strings.Join(config.ListToggleNames(), ", "))
		}
		if err != nil {
			return err
		}

		state := "off"
		if result.Enabled {
			state = "on"
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s (%s)\n", result.Name, state, result.Source)
		return err
	},
}
