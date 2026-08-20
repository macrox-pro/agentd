//go:build unix

package transport

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// Listen creates a Unix domain socket listener.
func Listen(path string) (net.Listener, error) {
	if err := os.MkdirAll(filepath.Dir(path), socketDirPerm); err != nil {
		return nil, fmt.Errorf("mkdir socket dir: %w", err)
	}
	_ = os.Remove(path)
	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("listen unix: %w", err)
	}
	if err := os.Chmod(path, socketFilePerm); err != nil {
		_ = ln.Close()
		return nil, fmt.Errorf("chmod socket: %w", err)
	}
	return ln, nil
}
