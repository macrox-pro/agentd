//go:build windows

package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

// Dial connects to a Windows named pipe.
func Dial(ctx context.Context, path string) (net.Conn, error) {
	conn, err := winio.DialPipeContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("dial pipe: %w", err)
	}
	return conn, nil
}
