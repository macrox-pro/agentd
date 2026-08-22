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

// Run dials the daemon, invokes HookService, then encodes the decision wire response.
// Returns a process exit code (0 success, 1 failure). Never writes debug logs to stdout.
func Run(ctx context.Context, opts Options) int {
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}

	id, err := provider.Parse(opts.Provider)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	protoProv, err := id.Proto()
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	payload, mode, err := readPayload(opts)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

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
		InvocationMode: mode,
		Cwd:            ResolveCWD(payload),
	}
	if opts.Timeout > 0 {
		req.Deadline = timestamppb.New(time.Now().Add(opts.Timeout))
	}

	resp, err := cli.Invoke(invokeCtx, req)
	if err != nil {
		fmt.Fprintln(stderr, "daemon not running")
		return 1
	}
	decision := resp.GetDecision()
	if decision == nil {
		decision = &agentdv1.Decision{Kind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION}
	}

	return encodeDecision(ctx, opts.Provider, opts.ArgvPayload, payload, decision, stdout, stderr)
}
