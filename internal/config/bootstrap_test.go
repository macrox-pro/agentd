package config_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestPrepareUserConfig_missing_creates_bootstrap(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")

	err := config.PrepareUserConfig(path, nil)
	require.NoError(t, err, "PrepareUserConfig(%q)", path)

	raw, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)
	body := string(raw)
	assert.Contains(t, body, "version: 1")
	assert.Contains(t, body, "fail_closed")
	assert.NotContains(t, body, "null:")

	store, err := config.Load(t.Context(), path)
	require.NoError(t, err, "Load(%q)", path)
	require.NotNil(t, store.Current())
}

func TestPrepareUserConfig_empty_file_unchanged(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, nil, 0o600), "WriteFile(%q)", path)
	before, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)

	err = config.PrepareUserConfig(path, nil)
	require.NoError(t, err, "PrepareUserConfig(%q)", path)

	after, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)
	assert.Equal(t, before, after, "PrepareUserConfig(%q)", path)
}

func TestPrepareUserConfig_valid_unchanged(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	content := "version: 1\npolicy:\n  fail: fail_closed\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600), "WriteFile(%q)", path)

	err := config.PrepareUserConfig(path, nil)
	require.NoError(t, err, "PrepareUserConfig(%q)", path)

	raw, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)
	assert.Equal(t, content, string(raw), "PrepareUserConfig(%q)", path)
}

func TestPrepareUserConfig_invalid_returns_error(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	content := "version: 1\npolicy:\n  fail: nope\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600), "WriteFile(%q)", path)

	var notify bytes.Buffer
	err := config.PrepareUserConfig(path, &notify)
	require.Error(t, err, "PrepareUserConfig(%q)", path)
	assert.Contains(t, notify.String(), path)
	assert.Contains(t, notify.String(), "invalid user config")

	raw, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)
	assert.Equal(t, content, string(raw), "PrepareUserConfig(%q)", path)
}

func TestPrepareUserConfig_invalid_parse_notify(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	content := ":\n  bad:\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600), "WriteFile(%q)", path)

	var notify bytes.Buffer
	err := config.PrepareUserConfig(path, &notify)
	require.Error(t, err, "PrepareUserConfig(%q)", path)
	assert.Contains(t, notify.String(), "invalid user config")
	assert.NotContains(t, notify.String(), "read config:")
}

func TestPrepareUserConfig_bootstrap_write_fails(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	path := filepath.Join(dir, "agentd.yaml")

	err := config.PrepareUserConfig(path, nil)
	require.Error(t, err, "PrepareUserConfig(%q)", path)
	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr), "Stat(%q)", path)
}

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
	assert.NotContains(t, notify.String(), "invalid user config")
}

func TestPrepareUserConfig_idempotent_twice(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")

	require.NoError(t, config.PrepareUserConfig(path, nil), "PrepareUserConfig(%q)", path)
	first, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)

	require.NoError(t, config.PrepareUserConfig(path, nil), "PrepareUserConfig(%q)", path)
	second, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)
	assert.Equal(t, string(first), string(second), "PrepareUserConfig(%q)", path)
}

func TestLoad_missing_does_not_create_file(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "missing.yaml")

	store, err := config.Load(t.Context(), path)
	require.NoError(t, err, "Load(%q)", path)
	require.NotNil(t, store.Current())

	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err), "Stat(%q)", path)
}

func TestLayerYAML_user_omitempty(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))

	store, err := config.Load(t.Context(), path)
	require.NoError(t, err, "Load(%q)", path)

	raw, err := store.LayerYAML(config.LayerUser, "", "")
	require.NoError(t, err, "LayerYAML(user)")
	assert.NotContains(t, string(raw), "null:")
	assert.Contains(t, string(raw), "version: 1")
}

func TestLayerYAML_user_missing_empty(t *testing.T) {
	t.Parallel()
	store, err := config.Load(t.Context(), filepath.Join(t.TempDir(), "missing.yaml"))
	require.NoError(t, err, "Load")

	raw, err := store.LayerYAML(config.LayerUser, "", "")
	require.NoError(t, err, "LayerYAML(user)")
	assert.Empty(t, raw)
}

func TestLayerYAML_user_null_keys_normalized(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	content := "version: 1\npolicy:\n  fail: fail_closed\nasync: null\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	store, err := config.Load(t.Context(), path)
	require.NoError(t, err, "Load(%q)", path)

	raw, err := store.LayerYAML(config.LayerUser, "", "")
	require.NoError(t, err, "LayerYAML(user)")
	body := string(raw)
	assert.NotContains(t, body, "null")
	assert.NotContains(t, body, "async:")
}

func TestPrepareUserConfig_directory(t *testing.T) {
	t.Parallel()
	path := t.TempDir()
	var notify bytes.Buffer
	err := config.PrepareUserConfig(path, &notify)
	require.Error(t, err, "PrepareUserConfig(%q)", path)
	assert.NotContains(t, notify.String(), "invalid user config")
}

func TestPrepareUserConfig_empty_path(t *testing.T) {
	t.Parallel()
	err := config.PrepareUserConfig("", nil)
	require.Error(t, err, "PrepareUserConfig(\"\")")
	assert.Contains(t, err.Error(), "home directory unavailable")
}

func TestBootstrapUserFileConfig_no_null_in_marshal(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "agentd.yaml")
	require.NoError(t, config.PrepareUserConfig(path, nil), "PrepareUserConfig(%q)", path)
	raw, err := os.ReadFile(path)
	require.NoError(t, err, "ReadFile(%q)", path)
	assert.False(t, strings.Contains(string(raw), "null"), "bootstrap yaml")
}
