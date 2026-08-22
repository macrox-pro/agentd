package hookedge

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/provider"
)

// Serve runs the long-lived OpenCode NDJSON bridge: each frame → daemon Invoke → wire reply.
// Never writes debug logs to stdout.
func Serve(ctx context.Context, opts Options) int {
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stdin := opts.Stdin
	if stdin == nil {
		stdin = strings.NewReader("")
	}

	id, err := provider.Parse(opts.Provider)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	if id != provider.OpenCode {
		fmt.Fprintf(stderr, "hook serve supports opencode only, got %q\n", id)
		return 1
	}

	cli, err := hookclient.Dial(ctx, opts.Socket)
	if err != nil {
		fmt.Fprintln(stderr, "daemon not running")
		return 1
	}
	defer cli.Close()

	r := agenthooks.New(agenthooks.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.Use(func(hctx context.Context, typed any, _ agenthooks.Next) (agenthooks.Decision, error) {
		base := agenthooks.EventOf(typed)
		if base == nil || len(base.Raw) == 0 {
			return agenthooks.NoDecision(), nil
		}
		invokeCtx := hctx
		var cancel context.CancelFunc
		if opts.Timeout > 0 {
			invokeCtx, cancel = context.WithTimeout(hctx, opts.Timeout)
			defer cancel()
		}
		req := &agentdv1.InvokeRequest{
			Provider:       agentdv1.Provider_PROVIDER_OPENCODE,
			RawPayload:     append([]byte(nil), base.Raw...),
			InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
		}
		if opts.Timeout > 0 {
			req.Deadline = timestamppb.New(time.Now().Add(opts.Timeout))
		}
		resp, err := cli.Invoke(invokeCtx, req)
		if err != nil {
			return nil, fmt.Errorf("daemon invoke: %w", err)
		}
		return fromProto(resp.GetDecision()), nil
	})

	return r.Run(ctx, []string{"serve", "--provider=opencode"}, stdin, stdout, stderr)
}
