package daemon

import (
	"context"
	"fmt"

	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/transport"
)

// ReloadResult is the config generation after a successful reload.
type ReloadResult struct {
	Generation  uint64
	Fingerprint string
}

// Reload asks the running daemon to reload config from disk.
func Reload(ctx context.Context, socket string) (ReloadResult, error) {
	if socket == "" {
		socket = transport.DefaultSocketPath()
	}
	cli, err := hookclient.Dial(ctx, socket)
	if err != nil {
		return ReloadResult{}, fmt.Errorf("daemon not running: %w", err)
	}
	defer cli.Close()

	resp, err := cli.Reload(ctx)
	if err != nil {
		return ReloadResult{}, err
	}
	var out ReloadResult
	if cfg := resp.GetConfig(); cfg != nil {
		out.Generation = cfg.GetGeneration()
		out.Fingerprint = cfg.GetFingerprint()
	}
	return out, nil
}
