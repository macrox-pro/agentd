package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Paths holds filesystem locations for one daemon instance.
type Paths struct {
	Socket string
	PID    string
	Lock   string
	Dir    string
}

// NewPaths derives state file paths from a socket path. The state directory is
// platform-specific (paths_unix.go / paths_windows.go / paths_other.go).
func NewPaths(socket string) Paths {
	dir := stateDir(socket)
	return Paths{
		Socket: socket,
		PID:    filepath.Join(dir, "agentd.pid"),
		Lock:   filepath.Join(dir, "agentd.lock"),
		Dir:    dir,
	}
}

func (p Paths) ensureDir() error {
	return os.MkdirAll(p.Dir, 0o700)
}

// WritePID writes the process ID to the pid file.
func (p Paths) WritePID(pid int) error {
	if err := p.ensureDir(); err != nil {
		return err
	}
	return os.WriteFile(p.PID, []byte(strconv.Itoa(pid)+"\n"), 0o600)
}

// ReadPID reads the PID file.
func (p Paths) ReadPID() (int, error) {
	b, err := os.ReadFile(p.PID)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrNotRunning
		}
		return 0, err
	}
	s := strings.TrimSpace(string(b))
	pid, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse pid: %w", err)
	}
	return pid, nil
}

// ClearPID removes the PID file so Stop waiters observe shutdown promptly.
func (p Paths) ClearPID() {
	_ = os.Remove(p.PID)
}

// RemoveStale removes PID and socket files after shutdown.
func (p Paths) RemoveStale() {
	p.ClearPID()
	_ = os.Remove(p.Socket)
}
