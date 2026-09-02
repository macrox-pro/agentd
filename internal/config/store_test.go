package config_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestStoreLoad(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		content string
		write   bool
		wantErr bool
		check   func(t *testing.T, store *config.Store)
	}{
		{
			name:  "missing file ok",
			write: false,
			check: func(t *testing.T, store *config.Store) {
				t.Helper()
				snap := store.Current()
				assert.Equal(t, uint64(1), snap.Generation, "Load(missing)")
				assert.NotEmpty(t, snap.Fingerprint, "Load(missing)")
				assert.True(t, snap.Guards.Secrets.Enabled, "Load(missing) secrets")
				assert.NotEmpty(t, snap.Routes, "Load(missing) routes")
				assert.Equal(t, 1024, snap.Async.QueueCapacity, "Load(missing) async")
			},
		},
		{
			name:    "invalid yaml",
			write:   true,
			content: ":\tinvalid",
			wantErr: true,
		},
		{
			name:    "valid yaml",
			write:   true,
			content: "version: 1\n",
			check: func(t *testing.T, store *config.Store) {
				t.Helper()
				snap := store.Current()
				assert.GreaterOrEqual(t, snap.Generation, uint64(1), "Load(valid)")
				assert.NotEmpty(t, snap.Fingerprint, "Load(valid)")
			},
		},
		{
			name:  "override policy and secrets",
			write: true,
			content: `version: 1
policy:
  fail: fail_open
guards:
  secrets:
    enabled: false
    action: deny
async:
  queue_capacity: 16
  worker_limit: 2
  target_timeout: 5s
`,
			check: func(t *testing.T, store *config.Store) {
				t.Helper()
				snap := store.Current()
				assert.Equal(t, config.FailOpen, snap.Policy.Fail, "policy.fail")
				assert.False(t, snap.Guards.Secrets.Enabled, "secrets.enabled")
				assert.Equal(t, config.GuardDeny, snap.Guards.Secrets.Action, "secrets.action")
				assert.Equal(t, 16, snap.Async.QueueCapacity, "async.queue_capacity")
				assert.Equal(t, 2, snap.Async.WorkerLimit, "async.worker_limit")
				assert.Equal(t, 5*time.Second, snap.Async.TargetTimeout, "async.target_timeout")
			},
		},
		{
			name:    "bad policy fail",
			write:   true,
			content: "version: 1\npolicy:\n  fail: nope\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "agentd.yaml")
			if tt.write {
				require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o600), "WriteFile(%q)", path)
			}

			store, err := config.Load(ctx, path)
			if tt.wantErr {
				require.Error(t, err, "Load(%q)", path)
				return
			}
			require.NoError(t, err, "Load(%q)", path)
			if tt.check != nil {
				tt.check(t, store)
			}
		})
	}
}

func TestReload(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600), "WriteFile(%q)", path)

	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "policy edit bumps generation",
			content: "version: 1\npolicy:\n  fail: fail_closed\n",
		},
		{
			name:    "second reload bumps again",
			content: "version: 1\npolicy:\n  fail: fail_open\n",
		},
	}

	var prevGen uint64
	var prevFP string
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o600), "WriteFile(%q)", path)
			require.NoError(t, store.Reload(ctx), "Reload(%q)", path)

			snap := store.Current()
			if prevGen > 0 {
				assert.Greater(t, snap.Generation, prevGen, "Reload(%q)", tt.name)
				assert.NotEqual(t, prevFP, snap.Fingerprint, "Reload(%q)", tt.name)
			}
			prevGen = snap.Generation
			prevFP = snap.Fingerprint
		})
	}
}

func TestReloadConcurrent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600), "WriteFile(%q)", path)

	store, err := config.Load(ctx, path)
	require.NoError(t, err, "Load(%q)", path)
	before := store.Current().Generation

	const workers = 8
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			errs <- store.Reload(ctx)
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err, "concurrent Reload")
	}

	snap := store.Current()
	require.NotNil(t, snap)
	assert.Greater(t, snap.Generation, before, "concurrent Reload")
	assert.NotEmpty(t, snap.Fingerprint)
	assert.NotEmpty(t, snap.Routes)
	assert.Equal(t, config.FailOpen, snap.Policy.Fail)
}

func TestEnsureProjectOverridesUser(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	projDir := filepath.Join(dir, "proj")
	require.NoError(t, os.MkdirAll(projDir, 0o700))
	projPath := filepath.Join(projDir, ".agentd.yaml")
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))
	require.NoError(t, os.WriteFile(projPath, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	store, err := config.Load(context.Background(), userPath)
	require.NoError(t, err)
	assert.Equal(t, config.FailClosed, store.Current().Policy.Fail)

	snap, err := store.EnsureProject(projDir, "")
	require.NoError(t, err)
	require.NotNil(t, snap)
	assert.Equal(t, config.FailOpen, snap.Policy.Fail, "project should override user")
	assert.Equal(t, config.FailClosed, store.Current().Policy.Fail, "base unchanged")
	assert.Equal(t, projPath, snap.ProjectPath)
}

func TestLoadWithRuntimeOverrides(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	runtimePath := filepath.Join(dir, "runtime.yaml")
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))
	require.NoError(t, os.WriteFile(runtimePath, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	store, err := config.LoadWith(context.Background(), config.LoadOptions{
		UserPath:    userPath,
		RuntimePath: runtimePath,
	})
	require.NoError(t, err)
	assert.Equal(t, config.FailOpen, store.Current().Policy.Fail)
	assert.Equal(t, runtimePath, store.RuntimePath())
}

func TestPatchRuntime(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))

	store, err := config.Load(context.Background(), userPath)
	require.NoError(t, err)
	before := store.Current().Generation

	require.NoError(t, store.PatchRuntime([]byte("version: 1\npolicy:\n  fail: fail_open\n")))
	snap := store.Current()
	assert.Equal(t, config.FailOpen, snap.Policy.Fail)
	assert.Greater(t, snap.Generation, before)
}

func TestSnapshotFor(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	projDir := filepath.Join(dir, "repo")
	require.NoError(t, os.MkdirAll(projDir, 0o700))
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(projDir, ".agentd.yaml"), []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	store, err := config.Load(context.Background(), userPath)
	require.NoError(t, err)

	base := store.SnapshotFor("", "")
	assert.Equal(t, config.FailClosed, base.Policy.Fail)

	proj := store.SnapshotFor(projDir, "")
	assert.Equal(t, config.FailOpen, proj.Policy.Fail)

	again := store.SnapshotFor(projDir, "")
	assert.Equal(t, proj.Generation, again.Generation)
}

func TestLayerYAML_project(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	projDir := filepath.Join(dir, "proj")
	require.NoError(t, os.MkdirAll(projDir, 0o700))
	projPath := filepath.Join(projDir, ".agentd.yaml")
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))
	require.NoError(t, os.WriteFile(projPath, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	store, err := config.Load(context.Background(), userPath)
	require.NoError(t, err, "Load(%q)", userPath)

	raw, err := store.LayerYAML(config.LayerProject, projDir, "")
	require.NoError(t, err, "LayerYAML(project)")
	assert.Contains(t, string(raw), "fail_open")
	assert.NotContains(t, string(raw), "null:")
}

func TestLayerYAML_merged(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	projDir := filepath.Join(dir, "proj")
	require.NoError(t, os.MkdirAll(projDir, 0o700))
	projPath := filepath.Join(projDir, ".agentd.yaml")
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))
	require.NoError(t, os.WriteFile(projPath, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	store, err := config.Load(context.Background(), userPath)
	require.NoError(t, err, "Load(%q)", userPath)

	raw, err := store.LayerYAML(config.LayerMerged, projDir, "")
	require.NoError(t, err, "LayerYAML(merged)")
	assert.Contains(t, string(raw), "fail_open", "merged_order: project overrides user")
}

func TestLayerYAML_project_cached(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	projDir := filepath.Join(dir, "proj")
	require.NoError(t, os.MkdirAll(projDir, 0o700))
	projPath := filepath.Join(projDir, ".agentd.yaml")
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\n"), 0o600))
	require.NoError(t, os.WriteFile(projPath, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	store, err := config.Load(context.Background(), userPath)
	require.NoError(t, err, "Load(%q)", userPath)

	uncached, err := store.LayerYAML(config.LayerProject, projDir, "")
	require.NoError(t, err, "LayerYAML(project)")

	_, err = store.EnsureProject(projDir, "")
	require.NoError(t, err, "EnsureProject(%q)", projDir)

	cached, err := store.LayerYAML(config.LayerProject, projDir, "")
	require.NoError(t, err, "LayerYAML(project)")
	assert.Equal(t, string(uncached), string(cached), "project_cached")
}

func TestStoreProjectCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "SnapshotFor and Reload remain race free",
			run: func(t *testing.T) {
				t.Helper()
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				projDir := filepath.Join(dir, "repo")
				require.NoError(t, os.MkdirAll(projDir, 0o700), "MkdirAll(repo)")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\n"), 0o600), "WriteFile(user)")
				require.NoError(t, os.WriteFile(filepath.Join(projDir, ".agentd.yaml"), []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600), "WriteFile(project)")

				store, err := config.Load(ctx, userPath)
				require.NoError(t, err, "Load(user)")

				const workers = 16
				var wg sync.WaitGroup
				errs := make(chan error, workers)
				wg.Add(workers)
				for range workers {
					go func() {
						defer wg.Done()
						for range 20 {
							if snap := store.SnapshotFor(projDir, ""); snap == nil {
								errs <- fmt.Errorf("nil snapshot")
								return
							}
							if err := store.Reload(ctx); err != nil {
								errs <- err
								return
							}
						}
					}()
				}
				wg.Wait()
				close(errs)
				for err := range errs {
					require.NoError(t, err, "concurrent SnapshotFor/Reload")
				}

				snap := store.SnapshotFor(projDir, "")
				require.NotNil(t, snap, "SnapshotFor after race")
				assert.Equal(t, config.FailOpen, snap.Policy.Fail, "SnapshotFor after race")
			},
		},
		{
			name: "concurrent first load returns one valid project snapshot",
			run: func(t *testing.T) {
				t.Helper()
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				projDir := filepath.Join(dir, "repo")
				require.NoError(t, os.MkdirAll(projDir, 0o700), "MkdirAll(repo)")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\n"), 0o600), "WriteFile(user)")
				require.NoError(t, os.WriteFile(filepath.Join(projDir, ".agentd.yaml"), []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600), "WriteFile(project)")

				store, err := config.Load(ctx, userPath)
				require.NoError(t, err, "Load(user)")

				const workers = 16
				var wg sync.WaitGroup
				errs := make(chan error, workers)
				snaps := make([]*config.Snapshot, workers)
				wg.Add(workers)
				for i := range workers {
					go func(idx int) {
						defer wg.Done()
						snap, loadErr := store.EnsureProject(projDir, "")
						if loadErr != nil {
							errs <- loadErr
							return
						}
						snaps[idx] = snap
					}(i)
				}
				wg.Wait()
				close(errs)
				for err := range errs {
					require.NoError(t, err, "concurrent EnsureProject")
				}

				var first *config.Snapshot
				for i, snap := range snaps {
					require.NotNil(t, snap, "EnsureProject worker %d", i)
					assert.Equal(t, config.FailOpen, snap.Policy.Fail, "EnsureProject worker %d", i)
					if first == nil {
						first = snap
					} else {
						assert.Equal(t, first.Generation, snap.Generation, "EnsureProject worker %d", i)
						assert.Equal(t, first.Fingerprint, snap.Fingerprint, "EnsureProject worker %d", i)
					}
				}
			},
		},
		{
			name: "reload replaces cached project snapshot consistently",
			run: func(t *testing.T) {
				t.Helper()
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				projDir := filepath.Join(dir, "repo")
				projPath := filepath.Join(projDir, ".agentd.yaml")
				require.NoError(t, os.MkdirAll(projDir, 0o700), "MkdirAll(repo)")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\n"), 0o600), "WriteFile(user)")
				require.NoError(t, os.WriteFile(projPath, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600), "WriteFile(project)")

				store, err := config.Load(ctx, userPath)
				require.NoError(t, err, "Load(user)")

				before, err := store.EnsureProject(projDir, "")
				require.NoError(t, err, "EnsureProject(before)")
				require.Equal(t, config.FailOpen, before.Policy.Fail, "EnsureProject(before)")

				require.NoError(t, os.WriteFile(projPath, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600), "WriteFile(project reload)")
				require.NoError(t, store.Reload(ctx), "Reload(project)")

				after := store.SnapshotFor(projDir, "")
				require.NotNil(t, after, "SnapshotFor(after)")
				assert.Equal(t, config.FailClosed, after.Policy.Fail, "SnapshotFor(after)")
				assert.Greater(t, after.Generation, before.Generation, "generation after reload")
			},
		},
		{
			name: "ProjectPaths is safe during concurrent project discovery",
			run: func(t *testing.T) {
				t.Helper()
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\n"), 0o600), "WriteFile(user)")

				store, err := config.Load(ctx, userPath)
				require.NoError(t, err, "Load(user)")

				const projects = 8
				projDirs := make([]string, projects)
				for i := range projects {
					projDir := filepath.Join(dir, fmt.Sprintf("repo%d", i))
					require.NoError(t, os.MkdirAll(projDir, 0o700), "MkdirAll(%d)", i)
					require.NoError(t, os.WriteFile(filepath.Join(projDir, ".agentd.yaml"), []byte("version: 1\n"), 0o600), "WriteFile(%d)", i)
					projDirs[i] = projDir
				}

				const workers = 16
				var wg sync.WaitGroup
				errs := make(chan error, workers)
				wg.Add(workers)
				for w := range workers {
					go func(worker int) {
						defer wg.Done()
						projDir := projDirs[worker%projects]
						if _, loadErr := store.EnsureProject(projDir, ""); loadErr != nil {
							errs <- loadErr
							return
						}
						paths := store.ProjectPaths()
						if len(paths) == 0 {
							errs <- fmt.Errorf("empty ProjectPaths")
						}
					}(w)
				}
				wg.Wait()
				close(errs)
				for err := range errs {
					require.NoError(t, err, "concurrent ProjectPaths")
				}

				assert.GreaterOrEqual(t, len(store.ProjectPaths()), 1, "ProjectPaths after discovery")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.run(t)
		})
	}
}
