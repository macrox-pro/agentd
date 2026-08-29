package cmd_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestDaemonEnableDisableCLI(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got := executeRoot(t, execOpts{args: []string{"daemon", "disable"}})
	require.NoError(t, got.err)
}

func TestDaemonEnable_partial_fail_prints_stderr(t *testing.T) {
	t.Cleanup(daemon.SetAutostartHooksForTest(&daemon.AutostartHooks{
		Register:    func(daemon.AutostartSpec) error { return nil },
		ResolveExe:  func() (string, error) { return "/usr/bin/agentd", nil },
		StartIfDown: func(context.Context, daemon.StartOptions) error { return errors.New("bad config") },
		ReadState: func() (daemon.AutostartReport, error) {
			return daemon.AutostartReport{Enabled: true}, nil
		},
	}))
	got := executeRoot(t, execOpts{args: []string{"daemon", "enable"}})
	require.Error(t, got.err)
	assert.Contains(t, got.out, "login autostart is enabled")
	assert.Contains(t, got.out, "daemon did not start now")
}

func TestDaemonEnable_unsupported_platform(t *testing.T) {
	t.Cleanup(daemon.SetAutostartHooksForTest(&daemon.AutostartHooks{
		Register:   func(daemon.AutostartSpec) error { return daemon.ErrAutostartUnsupported },
		ResolveExe: func() (string, error) { return "/usr/bin/agentd", nil },
	}))
	got := executeRoot(t, execOpts{args: []string{"daemon", "enable"}})
	require.Error(t, got.err)
	assert.Contains(t, got.out, "not supported on this platform")
}

func TestDaemonStatusJSON_includes_autostart(t *testing.T) {
	socket, _ := testSocketDir(t)
	got := executeRoot(t, execOpts{args: []string{"daemon", "status", "--json"}, socketPath: socket})
	require.NoError(t, got.err)
	assert.Contains(t, got.out, `"autostart"`)
	assert.Contains(t, got.out, `"enabled"`)
}

func TestDaemonEnable_maps_not_available_linux(t *testing.T) {
	if os.Getenv("GOOS") != "linux" {
		t.Skip("linux-only preflight")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_RUNTIME_DIR", "")
	got := executeRoot(t, execOpts{args: []string{"daemon", "enable"}})
	require.Error(t, got.err)
	assert.Contains(t, got.out, "unavailable")
}
