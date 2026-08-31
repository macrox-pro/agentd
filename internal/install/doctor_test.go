package install_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
)

func TestReport(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	env := testDiscoverEnv(t, root)
	bin := filepath.Join(root, "agentd")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	report, err := install.Report(context.Background(), install.DoctorOptions{
		Cwd:     env.Cwd,
		Home:    env.Home,
		Command: []string{bin},
		Env:     env,
	})
	require.NoError(t, err, "Report(empty)")
	assert.Empty(t, report.Findings)
	assert.False(t, report.DaemonReachable)
}

func TestReport_daemonUnreachableOK(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	root := t.TempDir()
	env := testDiscoverEnv(t, root)
	report, err := install.Report(ctx, install.DoctorOptions{
		Cwd:    env.Cwd,
		Home:   env.Home,
		Socket: filepath.Join(root, "no-such.sock"),
		Env:    env,
	})
	require.NoError(t, err, "Report(unreachable socket)")
	assert.False(t, report.DaemonReachable)
}
