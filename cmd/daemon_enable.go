package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/daemon"
	"github.com/macrox-pro/agentd/internal/version"
)

const enablePartialFailNotice = "agentd: login autostart is enabled; the daemon did not start now. Fix the error below, then run \"agentd daemon start\" or log in again. Check \"agentd daemon status --json\" (autostart.enabled)."

var daemonEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable agentd to start automatically at login",
	Long: `Register agentd with your OS so it starts when you log in, then start
the daemon now if it is not already running.

Disabling autostart later does not stop a running daemon — use "agentd daemon disable".`,
	Example: `  agentd daemon enable`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		err := daemon.Enable(cmd.Context(), daemon.AutostartOptions{
			Socket:     resolveSocket(),
			ConfigPath: resolveConfigPath(),
			Version:    version.String(),
		})
		if err != nil {
			if rep, stErr := daemon.AutostartStatus(); stErr == nil && rep.Enabled {
				fmt.Fprintln(cmd.ErrOrStderr(), enablePartialFailNotice)
			}
			return mapEnableErr(err)
		}
		return nil
	},
}

func mapEnableErr(err error) error {
	switch {
	case errors.Is(err, daemon.ErrAutostartUnsupported):
		return fmt.Errorf("login autostart is not supported on this platform")
	case errors.Is(err, daemon.ErrAutostartNotAvailable):
		return fmt.Errorf("login autostart is unavailable (user systemd session not running; try after login or enable systemd user session)")
	default:
		return err
	}
}
