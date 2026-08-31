package install_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
	"github.com/macrox-pro/agentd/internal/provider"
)

func TestPlan_hookStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, env install.DiscoverEnv, bin string) []install.Target
		mutate func(t *testing.T, targets []install.Target)
		want   install.HookStatus
	}{
		{
			name: "hook_status_missing",
			setup: func(t *testing.T, env install.DiscoverEnv, bin string) []install.Target {
				t.Helper()
				require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".cursor"), 0o755))
				return highConfidenceTargets(t, env)
			},
			want: install.HookStatusMissing,
		},
		{
			name: "hook_status_current",
			setup: func(t *testing.T, env install.DiscoverEnv, bin string) []install.Target {
				t.Helper()
				require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".cursor"), 0o755))
				targets := highConfidenceTargets(t, env)
				_, err := install.RunAll(context.Background(), targets, []string{bin}, false)
				require.NoError(t, err, "RunAll(%q)", "hook_status_current")
				return targets
			},
			want: install.HookStatusCurrent,
		},
		{
			name: "hook_status_stale",
			setup: func(t *testing.T, env install.DiscoverEnv, bin string) []install.Target {
				t.Helper()
				require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".codex"), 0o755))
				targets := highConfidenceTargets(t, env)
				_, err := install.RunAll(context.Background(), targets, []string{bin}, false)
				require.NoError(t, err, "RunAll(%q)", "hook_status_stale")
				return targets
			},
			mutate: func(t *testing.T, targets []install.Target) {
				t.Helper()
				require.NoError(t, os.Remove(filepath.Join(targets[0].Dir, "hooks.json")))
			},
			want: install.HookStatusStale,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			env := testDiscoverEnv(t, root)
			bin := filepath.Join(root, "agentd")
			require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))
			targets := tt.setup(t, env, bin)
			if tt.mutate != nil {
				tt.mutate(t, targets)
			}
			entries, err := install.Plan(context.Background(), targets, []string{bin})
			require.NoError(t, err, "Plan(%q)", tt.name)
			require.NotEmpty(t, entries, "Plan(%q)", tt.name)
			assert.Equal(t, tt.want, entries[0].Status, "Plan(%q) status", tt.name)
		})
	}
}

func TestPlan_multiTarget(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	env := testDiscoverEnv(t, root)
	require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".cursor"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(env.Home, ".claude"), 0o755))
	bin := filepath.Join(root, "agentd")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	targets := highConfidenceTargets(t, env)
	require.GreaterOrEqual(t, len(targets), 2)

	entries, err := install.Plan(context.Background(), targets, []string{bin})
	require.NoError(t, err, "Plan(multi)")
	assert.Len(t, entries, len(targets))
	providers := make([]provider.ID, len(entries))
	for i, e := range entries {
		providers[i] = e.Target.Provider
		assert.Equal(t, install.HookStatusMissing, e.Status)
	}
	assert.Contains(t, providers, provider.Cursor)
	assert.Contains(t, providers, provider.ClaudeCode)
}

func TestRunAll_twoProviders(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	env := testDiscoverEnv(t, root)
	require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".cursor"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(env.Home, ".claude"), 0o755))
	bin := filepath.Join(root, "agentd")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	targets := highConfidenceTargets(t, env)
	require.GreaterOrEqual(t, len(targets), 2)

	results, err := install.RunAll(context.Background(), targets, []string{bin}, false)
	require.NoError(t, err, "RunAll(two providers)")
	assert.Len(t, results, len(targets))
}

func TestRunAll_partialError(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	env := testDiscoverEnv(t, root)
	require.NoError(t, os.MkdirAll(filepath.Join(env.Cwd, ".cursor"), 0o755))
	bin := filepath.Join(root, "agentd")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	ok := highConfidenceTargets(t, env)
	require.NotEmpty(t, ok)
	targets := append(append([]install.Target(nil), ok...), install.Target{
		Provider: provider.Cursor,
		Scope:    "nope",
		Dir:      env.Cwd,
	})

	_, err := install.RunAll(context.Background(), targets, []string{bin}, false)
	require.Error(t, err, "RunAll(partial)")
	assert.Contains(t, err.Error(), "unknown scope")
	_, statErr := os.Stat(filepath.Join(ok[0].Dir, ".cursor", "hooks.json"))
	require.NoError(t, statErr, "first target written before partial error")
}

func highConfidenceTargets(t *testing.T, env install.DiscoverEnv) []install.Target {
	t.Helper()
	findings, err := install.Discover(context.Background(), env)
	require.NoError(t, err, "Discover")
	targets, err := install.TargetsFromHighConfidence(findings, env)
	require.NoError(t, err, "TargetsFromHighConfidence")
	return targets
}
