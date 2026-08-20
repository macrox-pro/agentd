//go:build unix

package config_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestDefaultRuntimePath(t *testing.T) {
	t.Parallel()
	got := config.DefaultRuntimePath()
	if got == "" {
		t.Skip("runtime path unavailable")
	}
	require.True(t, filepath.IsAbs(got), "DefaultRuntimePath()=%q", got)
	assert.Equal(t, "runtime.yaml", filepath.Base(got), "DefaultRuntimePath()=%q", got)
}
