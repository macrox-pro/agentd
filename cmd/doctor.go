package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/install"
)

var (
	doctorCWD  string
	doctorJSON bool
)

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().StringVar(&doctorCWD, "cwd", "", "working directory for discovery (default: current directory)")
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "emit JSON report")
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Report detected coding agents and hook install status",
	Long: `Discover installed coding agents and show hook install status.

Read-only: does not modify agent configuration. When the daemon socket is
reachable, the report includes daemon reachability.`,
	Example: `  agentd doctor
  agentd doctor --json
  agentd doctor --cwd /path/to/repo`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		cwd, err := resolveDoctorCWD(doctorCWD)
		if err != nil {
			return err
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("home directory: %w", err)
		}
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("resolve executable: %w", err)
		}
		exe, err = filepath.Abs(exe)
		if err != nil {
			return fmt.Errorf("abs executable: %w", err)
		}
		report, err := install.Report(cmd.Context(), install.DoctorOptions{
			Cwd:     cwd,
			Home:    home,
			Socket:  resolveSocket(),
			Command: []string{exe},
			Env: install.DiscoverEnv{
				Cwd:  cwd,
				Home: home,
			},
		})
		if err != nil {
			return err
		}
		if doctorJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(report)
		}
		return writeDoctorHuman(cmd.OutOrStdout(), report)
	},
}

func resolveDoctorCWD(flag string) (string, error) {
	if flag != "" {
		return filepath.Abs(flag)
	}
	return os.Getwd()
}

func writeDoctorHuman(w io.Writer, report install.DoctorReport) error {
	if len(report.Findings) == 0 {
		if _, err := fmt.Fprintln(w, "no coding agents detected"); err != nil {
			return err
		}
	} else {
		for _, f := range report.Findings {
			line := fmt.Sprintf("%s confidence=%s", f.Provider, f.Confidence)
			if f.ProjectDir != "" {
				line += fmt.Sprintf(" project=%s", f.ProjectDir)
			}
			if f.UserDir != "" {
				line += fmt.Sprintf(" user=%s", f.UserDir)
			}
			if f.Binary != "" {
				line += fmt.Sprintf(" binary=%s", f.Binary)
			}
			if f.Note != "" {
				line += fmt.Sprintf(" (%s)", f.Note)
			}
			if f.Plan != nil {
				line += fmt.Sprintf(" hooks=%s", f.Plan.Status)
			}
			if _, err := fmt.Fprintln(w, line); err != nil {
				return err
			}
		}
	}
	if report.DaemonReachable {
		_, err := fmt.Fprintln(w, "daemon: reachable")
		return err
	}
	_, err := fmt.Fprintln(w, "daemon: unreachable")
	return err
}
