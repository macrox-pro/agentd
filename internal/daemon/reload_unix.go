//go:build unix

package daemon

import (
	"os"
	"os/signal"
	"syscall"
)

func notifyReload(ch chan<- os.Signal) {
	signal.Notify(ch, syscall.SIGHUP)
}

func isReloadSignal(sig os.Signal) bool {
	return sig == syscall.SIGHUP
}
