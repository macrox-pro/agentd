//go:build windows

package transport

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"
)

// Listen creates a Windows named pipe listener restricted to the current user.
func Listen(path string) (net.Listener, error) {
	sddl, err := currentUserPipeSDDL()
	if err != nil {
		return nil, err
	}
	ln, err := winio.ListenPipe(path, &winio.PipeConfig{
		SecurityDescriptor: sddl,
	})
	if err != nil {
		return nil, fmt.Errorf("listen pipe: %w", err)
	}
	return ln, nil
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
