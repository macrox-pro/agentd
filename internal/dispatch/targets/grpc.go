package targets

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/speakeasy-api/agenthooks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/transport"
)

const defaultGRPCTimeout = 30 * time.Second

// GRPC forwards HookService.Invoke to a peer endpoint.
type GRPC struct {
	Logger *slog.Logger
	// InvokePeer, if set, replaces the default dial+Invoke (tests).
	InvokePeer func(ctx context.Context, endpoint string, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error)
}

// SyncRequest is one sync grpc forward.
type SyncRequest struct {
	Provider agentdv1.Provider
	Raw      []byte
	Target   config.CompiledTarget
}

// InvokeSync dials the peer and returns its decision.
func (t *GRPC) InvokeSync(ctx context.Context, req SyncRequest) (agenthooks.Decision, error) {
	resp, err := t.invoke(ctx, req.Target, req.Provider, req.Raw)
	if err != nil {
		return nil, err
	}
	return decisionFromProto(resp.GetDecision()), nil
}

// InvokeAsync forwards Invoke and discards the decision.
func (t *GRPC) InvokeAsync(ctx context.Context, req AsyncRequest) error {
	provider, err := providerFromName(req.Provider)
	if err != nil {
		return err
	}
	_, err = t.invoke(ctx, req.Target, provider, req.Raw)
	return err
}

func (t *GRPC) invoke(ctx context.Context, target config.CompiledTarget, provider agentdv1.Provider, raw []byte) (*agentdv1.InvokeResponse, error) {
	timeout := target.Timeout
	if timeout <= 0 {
		timeout = defaultGRPCTimeout
	}
	invokeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ireq := &agentdv1.InvokeRequest{
		Provider:       provider,
		RawPayload:     raw,
		InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
		Deadline:       timestamppb.New(time.Now().Add(timeout)),
	}

	if t.InvokePeer != nil {
		return t.InvokePeer(invokeCtx, target.Endpoint, ireq)
	}
	return dialAndInvoke(invokeCtx, target.Endpoint, ireq)
}

func dialAndInvoke(ctx context.Context, endpoint string, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	conn, err := grpc.NewClient("passthrough:///forward",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return transport.DialEndpoint(ctx, endpoint)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc target dial: %w", err)
	}
	defer conn.Close()
	cli := agentdv1.NewHookServiceClient(conn)
	resp, err := cli.Invoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc target invoke: %w", err)
	}
	return resp, nil
}

func decisionFromProto(d *agentdv1.Decision) agenthooks.Decision {
	if d == nil {
		return agenthooks.NoDecision()
	}
	switch d.GetKind() {
	case agentdv1.DecisionKind_DECISION_KIND_DENY:
		out := agenthooks.Deny(d.GetReason())
		if msg := d.GetSystemMessage(); msg != "" {
			out = out.WithSystemMessage(msg)
		}
		if c := d.GetContext(); c != "" {
			out = out.WithContext(c)
		}
		return out
	case agentdv1.DecisionKind_DECISION_KIND_ASK:
		out := agenthooks.AskUser(d.GetReason())
		if msg := d.GetSystemMessage(); msg != "" {
			out = out.WithSystemMessage(msg)
		}
		if c := d.GetContext(); c != "" {
			out = out.WithContext(c)
		}
		return out
	case agentdv1.DecisionKind_DECISION_KIND_ALLOW:
		out := agenthooks.Allow()
		if msg := d.GetSystemMessage(); msg != "" {
			out = out.WithSystemMessage(msg)
		}
		if c := d.GetContext(); c != "" {
			out = out.WithContext(c)
		}
		return out
	case agentdv1.DecisionKind_DECISION_KIND_BLOCK_PROMPT:
		return agenthooks.BlockPrompt(d.GetReason())
	default:
		return agenthooks.NoDecision()
	}
}

func providerFromName(name string) (agentdv1.Provider, error) {
	switch name {
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
	default:
		return 0, fmt.Errorf("unknown provider %q", name)
	}
}
