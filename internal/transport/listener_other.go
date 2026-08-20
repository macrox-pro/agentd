//go:build !unix && !windows

package transport

import (
	"errors"
	"net"
)

// DefaultSocketPath returns empty on unsupported platforms.
func DefaultSocketPath() string {
	return ""
}

// Listen is unsupported on this platform.
func Listen(string) (net.Listener, error) {
	return nil, errors.New("transport: unsupported platform")
}
