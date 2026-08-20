package daemon

import (
	"fmt"
	"os"
	"os/exec"
)

// detach re-execs the binary as a background child that runs with --foreground.
func detach(opts StartOptions) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("executable: %w", err)
	}
	cmd := exec.Command(exe, foregroundArgs(opts)...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := configureDetach(cmd); err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("detach start: %w", err)
	}
	_ = cmd.Process.Release()
	return nil
}

func foregroundArgs(opts StartOptions) []string {
	args := []string{"daemon", "start", "--foreground"}
	if opts.ConfigPath != "" {
		args = append(args, "--config", opts.ConfigPath)
	}
	if opts.Socket != "" {
		args = append(args, "--socket", opts.Socket)
	}
	return args
}
