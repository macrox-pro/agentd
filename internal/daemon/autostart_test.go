package daemon_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func testAutostartHooks(t *testing.T, h *daemon.AutostartHooks) {
	t.Helper()
	t.Cleanup(daemon.SetAutostartHooksForTest(h))
}

func TestAutostart_already_enabled_refresh_manifest(t *testing.T) {
	var calls []daemon.AutostartSpec
	testAutostartHooks(t, &daemon.AutostartHooks{
		Register: func(spec daemon.AutostartSpec) error {
			calls = append(calls, spec)
			return nil
		},
		ResolveExe:  func() (string, error) { return "/usr/bin/agentd", nil },
		StartIfDown: func(context.Context, daemon.StartOptions) error { return nil },
	})
	require.NoError(t, daemon.Enable(context.Background(), daemon.AutostartOptions{}))
	require.NoError(t, daemon.Enable(context.Background(), daemon.AutostartOptions{}))
	require.Len(t, calls, 2)
}

func TestAutostart_disable_idempotent(t *testing.T) {
	n := 0
	testAutostartHooks(t, &daemon.AutostartHooks{
		Unregister: func() error {
			n++
			return nil
		},
	})
	require.NoError(t, daemon.Disable())
	require.NoError(t, daemon.Disable())
	assert.Equal(t, 2, n)
}

func TestAutostart_enable_ignores_already_running(t *testing.T) {
	startCalled := false
	testAutostartHooks(t, &daemon.AutostartHooks{
		Register:    func(daemon.AutostartSpec) error { return nil },
		ResolveExe:  func() (string, error) { return "/usr/bin/agentd", nil },
		StartIfDown: func(context.Context, daemon.StartOptions) error { startCalled = true; return nil },
	})
	require.NoError(t, daemon.Enable(context.Background(), daemon.AutostartOptions{}))
	assert.True(t, startCalled)
}

func TestAutostart_start_fails_autostart_still_enabled(t *testing.T) {
	testAutostartHooks(t, &daemon.AutostartHooks{
		Register:    func(daemon.AutostartSpec) error { return nil },
		ResolveExe:  func() (string, error) { return "/usr/bin/agentd", nil },
		StartIfDown: func(context.Context, daemon.StartOptions) error { return errors.New("bad config") },
		ReadState: func() (daemon.AutostartReport, error) {
			return daemon.AutostartReport{Enabled: true, RegisteredExe: "/usr/bin/agentd"}, nil
		},
	})
	err := daemon.Enable(context.Background(), daemon.AutostartOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "daemon start")
	rep, err := daemon.AutostartStatus()
	require.NoError(t, err)
	assert.True(t, rep.Enabled)
}

func TestAutostart_stale_exe_after_go_install(t *testing.T) {
	testAutostartHooks(t, &daemon.AutostartHooks{
		ReadState: func() (daemon.AutostartReport, error) {
			return daemon.AutostartReport{
				Enabled:       true,
				RegisteredExe: "/old/path/agentd",
			}, nil
		},
	})
	rep, err := daemon.AutostartStatus()
	require.NoError(t, err)
	assert.True(t, rep.Stale)
}

func TestAutostart_home_unavailable(t *testing.T) {
	testAutostartHooks(t, &daemon.AutostartHooks{
		Register:   func(daemon.AutostartSpec) error { return errors.New("home directory unavailable") },
		ResolveExe: func() (string, error) { return "/usr/bin/agentd", nil },
	})
	err := daemon.Enable(context.Background(), daemon.AutostartOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "home directory unavailable")
}

func TestAutostart_exe_resolve_fails(t *testing.T) {
	testAutostartHooks(t, &daemon.AutostartHooks{
		ResolveExe: func() (string, error) { return "", errors.New("executable: no such file") },
	})
	err := daemon.Enable(context.Background(), daemon.AutostartOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "executable")
}

func TestAutostart_disable_does_not_stop_running_daemon(t *testing.T) {
	startCalled := false
	testAutostartHooks(t, &daemon.AutostartHooks{
		Unregister: func() error { return nil },
		StartIfDown: func(context.Context, daemon.StartOptions) error {
			startCalled = true
			return nil
		},
	})
	require.NoError(t, daemon.Disable())
	assert.False(t, startCalled)
}

func TestAutostartStatus_when_disabled(t *testing.T) {
	testAutostartHooks(t, &daemon.AutostartHooks{
		ReadState: func() (daemon.AutostartReport, error) {
			return daemon.AutostartReport{Enabled: false}, nil
		},
	})
	rep, err := daemon.AutostartStatus()
	require.NoError(t, err)
	assert.False(t, rep.Enabled)
}
