//go:build unix

package transport

import "net"

// DefaultSocketPath returns the platform default agentd socket path.
func DefaultSocketPath() string {
	// Resolved at runtime from XDG_RUNTIME_DIR in daemon (M1).
	return ""
}

// Listen creates a Unix domain socket listener for gRPC.
func Listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}
