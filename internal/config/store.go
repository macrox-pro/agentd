package config

import (
	"context"
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
	RuntimePath string
	ProjectPath string
	Policy      Policy
	Async       AsyncConfig
	Guards      Guards
	Routes      []CompiledRoute
}

// Layer identifies a config source for Get / show.
type Layer string

const (
	LayerUser    Layer = "user"
	LayerProject Layer = "project"
	LayerRuntime Layer = "runtime"
	LayerMerged  Layer = "merged"
)

// Store holds the current config snapshot for lock-free reads on the hot path.
type Store struct {
	snap        atomic.Pointer[Snapshot]
	userPath    string
	runtimePath string
	gen         atomic.Uint64
	reloadMu    sync.Mutex

	userRaw    []byte
	userFC     *fileConfig
	runtimeRaw []byte
	runtimeFC  *fileConfig
	mergedFC   *fileConfig // base merged (no project)

	projects map[string]*projectState // abs .agentd.yaml path → state

	watcher *Watcher
}

type projectState struct {
	path string
	raw  []byte
	fc   *fileConfig
	snap *Snapshot
}

// LoadOptions configures LoadWith.
type LoadOptions struct {
	UserPath    string
	RuntimePath string // empty skips runtime layer
}

// Load reads defaults merged with optional user YAML (no runtime path).
// A missing user file is not an error.
func Load(_ context.Context, userPath string) (*Store, error) {
	return LoadWith(context.Background(), LoadOptions{UserPath: userPath})
}

// LoadWith reads defaults ⊕ user ⊕ runtime into a Store.
func LoadWith(_ context.Context, opts LoadOptions) (*Store, error) {
	s := &Store{
		userPath:    opts.UserPath,
		runtimePath: opts.RuntimePath,
		projects:    map[string]*projectState{},
	}
	if err := s.reloadBase(); err != nil {
		return nil, err
	}
	return s, nil
}

// Current returns the active base snapshot (defaults ⊕ user ⊕ runtime).
func (s *Store) Current() *Snapshot {
	return s.snap.Load()
}

// UserPath returns the configured user config path.
func (s *Store) UserPath() string {
	if s == nil {
		return ""
	}
	return s.userPath
}

// RuntimePath returns the configured runtime overlay path.
func (s *Store) RuntimePath() string {
	if s == nil {
		return ""
	}
	return s.runtimePath
}

// ProjectPaths returns absolute paths of lazily loaded project configs.
func (s *Store) ProjectPaths() []string {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()
	out := make([]string, 0, len(s.projects))
	for p := range s.projects {
		out = append(out, p)
	}
	return out
}

// SnapshotFor returns a project-aware snapshot when cwd/projectRoot resolve to
// a project file; otherwise the base snapshot. Hot-path map lookup after first sighting.
func (s *Store) SnapshotFor(cwd, projectRoot string) *Snapshot {
	if s == nil {
		return nil
	}
	if cwd == "" && projectRoot == "" {
		return s.Current()
	}
	snap, err := s.EnsureProject(cwd, projectRoot)
	if err != nil || snap == nil {
		return s.Current()
	}
	return snap
}

// EnsureProject resolves, loads, and caches a project layer. Missing project is not an error
// (returns base snapshot, nil error).
func (s *Store) EnsureProject(cwd, projectRoot string) (*Snapshot, error) {
	path, ok := FindProjectConfig(cwd, projectRoot)
	if !ok {
		return s.Current(), nil
	}

	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	if ps, hit := s.projects[path]; hit && ps.snap != nil {
		return ps.snap, nil
	}
	ps, err := s.readProjectLocked(path)
	if err != nil {
		return nil, err
	}
	if err := s.compileOneProjectLocked(ps); err != nil {
		return nil, err
	}
	s.projects[path] = ps
	if s.watcher != nil {
		s.watcher.addFile(path)
	}
	return ps.snap, nil
}

// Reload re-reads user and runtime files and recompiles base + known projects.
func (s *Store) Reload(ctx context.Context) error {
	_ = ctx
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()
	return s.reloadAllLocked()
}

// PatchRuntime merges yamlPatch into the in-memory runtime layer and recompiles.
// Does not persist to disk (M7).
func (s *Store) PatchRuntime(yamlPatch []byte) error {
	if len(yamlPatch) == 0 {
		return fmt.Errorf("runtime patch: empty")
	}
	var patch fileConfig
	if err := yaml.Unmarshal(yamlPatch, &patch); err != nil {
		return fmt.Errorf("parse runtime patch: %w", err)
	}

	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	s.runtimeFC = mergeFile(s.runtimeFC, &patch)
	raw, err := yaml.Marshal(s.runtimeFC)
	if err != nil {
		return fmt.Errorf("marshal runtime: %w", err)
	}
	s.runtimeRaw = raw
	return s.recompileAllLocked()
}

// IgnoreSelfWrite marks path so the next watch events for it are skipped (atomic rename).
func (s *Store) IgnoreSelfWrite(path string) {
	if s == nil || s.watcher == nil {
		return
	}
	s.watcher.ignoreSelfWrite(path)
}

// LayerYAML returns YAML bytes for a config layer.
func (s *Store) LayerYAML(layer Layer, cwd, projectRoot string) ([]byte, error) {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	switch layer {
	case LayerUser:
		return append([]byte(nil), s.userRaw...), nil
	case LayerRuntime:
		return append([]byte(nil), s.runtimeRaw...), nil
	case LayerProject:
		path, ok := FindProjectConfig(cwd, projectRoot)
		if !ok {
			return nil, nil
		}
		if ps, hit := s.projects[path]; hit {
			return append([]byte(nil), ps.raw...), nil
		}
		return readOptionalYAML(path)
	case LayerMerged:
		path, ok := FindProjectConfig(cwd, projectRoot)
		var project *fileConfig
		if ok {
			if ps, hit := s.projects[path]; hit {
				project = ps.fc
			} else {
				raw, err := readOptionalYAML(path)
				if err != nil {
					return nil, err
				}
				if len(raw) > 0 {
					var fc fileConfig
					if err := yaml.Unmarshal(raw, &fc); err != nil {
						return nil, fmt.Errorf("parse project config %q: %w", path, err)
					}
					project = &fc
				}
			}
		}
		merged := mergeFile(baseFileConfig(), s.userFC)
		merged = mergeFile(merged, project)
		merged = mergeFile(merged, s.runtimeFC)
		return yaml.Marshal(merged)
	default:
		return nil, fmt.Errorf("unknown layer %q", layer)
	}
}

func (s *Store) reloadBase() error {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()
	return s.reloadAllLocked()
}

func (s *Store) reloadAllLocked() error {
	userRaw, userFC, err := readFileConfig(s.userPath)
	if err != nil {
		return err
	}
	runtimeRaw, runtimeFC, err := readFileConfig(s.runtimePath)
	if err != nil {
		return err
	}
	s.userRaw = userRaw
	s.userFC = userFC
	s.runtimeRaw = runtimeRaw
	s.runtimeFC = runtimeFC

	for path := range s.projects {
		ps, err := s.readProjectLocked(path)
		if err != nil {
			return err
		}
		s.projects[path] = ps
	}
	return s.recompileAllLocked()
}

func (s *Store) recompileAllLocked() error {
	pol, async, guards, routes, merged, err := CompileMerged(s.userFC, nil, s.runtimeFC)
	if err != nil {
		return fmt.Errorf("compile config: %w", err)
	}
	fp, err := Fingerprint(merged)
	if err != nil {
		return err
	}
	s.mergedFC = merged
	gen := s.gen.Add(1)
	base := &Snapshot{
		Generation:  gen,
		Fingerprint: fp,
		UserPath:    s.userPath,
		RuntimePath: s.runtimePath,
		Policy:      pol,
		Async:       async,
		Guards:      guards,
		Routes:      routes,
	}
	s.snap.Store(base)

	for path, ps := range s.projects {
		pol, async, guards, routes, merged, err := CompileMerged(s.userFC, ps.fc, s.runtimeFC)
		if err != nil {
			return fmt.Errorf("compile project config %q: %w", path, err)
		}
		fp, err := Fingerprint(merged)
		if err != nil {
			return err
		}
		ps.snap = &Snapshot{
			Generation:  gen,
			Fingerprint: fp,
			UserPath:    s.userPath,
			RuntimePath: s.runtimePath,
			ProjectPath: path,
			Policy:      pol,
			Async:       async,
			Guards:      guards,
			Routes:      routes,
		}
		s.projects[path] = ps
	}
	return nil
}

func (s *Store) compileOneProjectLocked(ps *projectState) error {
	pol, async, guards, routes, merged, err := CompileMerged(s.userFC, ps.fc, s.runtimeFC)
	if err != nil {
		return fmt.Errorf("compile project config %q: %w", ps.path, err)
	}
	fp, err := Fingerprint(merged)
	if err != nil {
		return err
	}
	gen := s.gen.Add(1)
	ps.snap = &Snapshot{
		Generation:  gen,
		Fingerprint: fp,
		UserPath:    s.userPath,
		RuntimePath: s.runtimePath,
		ProjectPath: ps.path,
		Policy:      pol,
		Async:       async,
		Guards:      guards,
		Routes:      routes,
	}
	return nil
}

func (s *Store) readProjectLocked(path string) (*projectState, error) {
	raw, fc, err := readFileConfig(path)
	if err != nil {
		return nil, err
	}
	return &projectState{path: path, raw: raw, fc: fc}, nil
}

func (s *Store) reloadProjectFile(path string) error {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()
	ps, err := s.readProjectLocked(path)
	if err != nil {
		return err
	}
	if err := s.compileOneProjectLocked(ps); err != nil {
		return err
	}
	s.projects[path] = ps
	return nil
}

func readFileConfig(path string) ([]byte, *fileConfig, error) {
	raw, err := readOptionalYAML(path)
	if err != nil {
		return nil, nil, err
	}
	if len(raw) == 0 {
		return nil, nil, nil
	}
	var fc fileConfig
	if err := yaml.Unmarshal(raw, &fc); err != nil {
		return nil, nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	return raw, &fc, nil
}

func readOptionalYAML(path string) ([]byte, error) {
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
