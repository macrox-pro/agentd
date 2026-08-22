package daemon_test

import (
	"os"
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
		{
			name: "read pid not running",
			run: func(t *testing.T, paths daemon.Paths) {
				t.Helper()
				_, err := paths.ReadPID()
				require.ErrorIs(t, err, daemon.ErrNotRunning, "ReadPID(missing)")
			},
		},
		{
			name: "invalid pid parse",
			run: func(t *testing.T, paths daemon.Paths) {
				t.Helper()
				require.NoError(t, os.MkdirAll(paths.Dir, 0o700), "MkdirAll(%q)", paths.Dir)
				require.NoError(t, os.WriteFile(paths.PID, []byte("not-a-pid\n"), 0o600), "WriteFile(%q)", paths.PID)

				_, err := paths.ReadPID()
				require.Error(t, err, "ReadPID(invalid)")
				assert.Contains(t, err.Error(), "parse pid", "ReadPID(invalid)")
			},
		},
		{
			name: "write pid dir error",
			run: func(t *testing.T, _ daemon.Paths) {
				t.Helper()
				dir := t.TempDir()
				block := filepath.Join(dir, "block")
				require.NoError(t, os.WriteFile(block, []byte("x"), 0o600), "WriteFile(%q)", block)
				paths := daemon.NewPaths(filepath.Join(block, "s.sock"))
				err := paths.WritePID(1)
				require.Error(t, err, "WritePID(blocked dir)")
			},
		},
		{
			name: "release nil lock",
			run: func(_ *testing.T, _ daemon.Paths) {
				daemon.ReleaseLock(nil)
			},
		},
		{
			name: "lock open error",
			run: func(t *testing.T, paths daemon.Paths) {
				t.Helper()
				require.NoError(t, os.MkdirAll(paths.Dir, 0o700), "MkdirAll(%q)", paths.Dir)
				require.NoError(t, os.Mkdir(paths.Lock, 0o700), "Mkdir(%q)", paths.Lock)
				_, err := paths.AcquireLock()
				require.Error(t, err, "AcquireLock(dir)")
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
