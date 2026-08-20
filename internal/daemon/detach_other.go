//go:build !unix && !windows

package daemon

import (
	"errors"
	"os/exec"
)

func configureDetach(*exec.Cmd) error {
	return errors.New("daemon detach: unsupported platform")
}
