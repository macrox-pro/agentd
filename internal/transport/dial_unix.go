//go:build unix

package transport

import (
	"context"
	"fmt"
	"net"
)

// Dial connects to a Unix domain socket.
func Dial(ctx context.Context, path string) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", path)
	if err != nil {
		return nil, fmt.Errorf("dial unix: %w", err)
	}
	return conn, nil
}
