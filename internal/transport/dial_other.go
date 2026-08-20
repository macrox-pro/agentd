//go:build !unix && !windows

package transport

import (
	"context"
	"net"
)

// Dial is unsupported on this platform.
func Dial(context.Context, string) (net.Conn, error) {
	return nil, ErrUnsupported
}
