package hookedge

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/provider"
)

// Notify handles Codex notify-style hooks (argv JSON). Always async semantics on the daemon.
// Prints no provider decision wire; exit 0 on success.
func Notify(ctx context.Context, opts Options) int {
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	id, err := provider.Parse(opts.Provider)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	if id != provider.Codex {
		fmt.Fprintf(stderr, "hook notify supports codex only, got %q\n", id)
		return 1
	}
	protoProv, err := id.Proto()
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	if opts.PayloadArg == "" {
		fmt.Fprintln(stderr, "notify payload required")
		return 1
	}
	payload := []byte(opts.PayloadArg)

	invokeCtx := ctx
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		invokeCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	cli, err := hookclient.Dial(invokeCtx, opts.Socket)
	if err != nil {
		fmt.Fprintln(stderr, "daemon not running")
		return 1
	}
	defer cli.Close()

	req := &agentdv1.InvokeRequest{
		Provider:       protoProv,
		RawPayload:     payload,
		InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY,
	}
	if opts.Timeout > 0 {
		req.Deadline = timestamppb.New(time.Now().Add(opts.Timeout))
	}

	if _, err := cli.Invoke(invokeCtx, req); err != nil {
		fmt.Fprintln(stderr, "daemon not running")
		return 1
	}
	return 0
}
