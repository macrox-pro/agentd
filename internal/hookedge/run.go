package hookedge

import (
	"bytes"
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
)

// Options configures a single hook run.
type Options struct {
	Socket      string
	Provider    string
	ArgvPayload bool
	Timeout     time.Duration
	PayloadArg  string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
}

// Run dials the daemon, invokes HookService, then encodes a NoDecision wire response.
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

	provider, err := parseProvider(opts.Provider)
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
		Provider:       provider,
		RawPayload:     payload,
		InvocationMode: mode,
	}
	if opts.Timeout > 0 {
		req.Deadline = timestamppb.New(time.Now().Add(opts.Timeout))
	}

	resp, err := cli.Invoke(invokeCtx, req)
	if err != nil {
		fmt.Fprintln(stderr, "daemon not running")
		return 1
	}
	if resp.GetDecision().GetKind() != agentdv1.DecisionKind_DECISION_KIND_NO_DECISION {
		fmt.Fprintln(stderr, "unsupported decision from daemon")
		return 1
	}

	return encodeDaemonNoDecision(ctx, opts.Provider, opts.ArgvPayload, payload, stdout, stderr)
}

func readPayload(opts Options) ([]byte, agentdv1.InvocationMode, error) {
	if opts.ArgvPayload {
		if opts.PayloadArg == "" {
			return nil, 0, fmt.Errorf("empty argv payload")
		}
		return []byte(opts.PayloadArg), agentdv1.InvocationMode_INVOCATION_MODE_ARGV, nil
	}
	stdin := opts.Stdin
	if stdin == nil {
		stdin = bytes.NewReader(nil)
	}
	payload, err := io.ReadAll(io.LimitReader(stdin, 8<<20))
	if err != nil {
		return nil, 0, fmt.Errorf("read stdin: %w", err)
	}
	if len(payload) == 0 {
		return nil, 0, fmt.Errorf("empty stdin")
	}
	return payload, agentdv1.InvocationMode_INVOCATION_MODE_STDIN, nil
}

func parseProvider(s string) (agentdv1.Provider, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "claude-code":
		return agentdv1.Provider_PROVIDER_CLAUDE_CODE, nil
	case "cursor":
		return agentdv1.Provider_PROVIDER_CURSOR, nil
	case "codex":
		return agentdv1.Provider_PROVIDER_CODEX, nil
	case "gemini":
		return agentdv1.Provider_PROVIDER_GEMINI, nil
	case "opencode":
		return agentdv1.Provider_PROVIDER_OPENCODE, nil
	case "kimicode", "kimi-code":
		return agentdv1.Provider_PROVIDER_KIMI_CODE, nil
	case "":
		return 0, fmt.Errorf("provider is required")
	default:
		return 0, fmt.Errorf("unknown provider %q", s)
	}
}

// encodeDaemonNoDecision encodes provider wire for a daemon-confirmed NoDecision.
// agenthooks does not export Encode; Runner.Run is the supported wire path.
// Handlers always return NoDecision so the edge cannot invent Allow/Deny —
// Decide remains in the daemon; richer decision mapping is M2.
func encodeDaemonNoDecision(ctx context.Context, provider string, argvPayload bool, payload []byte, stdout, stderr io.Writer) int {
	r := agenthooks.New(agenthooks.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.OnToolPre(func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		return agenthooks.NoDecision(), nil
	})
	r.OnPermission(func(context.Context, *agenthooks.PermissionEvent) (agenthooks.ToolPreDecision, error) {
		return agenthooks.NoDecision(), nil
	})
	r.OnToolPost(func(context.Context, *agenthooks.ToolPostEvent) (agenthooks.ToolPostDecision, error) {
		return agenthooks.Observed(), nil
	})
	r.OnToolError(func(context.Context, *agenthooks.ToolPostEvent) (agenthooks.ToolPostDecision, error) {
		return agenthooks.Observed(), nil
	})

	args := []string{"run", "--provider=" + provider}
	var stdin io.Reader
	if argvPayload {
		args = append(args, "--argv-payload", string(payload))
		stdin = bytes.NewReader(nil)
	} else {
		stdin = bytes.NewReader(payload)
	}
	return r.Run(ctx, args, stdin, stdout, stderr)
}
