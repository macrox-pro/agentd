package trajectory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

const replayInvokeTimeout = 30 * time.Second

// ErrReplayNoRaw is returned when no hook/invoked events have stored Raw.
var ErrReplayNoRaw = fmt.Errorf("policy replay requires trajectory.include_raw=true at record time")

// ReplayOptions configures an offline policy dry-run against stored Raw payloads.
type ReplayOptions struct {
	SessionsRoot string
	Provider     string
	SessionID    string
	Seq          uint64 // 0 = all hook/invoked with Raw
	Snap         *config.Snapshot
	Engine       *dispatch.Engine
}

// ReplayHit is one replayed hook/invoked event.
type ReplayHit struct {
	Seq             uint64 `json:"seq"`
	Kind            string `json:"kind,omitempty"`
	StoredDecision  string `json:"stored_decision,omitempty"`
	ReplayDecision  string `json:"replay_decision,omitempty"`
	Match           bool   `json:"match"`
	Error           string `json:"error,omitempty"`
}

// ReplayResult is the full policy replay output.
type ReplayResult struct {
	Provider  string      `json:"provider"`
	SessionID string      `json:"session_id"`
	Hits      []ReplayHit `json:"hits"`
}

// ReplayPolicy re-Invokes stored Raw through Engine (offline; no live agent).
func ReplayPolicy(ctx context.Context, opts ReplayOptions) (ReplayResult, error) {
	if opts.Snap == nil {
		return ReplayResult{}, fmt.Errorf("nil config snapshot")
	}
	if opts.Engine == nil {
		return ReplayResult{}, fmt.Errorf("nil dispatch engine")
	}
	root := opts.SessionsRoot
	if root == "" {
		root = DefaultSessionsDir()
	}
	path, err := FindSessionPath(root, opts.Provider, opts.SessionID)
	if err != nil {
		return ReplayResult{}, err
	}
	events, err := ReadEvents(path)
	if err != nil {
		return ReplayResult{}, fmt.Errorf("read session: %w", err)
	}

	prov := CanonicalProvider(opts.Provider)
	out := ReplayResult{Provider: prov, SessionID: opts.SessionID}
	decidedByPrev := map[uint64]string{}
	for _, e := range events {
		if e.Type != TypeHookDecided {
			continue
		}
		var d HookDecidedData
		_ = json.Unmarshal(e.Data, &d)
		// Pair with the preceding hook/invoked seq (same append batch uses adjacent seqs).
		if e.Seq > 1 {
			decidedByPrev[e.Seq-1] = d.Decision
		}
	}

	var candidates int
	for _, e := range events {
		if e.Type != TypeHookInvoked {
			continue
		}
		if opts.Seq != 0 && e.Seq != opts.Seq {
			continue
		}
		candidates++
		if len(e.Raw) == 0 {
			continue
		}
		hit := ReplayHit{
			Seq:            e.Seq,
			StoredDecision: decidedByPrev[e.Seq],
		}
		var inv HookInvokedData
		_ = json.Unmarshal(e.Data, &inv)
		hit.Kind = inv.Kind

		protoProv, err := providerFromName(e.Provider)
		if err != nil {
			hit.Error = err.Error()
			out.Hits = append(out.Hits, hit)
			continue
		}
		mode := invocationModeFromString(e.InvocationMode, e.Provider)
		res, err := opts.Engine.Invoke(ctx, dispatch.InvokeInput{
			Provider:       protoProv,
			RawPayload:     append([]byte(nil), e.Raw...),
			InvocationMode: mode,
			Snap:           opts.Snap,
			CWD:            e.CWD,
			ProjectRoot:    e.ProjectRoot,
			Deadline:       time.Now().Add(replayInvokeTimeout),
		})
		if err != nil {
			hit.Error = err.Error()
			out.Hits = append(out.Hits, hit)
			continue
		}
		replayKind := agentdv1.DecisionKind_DECISION_KIND_NO_DECISION.String()
		if res.Decision != nil {
			replayKind = res.Decision.GetKind().String()
		}
		hit.ReplayDecision = replayKind
		hit.Match = hit.StoredDecision == "" || hit.StoredDecision == replayKind
		out.Hits = append(out.Hits, hit)
	}

	if candidates > 0 && len(out.Hits) == 0 {
		return out, ErrReplayNoRaw
	}
	if opts.Seq != 0 && len(out.Hits) == 0 {
		return out, fmt.Errorf("no hook/invoked event with seq=%d and raw payload", opts.Seq)
	}
	if candidates == 0 {
		return out, fmt.Errorf("no hook/invoked events to replay")
	}
	return out, nil
}

func providerFromName(name string) (agentdv1.Provider, error) {
	switch CanonicalProvider(name) {
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
	case "kimi-code":
		return agentdv1.Provider_PROVIDER_KIMI_CODE, nil
	default:
		return 0, fmt.Errorf("unknown provider %q", name)
	}
}

func invocationModeFromString(mode, provider string) agentdv1.InvocationMode {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "argv":
		return agentdv1.InvocationMode_INVOCATION_MODE_ARGV
	case "notify":
		return agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY
	case "stdin":
		return agentdv1.InvocationMode_INVOCATION_MODE_STDIN
	default:
		if CanonicalProvider(provider) == "cursor" {
			return agentdv1.InvocationMode_INVOCATION_MODE_ARGV
		}
		return agentdv1.InvocationMode_INVOCATION_MODE_STDIN
	}
}
