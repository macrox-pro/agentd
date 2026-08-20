//go:build !unix && !windows

package transport

import "net"

// Listen is unsupported on this platform.
func Listen(string) (net.Listener, error) {
	return nil, ErrUnsupported
}
