//go:build windows

package daemon

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// AcquireLock creates an exclusive lock file using LockFileEx.
func (p Paths) AcquireLock() (*os.File, error) {
	if err := p.ensureDir(); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(p.Lock, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock: %w", err)
	}
	var ol windows.Overlapped
	err = windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &ol)
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("%w: %v", ErrAlreadyRunning, err)
	}
	return f, nil
}

// ReleaseLock unlocks and closes the lock file.
func ReleaseLock(f *os.File) {
	if f == nil {
		return
	}
	var ol windows.Overlapped
	_ = windows.UnlockFileEx(windows.Handle(f.Fd()), 0, 1, 0, &ol)
	_ = f.Close()
}
