package daemon

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// detach re-execs the binary as a background child that runs with --foreground
// and returns only after Health succeeds (or readiness fails).
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

	readyCtx, cancel := context.WithTimeout(context.Background(), readyTimeout)
	defer cancel()
	if err := waitHealth(readyCtx, opts.Socket); err != nil {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = Stop(stopCtx, opts.Socket, 2*time.Second)
		stopCancel()
		return fmt.Errorf("daemon failed to become ready: %w", err)
	}
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
	if opts.LogLevel != "" {
		args = append(args, "--log-level", opts.LogLevel)
	}
	if opts.LogFile != "" {
		args = append(args, "--log-file", opts.LogFile)
	}
	return args
}
