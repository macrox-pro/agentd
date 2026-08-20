//go:build unix

package daemon

import (
	"fmt"
	"os"
	"syscall"
)

// AcquireLock creates an exclusive lock file. Caller must ReleaseLock.
func (p Paths) AcquireLock() (*os.File, error) {
	if err := p.ensureDir(); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(p.Lock, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
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
	_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	_ = f.Close()
}
