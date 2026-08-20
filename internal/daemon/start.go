package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/transport"
)

// StartOptions configures daemon Start.
type StartOptions struct {
	Socket     string
	ConfigPath string
	Foreground bool
	Version    string
}

// Start runs the daemon. In foreground mode it blocks until shutdown.
// When Foreground is false it re-execs/detaches and returns after the child is launched.
func Start(ctx context.Context, opts StartOptions) error {
	if opts.Socket == "" {
		opts.Socket = transport.DefaultSocketPath()
	}
	if err := ensureNotRunning(ctx, opts.Socket); err != nil {
		return err
	}
	if !opts.Foreground {
		return detach(opts)
	}
	return runForeground(ctx, opts)
}

func ensureNotRunning(ctx context.Context, socket string) error {
	paths := NewPaths(socket)
	if err := rejectLivePID(paths); err != nil {
		return err
	}

	checkCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	if err := pingHealth(checkCtx, socket); err == nil {
		return ErrAlreadyRunning
	}
	return nil
}

func rejectLivePID(paths Paths) error {
	pid, err := paths.ReadPID()
	if err != nil {
		return nil
	}
	if processAlive(pid) {
		return ErrAlreadyRunning
	}
	paths.RemoveStale()
	return nil
}

func runForeground(ctx context.Context, opts StartOptions) error {
	paths := NewPaths(opts.Socket)
	if err := rejectLivePID(paths); err != nil {
		return err
	}

	lock, err := paths.AcquireLock()
	if err != nil {
		return err
	}
	defer ReleaseLock(lock)

	store, err := config.Load(ctx, opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ln, err := transport.Listen(opts.Socket)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer func() {
		_ = ln.Close()
		paths.RemoveStale()
	}()

	if err := paths.WritePID(os.Getpid()); err != nil {
		return fmt.Errorf("write pid: %w", err)
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	gs := server.New(server.Options{
		Store:      store,
		StartedAt:  time.Now().UTC(),
		Version:    opts.Version,
		OnShutdown: cancel,
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- gs.Serve(ln)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case <-runCtx.Done():
		gs.GracefulStop()
		return nil
	case <-sigCh:
		gs.GracefulStop()
		return nil
	case err := <-errCh:
		if err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	}
}

func pingHealth(ctx context.Context, socket string) error {
	cli, err := hookclient.Dial(ctx, socket)
	if err != nil {
		return err
	}
	defer cli.Close()
	_, err = cli.Health(ctx)
	return err
}
