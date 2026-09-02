//go:build windows

package daemon_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestPathsStateDirWindows(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, local string) (socket, wantDir string)
	}{
		{
			name: "pipe endpoint uses state dir",
			setup: func(_ *testing.T, local string) (string, string) {
				return `\\.\pipe\agentd-test-state`, filepath.Join(local, "agentd")
			},
		},
		{
			name: "file endpoint keeps socket dir",
			setup: func(t *testing.T, _ string) (string, string) {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "s.sock"), dir
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local := t.TempDir()
			t.Setenv("LOCALAPPDATA", local)
			socket, wantDir := tt.setup(t, local)

			paths := daemon.NewPaths(socket)
			assert.Equal(t, wantDir, paths.Dir, "NewPaths(%q).Dir", socket)
			assert.Equal(t, filepath.Join(wantDir, "agentd.pid"), paths.PID, "NewPaths(%q).PID", socket)
			assert.Equal(t, filepath.Join(wantDir, "agentd.lock"), paths.Lock, "NewPaths(%q).Lock", socket)

			const wantPID = 4321
			require.NoError(t, paths.WritePID(wantPID), "WritePID(%q)", socket)
			got, err := paths.ReadPID()
			require.NoError(t, err, "ReadPID(%q)", socket)
			assert.Equal(t, wantPID, got, "ReadPID(%q)", socket)
		})
	}
}
