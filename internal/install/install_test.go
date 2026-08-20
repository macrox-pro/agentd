package install_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
)

func TestInstallClaudeProject(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	bin := filepath.Join(dir, "agentd")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	err := install.Run(context.Background(), install.Options{
		Provider: "claude-code",
		Scope:    "project",
		Dir:      dir,
		Command:  []string{bin},
	})
	require.NoError(t, err)

	settings := filepath.Join(dir, ".claude", "settings.json")
	b, err := os.ReadFile(settings)
	require.NoError(t, err)
	body := string(b)
	assert.Contains(t, body, bin)
	assert.Contains(t, body, "agenthooks")
	assert.Contains(t, body, "run")
	assert.Contains(t, body, "--provider=claude-code")
}

func TestInstallOpenCodeShim(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	bin := filepath.Join(dir, "agentd")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	err := install.Run(context.Background(), install.Options{
		Provider: "opencode",
		Scope:    "project",
		Dir:      dir,
		Command:  []string{bin},
	})
	require.NoError(t, err)

	shim := filepath.Join(dir, ".opencode", "plugin", "agenthooks.ts")
	b, err := os.ReadFile(shim)
	require.NoError(t, err)
	body := string(b)
	assert.Contains(t, body, bin)
	assert.True(t, strings.Contains(body, `"agenthooks", "serve"`) || strings.Contains(body, "agenthooks\", \"serve\""))
	assert.Contains(t, body, "opencode")
}

func TestInstallUnknownProvider(t *testing.T) {
	t.Parallel()
	err := install.Run(context.Background(), install.Options{
		Provider: "nope",
		Command:  []string{"/bin/agentd"},
		Dir:      t.TempDir(),
	})
	require.Error(t, err)
}
