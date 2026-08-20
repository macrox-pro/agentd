//go:build unix

package transport

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	socketDirPerm  = 0o700
	socketFilePerm = 0o600

	runtimeSubdir = "agentd"
	socketName    = "agentd.sock"
	darwinCache   = "Library/Caches"
	linuxRun      = ".local/run"
)

// DefaultSocketPath returns the platform default agentd socket path.
func DefaultSocketPath() string {
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, runtimeSubdir, socketName)
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(os.TempDir(), runtimeSubdir, socketName)
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, darwinCache, runtimeSubdir, socketName)
	}
	return filepath.Join(home, linuxRun, runtimeSubdir, socketName)
}
