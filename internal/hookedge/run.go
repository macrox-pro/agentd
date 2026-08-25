package hookedge

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/decision"
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
	cwd := ResolveCWD(payload)

	invokeCtx := ctx
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		invokeCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	cli, err := hookclient.DialReady(invokeCtx, opts.Socket)
	if err != nil {
		return runOffline(ctx, opts, payload, cwd, stdout, stderr)
	}
	defer cli.Close()

	req := &agentdv1.InvokeRequest{
		Provider:       protoProv,
		RawPayload:     payload,
		InvocationMode: mode,
		Cwd:            cwd,
	}
	if opts.Timeout > 0 {
		req.Deadline = timestamppb.New(time.Now().Add(opts.Timeout))
	}

	resp, err := cli.Invoke(invokeCtx, req)
	if err != nil {
		return runOffline(ctx, opts, payload, cwd, stdout, stderr)
	}
	d := resp.GetDecision()
	if d == nil {
		d = decision.Neutral()
	}

	return encodeDecision(ctx, opts.Provider, opts.ArgvPayload, payload, d, stdout, stderr)
}

func runOffline(ctx context.Context, opts Options, payload []byte, cwd string, stdout, stderr io.Writer) int {
	if resolveOffline(opts, cwd, stderr) == config.FailClosed {
		return 1
	}
	return encodeDecision(ctx, opts.Provider, opts.ArgvPayload, payload, decision.Neutral(), stdout, stderr)
}
