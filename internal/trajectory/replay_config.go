package trajectory

import (
	"context"
	"fmt"
	"time"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

const replayQueueCloseTimeout = 2 * time.Second

// ReplayPolicyConfigOptions configures offline policy replay with config load and engine setup.
type ReplayPolicyConfigOptions struct {
	ConfigPath   string
	SessionsRoot string
	Provider     string
	SessionID    string
	Seq          uint64
}

// ReplayPolicyFromConfig loads config, builds a dispatch engine, and runs ReplayPolicy.
func ReplayPolicyFromConfig(ctx context.Context, opts ReplayPolicyConfigOptions) (ReplayResult, error) {
	store, err := config.Load(ctx, opts.ConfigPath)
	if err != nil {
		return ReplayResult{}, fmt.Errorf("load config: %w", err)
	}
	snap := store.Current()
	q := dispatch.NewQueue(snap.Async, nil)
	defer q.Close(replayQueueCloseTimeout)
	eng := dispatch.NewEngine(q, nil)

	return ReplayPolicy(ctx, ReplayOptions{
		SessionsRoot: opts.SessionsRoot,
		Provider:     opts.Provider,
		SessionID:    opts.SessionID,
		Seq:          opts.Seq,
		Snap:         snap,
		Engine:       eng,
	})
}
