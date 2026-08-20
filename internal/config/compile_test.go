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
			name: "route sync_timeout",
			content: `version: 1
dispatch:
  - name: tight
    match:
      kind: [tool.pre]
    mode: sync_only
    sync_timeout: 2s
    sync:
      - target: builtin
        guards: [secrets]
`,
			check: func(t *testing.T, snap *config.Snapshot) {
				t.Helper()
				require.GreaterOrEqual(t, len(snap.Routes), 1)
				assert.Equal(t, 2*time.Second, snap.Routes[0].SyncTimeout, "sync_timeout")
			},
		},
		{
			name: "accept grpc async",
			content: `version: 1
dispatch:
  - name: fwd
    match:
      kind: [tool.pre]
    mode: async_only
    async:
      - target: grpc
        endpoint: unix:///tmp/peer.sock
        timeout: 2s
`,
			check: func(t *testing.T, snap *config.Snapshot) {
				t.Helper()
				require.GreaterOrEqual(t, len(snap.Routes), 1)
				r := snap.Routes[0]
				require.Len(t, r.Async, 1)
				assert.Equal(t, config.TargetGRPC, r.Async[0].Kind)
				assert.Equal(t, "unix:///tmp/peer.sock", r.Async[0].Endpoint)
				assert.Equal(t, 2*time.Second, r.Async[0].Timeout)
				assert.Equal(t, config.FailClosed, r.Async[0].OnError)
			},
		},
		{
			name: "accept grpc sync",
			content: `version: 1
dispatch:
  - name: gate
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: grpc
        endpoint: /tmp/peer.sock
        on_error: fail_open
        merge: first_conclusive
`,
			check: func(t *testing.T, snap *config.Snapshot) {
				t.Helper()
				r := snap.Routes[0]
				require.Len(t, r.Sync, 1)
				assert.Equal(t, config.TargetGRPC, r.Sync[0].Kind)
				assert.Equal(t, config.FailOpen, r.Sync[0].OnError)
				assert.Equal(t, config.MergeFirstConclusive, r.Sync[0].Merge)
			},
		},
		{
			name: "reject grpc without endpoint",
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
			name: "reject grpc bad on_error",
			content: `version: 1
dispatch:
  - name: fwd
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: grpc
        endpoint: /tmp/peer.sock
        on_error: nope
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
        guards: [nope]
`,
			wantErr: true,
		},
		{
			name: "accept shell mcp paths guards",
			content: `version: 1
guards:
  shell:
    enabled: true
    deny_patterns: ["rm -rf /"]
    ask_on: [curl]
  mcp:
    enabled: true
    deny_servers: ["untrusted-*"]
  paths:
    enabled: true
    deny_read: ["/etc/shadow"]
    deny_write: ["**/.env"]
dispatch:
  - name: gate
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [secrets, shell, mcp, paths]
`,
			check: func(t *testing.T, snap *config.Snapshot) {
				t.Helper()
				assert.True(t, snap.Guards.Shell.Enabled, "shell.enabled")
				assert.Equal(t, []string{"rm -rf /"}, snap.Guards.Shell.DenyPatterns)
				assert.Equal(t, []string{"curl"}, snap.Guards.Shell.AskOn)
				assert.True(t, snap.Guards.MCP.Enabled, "mcp.enabled")
				assert.Equal(t, []string{"untrusted-*"}, snap.Guards.MCP.DenyServers)
				assert.True(t, snap.Guards.Paths.Enabled, "paths.enabled")
				assert.Equal(t, []string{"/etc/shadow"}, snap.Guards.Paths.DenyRead)
				assert.Equal(t, []string{"**/.env"}, snap.Guards.Paths.DenyWrite)
				found := false
				for _, r := range snap.Routes {
					if r.Name != "gate" {
						continue
					}
					found = true
					require.Len(t, r.Sync, 1)
					assert.Equal(t, []string{"secrets", "shell", "mcp", "paths"}, r.Sync[0].Guards)
				}
				assert.True(t, found, "gate route")
			},
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
