package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	persistDebounce  = 500 * time.Millisecond
	runtimeTmpSuffix = ".tmp"
	runtimeFilePerm  = 0o600
	runtimeDirPerm   = 0o700
)

// schedulePersistLocked queues a debounced flush of runtime.yaml.
// Caller must hold reloadMu.
func (s *Store) schedulePersistLocked() {
	if s == nil || s.runtimePath == "" {
		return
	}
	if s.persistTimer != nil {
		s.persistTimer.Stop()
	}
	s.persistTimer = time.AfterFunc(persistDebounce, s.flushRuntimeAsync)
}

func (s *Store) flushRuntimeAsync() {
	if err := s.FlushRuntime(); err != nil {
		slog.Default().Warn("runtime persist failed", "error", err)
	}
}

// FlushRuntime writes the in-memory runtime overlay to disk immediately.
func (s *Store) FlushRuntime() error {
	if s == nil {
		return nil
	}
	s.reloadMu.Lock()
	if s.persistTimer != nil {
		s.persistTimer.Stop()
		s.persistTimer = nil
	}
	path := s.runtimePath
	raw := append([]byte(nil), s.runtimeRaw...)
	s.reloadMu.Unlock()

	if path == "" {
		return nil
	}
	if len(raw) == 0 {
		// Still persist empty document so reloads see cleared overlay.
		raw = []byte("version: 1\n")
	}
	return writeRuntimeAtomic(s, path, raw)
}

func writeRuntimeAtomic(s *Store, path string, raw []byte) error {
	if s != nil {
		s.IgnoreSelfWrite(path)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, runtimeDirPerm); err != nil {
		return fmt.Errorf("mkdir runtime dir: %w", err)
	}
	tmp := path + runtimeTmpSuffix
	if err := os.WriteFile(tmp, raw, runtimeFilePerm); err != nil {
		return fmt.Errorf("write runtime tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename runtime: %w", err)
	}
	return nil
}
