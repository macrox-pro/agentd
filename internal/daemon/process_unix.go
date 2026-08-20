//go:build unix

package daemon

import "syscall"

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	return syscall.Kill(pid, 0) == nil
}

func signalTerminate(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
