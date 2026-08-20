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

type watchKind int

const (
	watchUser watchKind = iota
	watchRuntime
	watchProject
)

// Watcher debounces fsnotify events on user, runtime, and project config files.
type Watcher struct {
	store    *Store
	debounce time.Duration
	log      *slog.Logger

	mu      sync.Mutex
	w       *fsnotify.Watcher
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	files   map[string]watchKind // abs path → kind
	dirs    map[string]int       // dir → refcount
	pending map[string]struct{}
	ignore  map[string]time.Time // self-write ignore until
}

// WatchOptions configures Watch.
type WatchOptions struct {
	Debounce time.Duration
	Log      *slog.Logger
}

// Watch starts watching the store's user and runtime config paths.
func (s *Store) Watch(opts WatchOptions) (*Watcher, error) {
	if s == nil {
		return &Watcher{}, nil
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
	ctx, cancel := context.WithCancel(context.Background())
	w := &Watcher{
		store:    s,
		debounce: debounce,
		log:      log,
		w:        fw,
		cancel:   cancel,
		files:    map[string]watchKind{},
		dirs:     map[string]int{},
		pending:  map[string]struct{}{},
		ignore:   map[string]time.Time{},
	}
	s.reloadMu.Lock()
	s.watcher = w
	s.reloadMu.Unlock()

	w.addFile(s.userPath)
	if s.runtimePath != "" {
		w.addFile(s.runtimePath)
	}
	for _, p := range s.ProjectPaths() {
		w.addFile(p)
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
	if w.store != nil {
		w.store.reloadMu.Lock()
		if w.store.watcher == w {
			w.store.watcher = nil
		}
		w.store.reloadMu.Unlock()
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

func (w *Watcher) addFile(path string) {
	if w == nil || path == "" {
		return
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	kind := watchProject
	if w.store != nil {
		if abs == absPath(w.store.userPath) {
			kind = watchUser
		} else if abs == absPath(w.store.runtimePath) {
			kind = watchRuntime
		}
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.files[abs]; ok {
		return
	}
	fw := w.w
	if fw == nil {
		return
	}
	dir := filepath.Dir(abs)
	if w.dirs[dir] == 0 {
		if err := fw.Add(dir); err != nil {
			if w.log != nil {
				w.log.Warn("config watch add", "path", dir, "error", err)
			}
			return
		}
	}
	w.dirs[dir]++
	w.files[abs] = kind
}

func (w *Watcher) ignoreSelfWrite(path string) {
	if w == nil || path == "" {
		return
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	// Cover debounce window plus rename settle.
	w.ignore[abs] = time.Now().Add(w.debounce + time.Second)
}

func absPath(p string) string {
	if p == "" {
		return ""
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

func (w *Watcher) loop(ctx context.Context, fw *fsnotify.Watcher) {
	defer w.wg.Done()
	var timer *time.Timer
	var timerC <-chan time.Time

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
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Remove) == 0 {
				continue
			}
			abs, err := filepath.Abs(ev.Name)
			if err != nil {
				continue
			}
			w.mu.Lock()
			kind, tracked := w.files[abs]
			if !tracked {
				w.mu.Unlock()
				continue
			}
			if until, skip := w.ignore[abs]; skip {
				if time.Now().Before(until) {
					w.mu.Unlock()
					continue
				}
				delete(w.ignore, abs)
			}
			_ = kind
			w.pending[abs] = struct{}{}
			w.mu.Unlock()

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
			w.flushPending()
		}
	}
}

func (w *Watcher) flushPending() {
	w.mu.Lock()
	pending := w.pending
	w.pending = map[string]struct{}{}
	files := make(map[string]watchKind, len(pending))
	for p := range pending {
		if k, ok := w.files[p]; ok {
			files[p] = k
		}
	}
	w.mu.Unlock()

	needBase := false
	var projects []string
	for p, kind := range files {
		switch kind {
		case watchUser, watchRuntime:
			needBase = true
		case watchProject:
			projects = append(projects, p)
		}
	}
	if needBase {
		if err := w.store.Reload(context.Background()); err != nil && w.log != nil {
			w.log.Warn("config reload failed", "error", err)
		}
		return
	}
	for _, p := range projects {
		if err := w.store.reloadProjectFile(p); err != nil && w.log != nil {
			w.log.Warn("project config reload failed", "path", p, "error", err)
		}
	}
}
