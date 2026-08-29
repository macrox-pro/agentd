package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestDaemonDisable_success(t *testing.T) {
	t.Cleanup(daemon.SetAutostartHooksForTest(&daemon.AutostartHooks{
		Unregister: func() error { return nil },
	}))
	got := executeRoot(t, execOpts{args: []string{"daemon", "disable"}})
	require.NoError(t, got.err)
}

func TestDaemonDisable_idempotent(t *testing.T) {
	n := 0
	t.Cleanup(daemon.SetAutostartHooksForTest(&daemon.AutostartHooks{
		Unregister: func() error {
			n++
			return nil
		},
	}))
	require.NoError(t, executeRoot(t, execOpts{args: []string{"daemon", "disable"}}).err)
	require.NoError(t, executeRoot(t, execOpts{args: []string{"daemon", "disable"}}).err)
	assert.Equal(t, 2, n)
}
