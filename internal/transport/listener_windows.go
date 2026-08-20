//go:build windows

package transport

import (
	"context"
	"net"
	"os/user"

	"github.com/Microsoft/go-winio"
)

// DefaultSocketPath returns the platform default agentd named pipe path.
func DefaultSocketPath() string {
	u, err := user.Current()
	if err != nil || u.Uid == "" {
		return `\\.\pipe\agentd`
	}
	return `\\.\pipe\agentd-` + u.Uid
}

// Listen creates a Windows named pipe listener.
func Listen(path string) (net.Listener, error) {
	return winio.ListenPipe(path, nil)
}

// Dial connects to a Windows named pipe.
func Dial(ctx context.Context, path string) (net.Conn, error) {
	return winio.DialPipeContext(ctx, path)
}
