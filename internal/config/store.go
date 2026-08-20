package config

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"go.yaml.in/yaml/v3"
)

// Snapshot is an immutable compiled configuration generation.
type Snapshot struct {
	Generation  uint64
	Fingerprint string
	UserPath    string
	Policy      Policy
	Async       AsyncConfig
	Guards      Guards
	Routes      []CompiledRoute
}

// Store holds the current config snapshot for lock-free reads on the hot path.
type Store struct {
	snap     atomic.Pointer[Snapshot]
	userPath string
	gen      atomic.Uint64
	reloadMu sync.Mutex
}

// Load reads defaults merged with an optional user YAML file.
// A missing user file is not an error.
func Load(_ context.Context, userPath string) (*Store, error) {
	s := &Store{userPath: userPath}
	if err := s.reload(); err != nil {
		return nil, err
	}
	return s, nil
}

// Current returns the active snapshot.
func (s *Store) Current() *Snapshot {
	return s.snap.Load()
}

// Reload re-reads the user config file and swaps the snapshot.
func (s *Store) Reload(ctx context.Context) error {
	_ = ctx
	return s.reload()
}

func (s *Store) reload() error {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	raw, err := readUserYAML(s.userPath)
	if err != nil {
		return err
	}
	var user *fileConfig
	if len(raw) > 0 {
		var fc fileConfig
		if err := yaml.Unmarshal(raw, &fc); err != nil {
			return fmt.Errorf("parse config %q: %w", s.userPath, err)
		}
		user = &fc
	}
	pol, async, guards, routes, err := Compile(user)
	if err != nil {
		return fmt.Errorf("compile config %q: %w", s.userPath, err)
	}
	// Fingerprint is sha256 of raw user YAML; DESIGN canonical merged JSON is deferred.
	sum := sha256.Sum256(raw)
	fp := hex.EncodeToString(sum[:])
	gen := s.gen.Add(1)
	snap := &Snapshot{
		Generation:  gen,
		Fingerprint: fp,
		UserPath:    s.userPath,
		Policy:      pol,
		Async:       async,
		Guards:      guards,
		Routes:      routes,
	}
	s.snap.Store(snap)
	return nil
}

func readUserYAML(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	return b, nil
}
