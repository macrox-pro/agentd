package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestFingerprintStableAcrossLoads(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "user.yaml")
	require.NoError(t, os.WriteFile(path, []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	a, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	b, err := config.Load(context.Background(), path)
	require.NoError(t, err)
	assert.Equal(t, a.Current().Fingerprint, b.Current().Fingerprint)

	require.NoError(t, os.WriteFile(path, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))
	require.NoError(t, a.Reload(context.Background()))
	assert.NotEqual(t, b.Current().Fingerprint, a.Current().Fingerprint)
}
