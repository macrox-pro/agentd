package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/config"
)

var (
	configMerged bool
	configLayer  string
	configCWD    string
)

func init() {
	configShowCmd.Flags().BoolVar(&configMerged, "merged", false, "show the effective merged settings")
	configShowCmd.Flags().StringVar(&configLayer, "layer", "", "show one layer: user, project, or runtime")
	configShowCmd.Flags().StringVar(&configCWD, "cwd", "", "project directory used to find project settings")
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print configuration settings",
	Long: `Print agentd configuration settings.

Use --layer to show one source file. Use --merged for the effective settings
after combining sources. Runs offline (no daemon required).`,
	Example: `  agentd config show --merged
  agentd config show --layer user
  agentd config show --layer project --cwd /path/to/repo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		if configMerged && configLayer != "" {
			return fmt.Errorf("use either --merged or --layer, not both")
		}
		var layer config.Layer
		switch {
		case configMerged || configLayer == "":
			layer = config.LayerMerged
		case configLayer == "user":
			layer = config.LayerUser
		case configLayer == "project":
			layer = config.LayerProject
		case configLayer == "runtime":
			layer = config.LayerRuntime
		default:
			return fmt.Errorf("unknown layer %q (want user, project, or runtime)", configLayer)
		}

		store, err := config.LoadWith(cmd.Context(), config.LoadOptions{
			UserPath:    resolveConfigPath(),
			RuntimePath: config.DefaultRuntimePath(),
		})
		if err != nil {
			return err
		}
		if configCWD != "" {
			if _, err := store.EnsureProject(configCWD, ""); err != nil {
				return err
			}
		}
		raw, err := store.LayerYAML(layer, configCWD, "")
		if err != nil {
			return err
		}
		_, err = cmd.OutOrStdout().Write(raw)
		if len(raw) > 0 && raw[len(raw)-1] != '\n' {
			_, _ = fmt.Fprintln(cmd.OutOrStdout())
		}
		return err
	},
}
