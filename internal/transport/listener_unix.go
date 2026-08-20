//go:build unix

package transport

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

// DefaultSocketPath returns the platform default agentd socket path.
func DefaultSocketPath() string {
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, "agentd", "agentd.sock")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(os.TempDir(), "agentd", "agentd.sock")
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Caches", "agentd", "agentd.sock")
	}
	return filepath.Join(home, ".local", "run", "agentd", "agentd.sock")
}

// Listen creates a Unix domain socket listener.
func Listen(path string) (net.Listener, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("mkdir socket dir: %w", err)
	}
	_ = os.Remove(path)
	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(path, 0o600); err != nil {
		_ = ln.Close()
		return nil, fmt.Errorf("chmod socket: %w", err)
	}
	return ln, nil
}

// Dial connects to a Unix domain socket.
func Dial(ctx context.Context, path string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "unix", path)
}
