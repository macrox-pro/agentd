//go:build windows

package transport

import (
	"net"

	"github.com/Microsoft/go-winio"
)

// DefaultSocketPath returns the platform default agentd named pipe address.
func DefaultSocketPath() string {
	return `\\.\pipe\agentd`
}

// Listen creates a Windows named pipe listener for gRPC.
func Listen(path string) (net.Listener, error) {
	return winio.ListenPipe(path, nil)
}
