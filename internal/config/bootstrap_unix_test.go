//go:build unix

package config_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestPrepareUserConfig_unreadable(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o400), "WriteFile(%q)", path)
	require.NoError(t, os.Chmod(path, 0o000), "Chmod(%q)", path)
	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })

	var notify bytes.Buffer
	err := config.PrepareUserConfig(path, &notify)
	require.Error(t, err, "PrepareUserConfig(%q)", path)
	assert.False(t, errors.Is(err, config.ErrParseConfig))
	assert.NotContains(t, notify.String(), "invalid user config")
}
