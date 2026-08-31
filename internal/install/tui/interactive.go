package tui

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	noTUIEnv = "AGENTD_NO_TUI"
	ciEnv    = "CI"
)

// ErrNonInteractive is returned when a TUI command runs without a terminal.
var ErrNonInteractive = errors.New("non-interactive environment")

// Interactive reports whether the wizard may run on stdout.
func Interactive(getenv func(string) string, stdout *os.File) bool {
	if getenv == nil {
		getenv = os.Getenv
	}
	if getenv(noTUIEnv) == "1" {
		return false
	}
	if strings.EqualFold(getenv(ciEnv), "true") {
		return false
	}
	if stdout == nil {
		return false
	}
	return term.IsTerminal(int(stdout.Fd()))
}
