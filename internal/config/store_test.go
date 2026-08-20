package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestLoad(t *testing.T) {
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

func TestCompile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		wantMode map[string]config.DispatchMode
		wantSync bool
	}{
		{
			name: "defaults",
			wantMode: map[string]config.DispatchMode{
				"tool.pre":     config.ModeParallel,
				"notification": config.ModeAsyncOnly,
				"agent.stop":   config.ModeAfterSync, // normalized from sync_then_async
			},
			wantSync: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pol, async, guards, routes, err := config.Compile(nil)
			require.NoError(t, err, "Compile(nil)")
			assert.Equal(t, config.FailClosed, pol.Fail, "Compile policy")
			assert.Equal(t, 1024, async.QueueCapacity, "Compile async")
			assert.True(t, guards.Secrets.Enabled, "Compile secrets")
			byKind := map[string]config.CompiledRoute{}
			for _, r := range routes {
				byKind[r.Kind] = r
			}
			for kind, mode := range tt.wantMode {
				r, ok := byKind[kind]
				require.True(t, ok, "route %q", kind)
				assert.Equal(t, mode, r.Mode, "route %q mode", kind)
				if tt.wantSync && mode != config.ModeAsyncOnly {
					assert.NotEmpty(t, r.Sync, "route %q sync", kind)
				}
			}
		})
	}
}

func TestLoadDispatchRoutes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name    string
		content string
		wantErr bool
		check   func(t *testing.T, snap *config.Snapshot)
	}{
		{
			name: "named route with file async",
			content: `version: 1
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: file
        path: /tmp/audit.jsonl
      - target: log
        level: info
`,
			check: func(t *testing.T, snap *config.Snapshot) {
				t.Helper()
				require.GreaterOrEqual(t, len(snap.Routes), 1)
				r := snap.Routes[0]
				assert.Equal(t, "gate-and-audit", r.Name)
				assert.False(t, r.Default)
				assert.Equal(t, config.ModeParallel, r.Mode)
				require.Len(t, r.Async, 2)
				assert.Equal(t, config.TargetFile, r.Async[0].Kind)
				assert.Equal(t, config.TargetLog, r.Async[1].Kind)
			},
		},
		{
			name: "reject grpc",
			content: `version: 1
dispatch:
  - name: fwd
    match:
      kind: [tool.pre]
    mode: async_only
    async:
      - target: grpc
`,
			wantErr: true,
		},
		{
			name: "reject sync http",
			content: `version: 1
dispatch:
  - name: bad
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: http
        url: http://example.com
`,
			wantErr: true,
		},
		{
			name: "reject unknown guard",
			content: `version: 1
dispatch:
  - name: bad
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [shell]
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "agentd.yaml")
			require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o600))
			store, err := config.Load(ctx, path)
			if tt.wantErr {
				require.Error(t, err, "Load(%q)", tt.name)
				return
			}
			require.NoError(t, err, "Load(%q)", tt.name)
			if tt.check != nil {
				tt.check(t, store.Current())
			}
		})
	}
}

func TestNormalizeMode(t *testing.T) {
	t.Parallel()
	assert.Equal(t, config.ModeAfterSync, config.NormalizeMode(config.ModeSyncThenAsync))
	assert.Equal(t, config.ModeParallel, config.NormalizeMode(config.ModeParallel))
}
