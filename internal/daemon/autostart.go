package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/macrox-pro/agentd/internal/transport"
)

// AutostartHooks overrides platform autostart I/O for tests. Nil fields use real OS hooks.
type AutostartHooks struct {
	Register    func(AutostartSpec) error
	Unregister  func() error
	ReadState   func() (AutostartReport, error)
	ResolveExe  func() (string, error)
	StartIfDown func(context.Context, StartOptions) error
}

var (
	autostartTestHooksMu sync.Mutex
	autostartTestHooks   *AutostartHooks
)

// Enable registers login autostart and starts the daemon when it is not running.
func Enable(ctx context.Context, opts AutostartOptions) error {
	exe, err := resolveAutostartExeForEnable()
	if err != nil {
		return err
	}
	if opts.Socket == "" {
		opts.Socket = transport.DefaultSocketPath()
	}
	startOpts := StartOptions{
		Socket:     opts.Socket,
		ConfigPath: opts.ConfigPath,
		Version:    opts.Version,
	}
	spec := AutostartSpec{
		Exe:  exe,
		Args: serviceStartArgs(startOpts),
	}
	if err := callRegisterAutostart(spec); err != nil {
		return err
	}
	if err := callStartIfDown(ctx, startOpts); err != nil {
		return fmt.Errorf("daemon start: %w", err)
	}
	return nil
}

// Disable removes login autostart registration. It does not stop a running daemon.
func Disable() error {
	return callUnregisterAutostart()
}

// AutostartStatus reports whether login autostart is registered on this machine.
func AutostartStatus() (AutostartReport, error) {
	rep, err := callReadAutostartState()
	if err != nil {
		return rep, err
	}
	if !rep.Enabled || rep.RegisteredExe == "" {
		return rep, nil
	}
	cur, err := os.Executable()
	if err != nil {
		return rep, nil
	}
	cur, err = filepath.Abs(cur)
	if err != nil {
		return rep, nil
	}
	rep.Stale = !samePath(rep.RegisteredExe, cur)
	return rep, nil
}

func resolveAutostartExe() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("executable: %w", err)
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return "", fmt.Errorf("abs executable: %w", err)
	}
	return exe, nil
}

func resolveAutostartExeForEnable() (string, error) {
	if h := autostartTestHooksLocked(); h != nil && h.ResolveExe != nil {
		return h.ResolveExe()
	}
	return resolveAutostartExe()
}

func callRegisterAutostart(spec AutostartSpec) error {
	if h := autostartTestHooksLocked(); h != nil && h.Register != nil {
		return h.Register(spec)
	}
	return registerAutostart(spec)
}

func callUnregisterAutostart() error {
	if h := autostartTestHooksLocked(); h != nil && h.Unregister != nil {
		return h.Unregister()
	}
	return unregisterAutostart()
}

func callReadAutostartState() (AutostartReport, error) {
	if h := autostartTestHooksLocked(); h != nil && h.ReadState != nil {
		return h.ReadState()
	}
	return readAutostartState()
}

func callStartIfDown(ctx context.Context, opts StartOptions) error {
	if h := autostartTestHooksLocked(); h != nil && h.StartIfDown != nil {
		return h.StartIfDown(ctx, opts)
	}
	return startIfDown(ctx, opts)
}

func autostartTestHooksLocked() *AutostartHooks {
	autostartTestHooksMu.Lock()
	defer autostartTestHooksMu.Unlock()
	return autostartTestHooks
}

func startIfDown(ctx context.Context, opts StartOptions) error {
	rep, err := Status(ctx, opts.Socket)
	if err != nil {
		return err
	}
	if rep.Running {
		return nil
	}
	err = Start(ctx, opts)
	if errors.Is(err, ErrAlreadyRunning) {
		return nil
	}
	return err
}

func samePath(a, b string) bool {
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	return a == b
}

// SetAutostartHooksForTest overrides platform autostart hooks. For tests only.
func SetAutostartHooksForTest(h *AutostartHooks) func() {
	autostartTestHooksMu.Lock()
	prev := autostartTestHooks
	autostartTestHooks = h
	autostartTestHooksMu.Unlock()
	return func() {
		autostartTestHooksMu.Lock()
		autostartTestHooks = prev
		autostartTestHooksMu.Unlock()
	}
}
