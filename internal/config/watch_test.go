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

func TestWatchReload(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600))

	store, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	before := store.Current().Generation

	w, err := store.Watch(config.WatchOptions{Debounce: 50 * time.Millisecond})
	require.NoError(t, err)
	t.Cleanup(func() { _ = w.Close() })

	require.NoError(t, os.WriteFile(path, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	deadline := time.Now().Add(2 * time.Second)
	for {
		snap := store.Current()
		if snap.Generation > before && snap.Policy.Fail == config.FailOpen {
			assert.Greater(t, snap.Generation, before)
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("watch did not reload: gen=%d fail=%s", snap.Generation, snap.Policy.Fail)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
