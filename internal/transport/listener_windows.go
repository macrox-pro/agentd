//go:build windows

package transport

import (
	"context"
	"fmt"
	"net"
	"os/user"

	"github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"
)

// DefaultSocketPath returns the platform default agentd named pipe path.
func DefaultSocketPath() string {
	u, err := user.Current()
	if err != nil || u.Uid == "" {
		return `\\.\pipe\agentd`
	}
	return `\\.\pipe\agentd-` + u.Uid
}

// Listen creates a Windows named pipe listener restricted to the current user.
func Listen(path string) (net.Listener, error) {
	sddl, err := currentUserPipeSDDL()
	if err != nil {
		return nil, err
	}
	return winio.ListenPipe(path, &winio.PipeConfig{
		SecurityDescriptor: sddl,
	})
}

// Dial connects to a Windows named pipe.
func Dial(ctx context.Context, path string) (net.Conn, error) {
	return winio.DialPipeContext(ctx, path)
}

func currentUserPipeSDDL() (string, error) {
	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return "", fmt.Errorf("open process token: %w", err)
	}
	defer token.Close()

	tokUser, err := token.GetTokenUser()
	if err != nil {
		return "", fmt.Errorf("token user: %w", err)
	}
	// D:P = protected DACL; (A;;GA;;;SID) = allow generic-all to current user only.
	return fmt.Sprintf("D:P(A;;GA;;;%s)", tokUser.User.Sid.String()), nil
}
