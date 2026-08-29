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
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/transport"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

const (
	readyTimeout = 5 * time.Second
	readyPoll    = 25 * time.Millisecond
)

// StartOptions configures daemon Start.
type StartOptions struct {
	Socket     string
	ConfigPath string
	Foreground bool
	Version    string
	LogLevel   string
	LogFile    string
}

// Start runs the daemon. In foreground mode it blocks until shutdown.
// When Foreground is false it re-execs/detaches and returns only after the
// child answers Health (or a readiness timeout).
func Start(ctx context.Context, opts StartOptions) error {
	if opts.Socket == "" {
		opts.Socket = transport.DefaultSocketPath()
	}
	if err := ensureNotRunning(ctx, opts.Socket); err != nil {
		return err
	}
	userPath := opts.ConfigPath
	if userPath == "" {
		userPath = config.DefaultUserPath()
	}
	if userPath == "" {
		return fmt.Errorf("user config path: home directory unavailable")
	}
	if err := config.PrepareUserConfig(userPath, os.Stderr); err != nil {
		return err
	}
	opts.ConfigPath = userPath
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

// rejectLivePID returns ErrAlreadyRunning when a live PID file is present.
// It never removes socket/PID — cleanup happens only under the lock.
func rejectLivePID(paths Paths) error {
	pid, err := paths.ReadPID()
	if err != nil {
		return nil
	}
	if processAlive(pid) {
		return ErrAlreadyRunning
	}
	return nil
}

func cleanStaleUnderLock(paths Paths) {
	paths.RemoveStale()
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

	cleanStaleUnderLock(paths)

	store, err := config.LoadWith(ctx, config.LoadOptions{
		UserPath:    opts.ConfigPath,
		RuntimePath: config.DefaultRuntimePath(),
	})
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

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	snap := store.Current()
	log, logCleanup, err := SetupLog(SetupLogOptions{
		Logging:    snap.Logging,
		Foreground: opts.Foreground,
		LogLevel:   opts.LogLevel,
		LogFile:    opts.LogFile,
	})
	if err != nil {
		return fmt.Errorf("setup log: %w", err)
	}
	defer func() {
		log.Info("daemon shutdown complete")
		logCleanup()
	}()

	store.SetLogger(log)
	queue := dispatch.NewQueue(snap.Async, log)
	engine := dispatch.NewEngine(queue, log)
	defer queue.Close(5 * time.Second)

	trajCfg := snap.Trajectory
	recorder := trajectory.NewRecorder(trajectory.DefaultSessionsDir(), trajCfg.QueueCapacity, log)
	defer recorder.Close(5 * time.Second)

	watcher, err := store.Watch(config.WatchOptions{Log: log})
	if err != nil {
		return fmt.Errorf("watch config: %w", err)
	}
	defer func() { _ = watcher.Close() }()

	importWatcher := NewImportWatcher(store, recorder.Hub(), log)
	importWatcher.Start(runCtx)
	defer importWatcher.Stop()

	gs := server.New(server.Options{
		Store:      store,
		Engine:     engine,
		Recorder:   recorder,
		Logger:     log,
		StartedAt:  time.Now().UTC(),
		Version:    opts.Version,
		OnShutdown: cancel,
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- gs.Serve(ln)
	}()

	readyCtx, readyCancel := context.WithTimeout(ctx, readyTimeout)
	err = waitHealth(readyCtx, opts.Socket)
	readyCancel()
	if err != nil {
		gs.GracefulStop()
		return fmt.Errorf("daemon failed to become ready: %w", err)
	}

	if err := paths.WritePID(os.Getpid()); err != nil {
		gs.GracefulStop()
		return fmt.Errorf("write pid: %w", err)
	}

	cur := store.Current()
	log.Info("daemon ready",
		"version", opts.Version,
		"socket", opts.Socket,
		"generation", cur.Generation,
		"fingerprint", cur.Fingerprint,
	)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	reloadCh := make(chan os.Signal, 1)
	notifyReload(reloadCh)
	defer signal.Stop(reloadCh)

	shutdown := func() {
		log.Info("daemon shutdown started")
		// Drop PID before signal.Stop / queue drain so Stop does not fall through
		// to SIGTERM against this process (foreground PID == os.Getpid()).
		paths.ClearPID()
		_ = store.FlushRuntime()
		gs.GracefulStop()
	}

	for {
		select {
		case <-runCtx.Done():
			shutdown()
			return nil
		case sig := <-sigCh:
			if isReloadSignal(sig) {
				if err := store.Reload(runCtx); err != nil {
					log.Warn("config reload failed", "error", err)
				} else {
					cur := store.Current()
					log.Info("config reload succeeded",
						"generation", cur.Generation,
						"fingerprint", cur.Fingerprint,
					)
					importWatcher.Start(runCtx)
				}
				continue
			}
			shutdown()
			return nil
		case <-reloadCh:
			if err := store.Reload(context.Background()); err != nil {
				log.Warn("config reload failed", "error", err)
			} else {
				cur := store.Current()
				log.Info("config reload succeeded",
					"generation", cur.Generation,
					"fingerprint", cur.Fingerprint,
				)
				importWatcher.Start(runCtx)
			}
		case err := <-errCh:
			paths.ClearPID()
			_ = store.FlushRuntime()
			if err != nil && err != grpc.ErrServerStopped {
				return err
			}
			return nil
		}
	}
}

func waitHealth(ctx context.Context, socket string) error {
	for {
		checkCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		err := pingHealth(checkCtx, socket)
		cancel()
		if err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return fmt.Errorf("health: %w", ctx.Err())
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("health: %w", ctx.Err())
		case <-time.After(readyPoll):
		}
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
