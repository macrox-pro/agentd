package daemon_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestClaudeProjectsRootConfigured(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "/custom/projects", daemon.ClaudeProjectsRootForTest("/custom/projects"))
}

func TestClaudeProjectsRootDefault(t *testing.T) {
	t.Parallel()
	root := daemon.ClaudeProjectsRootForTest("")
	assert.Contains(t, root, ".claude")
	assert.Contains(t, root, "projects")
}
