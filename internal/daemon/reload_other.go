//go:build !unix

package daemon

import "os"

func notifyReload(chan<- os.Signal) {}

func isReloadSignal(os.Signal) bool { return false }
