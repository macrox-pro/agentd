//go:build !unix && !windows

package daemon

import (
	"errors"
	"os"
)

func (Paths) AcquireLock() (*os.File, error) {
	return nil, errors.New("daemon lock: unsupported platform")
}

func ReleaseLock(*os.File) {}
