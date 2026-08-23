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
	"github.com/macrox-pro/agentd/internal/provider"
)

const replayInvokeTimeout = 30 * time.Second

// ReplayOptions configures an offline policy dry-run against stored Raw payloads.
type ReplayOptions struct {
	SessionsRoot string
	Provider     string
	SessionID    string
	Seq          uint64 // 0 = all hook/invoked with Raw
	Snap         *config.Snapshot
	Engine       dispatch.Invoker
}

// ReplayHit is one replayed hook/invoked event.
type ReplayHit struct {
	Seq            uint64 `json:"seq"`
	Kind           string `json:"kind,omitempty"`
	StoredDecision string `json:"stored_decision,omitempty"`
	ReplayDecision string `json:"replay_decision,omitempty"`
	Match          bool   `json:"match"`
	Error          string `json:"error,omitempty"`
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
		return ReplayResult{}, ErrNilConfigSnap
	}
	if opts.Engine == nil {
		return ReplayResult{}, ErrNilEngine
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

		protoProv, err := eventProviderProto(e.Provider)
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
		return out, fmt.Errorf("%w: %d", ErrReplaySeqNotFound, opts.Seq)
	}
	if candidates == 0 {
		return out, ErrReplayNoEvents
	}
	return out, nil
}

func eventProviderProto(name string) (agentdv1.Provider, error) {
	id, ok := provider.Lookup(name)
	if !ok {
		return 0, fmt.Errorf("%w %q", provider.ErrUnknownProvider, name)
	}
	return id.Proto()
}

func invocationModeFromString(mode, providerName string) agentdv1.InvocationMode {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "argv":
		return agentdv1.InvocationMode_INVOCATION_MODE_ARGV
	case "notify":
		return agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY
	case "stdin":
		return agentdv1.InvocationMode_INVOCATION_MODE_STDIN
	default:
		if CanonicalProvider(providerName) == string(provider.Cursor) {
			return agentdv1.InvocationMode_INVOCATION_MODE_ARGV
		}
		return agentdv1.InvocationMode_INVOCATION_MODE_STDIN
	}
}
