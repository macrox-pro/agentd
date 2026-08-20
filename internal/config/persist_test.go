package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestPersistFlushWritesRuntime(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	runtimePath := filepath.Join(dir, "agentd", "runtime.yaml")
	store, err := config.LoadWith(t.Context(), config.LoadOptions{RuntimePath: runtimePath})
	require.NoError(t, err)

	fp := config.ApprovalFingerprint(config.ApprovalKindShell, "Bash", "curl")
	require.NoError(t, store.RecordDecision(config.RecordDecisionOptions{
		Fingerprint: fp,
		Scope:       config.ApprovalScopeSession,
		SessionID:   "s1",
	}))
	require.NoError(t, store.FlushRuntime())

	raw, err := os.ReadFile(runtimePath)
	require.NoError(t, err)
	assert.Contains(t, string(raw), fp)
	assert.Contains(t, string(raw), "session_id: s1")

	// Reload from disk preserves approvals.
	store2, err := config.LoadWith(t.Context(), config.LoadOptions{RuntimePath: runtimePath})
	require.NoError(t, err)
	assert.True(t, store2.Current().Approvals.HasApproval(
		config.ApprovalKindShell, fp, "", "s1", time.Now().UTC()))
}

func TestPersistDebounced(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	runtimePath := filepath.Join(dir, "runtime.yaml")
	store, err := config.LoadWith(t.Context(), config.LoadOptions{RuntimePath: runtimePath})
	require.NoError(t, err)

	require.NoError(t, store.PatchRuntime([]byte("version: 1\nasync:\n  queue_capacity: 9\n")))
	// Before debounce window the file may not exist yet; wait then check.
	deadline := time.Now().Add(2 * time.Second)
	var raw []byte
	for time.Now().Before(deadline) {
		raw, err = os.ReadFile(runtimePath)
		if err == nil && len(raw) > 0 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	require.NoError(t, err)
	assert.Contains(t, string(raw), "queue_capacity: 9")
}
