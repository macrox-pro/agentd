//go:build windows

package daemon

import (
	"os"
	"path/filepath"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/transport"
)

const stateFallbackSubdir = "agentd"

// stateDir keeps pid and lock next to a file socket, but a named pipe has no
// filesystem parent: pipe endpoints store state in the per-user state directory.
func stateDir(socket string) string {
	if !transport.IsPipePath(socket) {
		return filepath.Dir(socket)
	}
	if dir := config.DefaultStateDir(); dir != "" {
		return dir
	}
	return filepath.Join(os.TempDir(), stateFallbackSubdir)
}
