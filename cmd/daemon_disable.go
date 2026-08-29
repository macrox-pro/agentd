package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/daemon"
)

var daemonDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable agentd login autostart",
	Long: `Remove OS login autostart registration for agentd.

This does not stop a daemon that is already running. Use "agentd daemon stop"
to shut down the running process.`,
	Example: `  agentd daemon disable`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		if err := daemon.Disable(); err != nil {
			return mapDisableErr(err)
		}
		return nil
	},
}

func mapDisableErr(err error) error {
	if errors.Is(err, daemon.ErrAutostartUnsupported) {
		return fmt.Errorf("login autostart is not supported on this platform")
	}
	return err
}
