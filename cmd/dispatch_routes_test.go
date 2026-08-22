package cmd_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/cmd"
	"github.com/macrox-pro/agentd/internal/config"
)

func TestDispatchRoutesHuman(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), ".agentd.yaml")

	root := cmd.RootCommand()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--config", cfg, "dispatch", "routes"})
	err := root.Execute()
	require.NoError(t, err, "dispatch routes")
	assert.Contains(t, buf.String(), "mode=", "dispatch routes human output")
}

func TestDispatchRoutesJSON(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), ".agentd.yaml")

	root := cmd.RootCommand()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--config", cfg, "dispatch", "routes", "--json"})
	err := root.Execute()
	require.NoError(t, err, "dispatch routes --json")

	var routes []config.CompiledRoute
	require.NoError(t, json.Unmarshal(buf.Bytes(), &routes), "dispatch routes json")
	require.NotEmpty(t, routes, "dispatch routes json")
	assert.NotEmpty(t, routes[0].Name, "dispatch routes json")
}
