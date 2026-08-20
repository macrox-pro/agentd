//go:build windows

package transport

import (
	"fmt"

	"golang.org/x/sys/windows"
)

const (
	pipePrefix = `\\.\pipe\agentd`
	pipeSep    = `-`
)

// DefaultSocketPath returns the platform default agentd named pipe path.
// Uses the current process user SID (DESIGN: agentd-<user-sid>).
func DefaultSocketPath() string {
	sid, err := currentUserSID()
	if err != nil || sid == "" {
		return pipePrefix
	}
	return pipePrefix + pipeSep + sid
}

func currentUserSID() (string, error) {
	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return "", fmt.Errorf("open process token: %w", err)
	}
	defer token.Close()

	tokUser, err := token.GetTokenUser()
	if err != nil {
		return "", fmt.Errorf("token user: %w", err)
	}
	return tokUser.User.Sid.String(), nil
}
