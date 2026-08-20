//go:build unix

package daemon

import (
	"os/exec"
	"syscall"
)

func configureDetach(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	return nil
}
