//go:build !unix && !windows

package transport

import (
	"context"
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

// Dial is unsupported on this platform.
func Dial(context.Context, string) (net.Conn, error) {
	return nil, errors.New("transport: unsupported platform")
}
