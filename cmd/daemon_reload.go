package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/daemon"
)

var daemonReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload settings from disk",
	Long: `Reload agentd settings from the config file without restarting the service.

Use this after editing your config by hand when you want changes applied
immediately.`,
	Example: `  agentd daemon reload`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		res, err := daemon.Reload(cmd.Context(), resolveSocket())
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "reloaded generation=%d fingerprint=%s\n",
			res.Generation, res.Fingerprint)
		return nil
	},
}
