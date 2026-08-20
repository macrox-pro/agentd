package config_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestDefaultUserPath(t *testing.T) {
	t.Parallel()

	got := config.DefaultUserPath()
	if got == "" {
		t.Skip("home directory unavailable")
	}
	require.True(t, filepath.IsAbs(got), "DefaultUserPath()=%q", got)
	assert.Equal(t, ".agentd.yaml", filepath.Base(got), "DefaultUserPath()=%q", got)
}
