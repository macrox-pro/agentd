package daemon_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestDaemonPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(t *testing.T, paths daemon.Paths)
	}{
		{
			name: "exclusive lock",
			run: func(t *testing.T, paths daemon.Paths) {
				t.Helper()
				lock, err := paths.AcquireLock()
				require.NoError(t, err, "AcquireLock()")
				t.Cleanup(func() { daemon.ReleaseLock(lock) })

				_, err = paths.AcquireLock()
				require.Error(t, err, "second AcquireLock()")
			},
		},
		{
			name: "pid round trip",
			run: func(t *testing.T, paths daemon.Paths) {
				t.Helper()
				const wantPID = 12345
				require.NoError(t, paths.WritePID(wantPID), "WritePID(%d)", wantPID)

				got, err := paths.ReadPID()
				require.NoError(t, err, "ReadPID()")
				assert.Equal(t, wantPID, got, "ReadPID()")
			},
		},
		{
			name: "stale pid removed",
			run: func(t *testing.T, paths daemon.Paths) {
				t.Helper()
				require.NoError(t, paths.WritePID(999999999), "WritePID(stale)")

				_, err := paths.ReadPID()
				require.NoError(t, err, "ReadPID(stale)")

				paths.RemoveStale()
				_, err = paths.ReadPID()
				require.ErrorIs(t, err, daemon.ErrNotRunning, "ReadPID() after RemoveStale")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			paths := daemon.NewPaths(filepath.Join(dir, "agentd.sock"))
			tt.run(t, paths)
		})
	}
}
