//go:build !unix && !windows

package daemon

import "errors"

func processAlive(int) bool { return false }

func signalTerminate(int) error {
	return errors.New("daemon signal: unsupported platform")
}
