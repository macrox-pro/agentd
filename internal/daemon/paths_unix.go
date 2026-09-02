//go:build unix

package daemon

import "path/filepath"

func stateDir(socket string) string {
	return filepath.Dir(socket)
}
