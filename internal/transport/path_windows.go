//go:build windows

package transport

import (
	"os/user"
)

const (
	pipePrefix = `\\.\pipe\agentd`
	pipeSep    = `-`
)

// DefaultSocketPath returns the platform default agentd named pipe path.
func DefaultSocketPath() string {
	u, err := user.Current()
	if err != nil || u.Uid == "" {
		return pipePrefix
	}
	return pipePrefix + pipeSep + u.Uid
}
