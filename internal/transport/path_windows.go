//go:build windows

package transport

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	pipeNamespace = `\\.\pipe\`
	pipeBaseName  = "agentd"
	pipeSep       = `-`
)

// DefaultSocketPath returns the platform default agentd named pipe path.
// Uses the current process user SID (DESIGN: agentd-<user-sid>).
func DefaultSocketPath() string {
	sid, err := currentUserSID()
	if err != nil || sid == "" {
		return pipeNamespace + pipeBaseName
	}
	return pipeNamespace + pipeBaseName + pipeSep + sid
}

// IsPipePath reports whether path names a Windows named pipe rather than a file.
func IsPipePath(path string) bool {
	if len(path) <= len(pipeNamespace) {
		return false
	}
	// The pipe namespace is case-insensitive: \\.\PIPE\name is the same pipe.
	return strings.EqualFold(path[:len(pipeNamespace)], pipeNamespace)
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
