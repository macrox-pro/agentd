package config

import (
	"context"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const defaultWatchDebounce = 300 * time.Millisecond

// Watcher debounces fsnotify events on the user config file and reloads the store.
type Watcher struct {
	store    *Store
	path     string
	debounce time.Duration
	log      *slog.Logger

	mu     sync.Mutex
	w      *fsnotify.Watcher
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// WatchOptions configures Watch.
type WatchOptions struct {
	Debounce time.Duration
	Log      *slog.Logger
}

// Watch starts watching the store's user config path. A missing path is a no-op.
func (s *Store) Watch(opts WatchOptions) (*Watcher, error) {
	if s == nil || s.userPath == "" {
		return &Watcher{store: s}, nil
	}
	debounce := opts.Debounce
	if debounce <= 0 {
		debounce = defaultWatchDebounce
	}
	log := opts.Log
	if log == nil {
		log = slog.Default()
	}
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(s.userPath)
	if err := fw.Add(dir); err != nil {
		_ = fw.Close()
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	w := &Watcher{
		store:    s,
		path:     s.userPath,
		debounce: debounce,
		log:      log,
		w:        fw,
		cancel:   cancel,
	}
	w.wg.Add(1)
	go w.loop(ctx, fw)
	return w, nil
}

// Close stops the watcher.
func (w *Watcher) Close() error {
	if w == nil {
		return nil
	}
	w.mu.Lock()
	cancel := w.cancel
	fw := w.w
	w.cancel = nil
	w.w = nil
	w.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	var err error
	if fw != nil {
		err = fw.Close()
	}
	w.wg.Wait()
	return err
}

func (w *Watcher) loop(ctx context.Context, fw *fsnotify.Watcher) {
	defer w.wg.Done()
	var timer *time.Timer
	var timerC <-chan time.Time
	base := filepath.Base(w.path)

	for {
		select {
		case <-ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			return
		case ev, ok := <-fw.Events:
			if !ok {
				return
			}
			if filepath.Base(ev.Name) != base {
				continue
			}
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Remove) == 0 {
				continue
			}
			if timer == nil {
				timer = time.NewTimer(w.debounce)
				timerC = timer.C
			} else {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(w.debounce)
			}
		case err, ok := <-fw.Errors:
			if !ok {
				return
			}
			if w.log != nil && err != nil {
				w.log.Warn("config watch error", "error", err)
			}
		case <-timerC:
			timer = nil
			timerC = nil
			if err := w.store.Reload(context.Background()); err != nil && w.log != nil {
				w.log.Warn("config reload failed", "error", err)
			}
		}
	}
}
