package hookedge

import (
	"fmt"
	"io"

	"github.com/macrox-pro/agentd/internal/config"
)

const daemonNotRunningMsg = "daemon not running"

// resolveOffline prints daemonNotRunningMsg and returns FailOpen or FailClosed
// from config.OfflineFor. Load/parse errors and FailClosed map to FailClosed.
func resolveOffline(opts Options, cwd string, stderr io.Writer) config.FailMode {
	if stderr == nil {
		stderr = io.Discard
	}
	fmt.Fprintln(stderr, daemonNotRunningMsg)
	mode, err := config.OfflineFor(config.LoadOptions{UserPath: opts.ConfigPath}, cwd)
	if err != nil || mode == config.FailClosed {
		return config.FailClosed
	}
	return config.FailOpen
}
