package daemon

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

const importWatchDebounce = 500 * time.Millisecond

// ImportWatcher debounces Claude transcript imports (async_side; never blocks Invoke).
type ImportWatcher struct {
	store  *config.Store
	hub    *trajectory.Hub
	log    *slog.Logger
	mu     sync.Mutex
	cancel context.CancelFunc
}

// NewImportWatcher returns a watcher tied to config reloads.
func NewImportWatcher(store *config.Store, hub *trajectory.Hub, log *slog.Logger) *ImportWatcher {
	return &ImportWatcher{store: store, hub: hub, log: log}
}

// Start begins watching when claude import is enabled in the current snapshot.
func (w *ImportWatcher) Start(ctx context.Context) {
	if w == nil || w.store == nil {
		return
	}
	w.Stop()
	runCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	go w.loop(runCtx)
}

// Stop ends the background watcher.
func (w *ImportWatcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
}

func (w *ImportWatcher) loop(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}
		snap := w.store.Current()
		claude := snap.Trajectory.ClaudeImport()
		if !claude.Enabled {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				continue
			}
		}
		root := claudeProjectsRoot(claude.Path)
		if root == "" {
			w.log.Warn("claude import enabled but projects root unavailable")
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				continue
			}
		}
		w.runWatch(ctx, snap.Trajectory, root)
		return
	}
}

func (w *ImportWatcher) runWatch(ctx context.Context, cfg config.TrajectoryConfig, root string) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		w.log.Warn("import watcher: fsnotify", "error", err)
		return
	}
	defer fw.Close()
	if err := watchTree(fw, root); err != nil {
		w.log.Warn("import watcher: add root", "error", err, "root", root)
		return
	}

	debounce := map[string]*time.Timer{}
	var debounceMu sync.Mutex

	for {
		select {
		case <-ctx.Done():
			debounceMu.Lock()
			for _, t := range debounce {
				t.Stop()
			}
			debounceMu.Unlock()
			return
		case ev, ok := <-fw.Events:
			if !ok {
				return
			}
			if ev.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}
			if !strings.HasSuffix(ev.Name, ".jsonl") {
				continue
			}
			path := ev.Name
			debounceMu.Lock()
			if t, exists := debounce[path]; exists {
				t.Stop()
			}
			debounce[path] = time.AfterFunc(importWatchDebounce, func() {
				w.importFile(cfg, path)
			})
			debounceMu.Unlock()
		case err, ok := <-fw.Errors:
			if !ok {
				return
			}
			w.log.Warn("import watcher", "error", err)
		}
	}
}

func (w *ImportWatcher) importFile(cfg config.TrajectoryConfig, transcriptPath string) {
	sid := strings.TrimSuffix(filepath.Base(transcriptPath), ".jsonl")
	sessionsRoot := trajectory.DefaultSessionsDir()
	if sessionsRoot == "" {
		return
	}
	cp, err := trajectory.LoadImportCheckpoint(trajectory.ImportSidecarPath(sessionsRoot, "claude-code", sid))
	if err != nil {
		w.log.Warn("import watcher: checkpoint", "error", err, "session", sid)
		return
	}
	startIndex := 0
	if cp.SourcePath != "" {
		startIndex = cp.LastLineIndex + 1
	}
	result, err := importer.ImportClaude(importer.ImportOptions{
		SessionID:      sid,
		TranscriptPath: transcriptPath,
		ProjectsRoot:   cfg.ClaudeImport().Path,
		StartIndex:     startIndex,
		Cfg:            cfg,
	})
	if err != nil {
		w.log.Warn("import watcher: import", "error", err, "session", sid)
		return
	}
	if len(result.Events) == 0 {
		return
	}
	key := trajectory.ResolveSessionKey("claude-code", sid, "", "")
	if err := trajectory.AppendImported(sessionsRoot, key, result.Events); err != nil {
		w.log.Warn("import watcher: append", "error", err, "session", sid)
		return
	}
	if w.hub != nil && len(result.Events) > 0 {
		w.hub.Publish(result.Events)
	}
	st, err := os.Stat(transcriptPath)
	if err != nil {
		w.log.Warn("import watcher: stat", "error", err)
		return
	}
	cp.LastLineIndex = result.LastLineIndex
	cp.SourcePath = transcriptPath
	cp.SourceModTime = st.ModTime().UTC()
	if err := trajectory.SaveImportCheckpoint(trajectory.ImportSidecarPath(sessionsRoot, "claude-code", sid), cp); err != nil {
		w.log.Warn("import watcher: save checkpoint", "error", err)
	}
}

func claudeProjectsRoot(configured string) string {
	return ClaudeProjectsRootForTest(configured)
}

// ClaudeProjectsRootForTest resolves Claude projects directory (exported for tests).
func ClaudeProjectsRootForTest(configured string) string {
	if configured != "" {
		return configured
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "projects")
}

func watchTree(w *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if err := w.Add(path); err != nil {
			return err
		}
		return nil
	})
}
