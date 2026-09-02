//go:build !unix && !windows

package daemon

import "path/filepath"

func stateDir(socket string) string {
	return filepath.Dir(socket)
}
