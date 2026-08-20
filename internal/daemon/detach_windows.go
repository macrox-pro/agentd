//go:build windows

package daemon

import (
	"os/exec"
	"syscall"
)

// DETACHED_PROCESS is not exported by syscall on all Go versions.
const detachedProcess = 0x00000008

func configureDetach(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | detachedProcess,
	}
	return nil
}
