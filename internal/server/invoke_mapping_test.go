package server_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/decision"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
)

type fakeSnapshotSource struct {
	snap *config.Snapshot
}

func (f fakeSnapshotSource) SnapshotFor(_, _ string) *config.Snapshot {
	return f.snap
}

type fakeInvoker struct {
	result dispatch.InvokeResult
	err    error
}

func (f fakeInvoker) Invoke(context.Context, dispatch.InvokeInput) (dispatch.InvokeResult, error) {
	return f.result, f.err
}

type recordingInvoker struct {
	mu   sync.Mutex
	last dispatch.InvokeInput
}

func (r *recordingInvoker) Invoke(_ context.Context, in dispatch.InvokeInput) (dispatch.InvokeResult, error) {
	r.mu.Lock()
	r.last = in
	r.mu.Unlock()
	return dispatch.InvokeResult{Decision: decision.Neutral()}, nil
}

func (r *recordingInvoker) mode() agentdv1.InvocationMode {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.last.InvocationMode
}

func TestHookServiceInvokeMapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseSnap := &config.Snapshot{Generation: 42, Fingerprint: "fp-test"}

	tests := []struct {
		name       string
		snap       server.SnapshotSource
		inv        server.Invoker
		wantKind   agentdv1.DecisionKind
		wantGen    uint64
		wantFP     string
		wantAsync  uint32
		wantReason string
	}{
		{
			name:     "snap_nil_neutral",
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
			wantGen:  0,
		},
		{
			name:     "engine_nil_neutral",
			snap:     fakeSnapshotSource{snap: baseSnap},
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
			wantGen:  42,
			wantFP:   "fp-test",
		},
		{
			name:     "engine_error_neutral",
			snap:     fakeSnapshotSource{snap: baseSnap},
			inv:      fakeInvoker{err: assert.AnError},
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
			wantGen:  42,
		},
		{
			name: "engine_success_maps_decision",
			snap: fakeSnapshotSource{snap: baseSnap},
			inv: fakeInvoker{result: dispatch.InvokeResult{
				Decision: &agentdv1.Decision{
					Kind:   agentdv1.DecisionKind_DECISION_KIND_ASK,
					Reason: "confirm shell",
				},
			}},
			wantKind:   agentdv1.DecisionKind_DECISION_KIND_ASK,
			wantGen:    42,
			wantReason: "confirm shell",
		},
		{
			name:     "snap_generation_in_response",
			snap:     fakeSnapshotSource{snap: &config.Snapshot{Generation: 99, Fingerprint: "proj-fp"}},
			inv:      fakeInvoker{result: dispatch.InvokeResult{Decision: decision.Neutral()}},
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
			wantGen:  99,
			wantFP:   "proj-fp",
		},
		{
			name: "async_count_forwarded",
			snap: fakeSnapshotSource{snap: baseSnap},
			inv: fakeInvoker{result: dispatch.InvokeResult{
				Decision:             decision.Neutral(),
				AsyncDispatchedCount: 3,
			}},
			wantKind:  agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
			wantGen:   42,
			wantAsync: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			conn := dialHook(t, server.NewHookService(tt.snap, tt.inv, nil, nil))
			hook := agentdv1.NewHookServiceClient(conn)

			resp, err := hook.Invoke(ctx, &agentdv1.InvokeRequest{
				Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
				RawPayload: []byte(`{}`),
			})
			require.NoError(t, err, "Invoke(%q)", tt.name)
			require.NotNil(t, resp, "Invoke(%q)", tt.name)
			assert.Equal(t, tt.wantKind, resp.GetDecision().GetKind(), "Invoke(%q) kind", tt.name)
			assert.Equal(t, tt.wantGen, resp.GetConfig().GetGeneration(), "Invoke(%q) generation", tt.name)
			if tt.wantFP != "" {
				assert.Equal(t, tt.wantFP, resp.GetConfig().GetFingerprint(), "Invoke(%q) fingerprint", tt.name)
			}
			if tt.wantReason != "" {
				assert.Equal(t, tt.wantReason, resp.GetDecision().GetReason(), "Invoke(%q) reason", tt.name)
			}
			if tt.wantAsync != 0 {
				assert.Equal(t, tt.wantAsync, resp.GetAsyncDispatchedCount(), "Invoke(%q) async", tt.name)
			}
		})
	}
}

func TestHookServiceInvocationMode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseSnap := &config.Snapshot{Generation: 1, Fingerprint: "fp"}

	tests := []struct {
		name     string
		provider agentdv1.Provider
		in       agentdv1.InvocationMode
		want     agentdv1.InvocationMode
	}{
		{
			name:     "cursor unspecified to argv",
			provider: agentdv1.Provider_PROVIDER_CURSOR,
			in:       agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
			want:     agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
		},
		{
			name:     "cursor argv preserved",
			provider: agentdv1.Provider_PROVIDER_CURSOR,
			in:       agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
			want:     agentdv1.InvocationMode_INVOCATION_MODE_ARGV,
		},
		{
			name:     "claude unspecified stays unspecified",
			provider: agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			in:       agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
			want:     agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
		},
		{
			name:     "codex unspecified to stdin",
			provider: agentdv1.Provider_PROVIDER_CODEX,
			in:       agentdv1.InvocationMode_INVOCATION_MODE_UNSPECIFIED,
			want:     agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rec := &recordingInvoker{}
			conn := dialHook(t, server.NewHookService(fakeSnapshotSource{snap: baseSnap}, rec, nil, nil))
			hook := agentdv1.NewHookServiceClient(conn)

			_, err := hook.Invoke(ctx, &agentdv1.InvokeRequest{
				Provider:       tt.provider,
				RawPayload:     []byte(`{}`),
				InvocationMode: tt.in,
			})
			require.NoError(t, err, "Invoke(%q)", tt.name)
			assert.Equal(t, tt.want, rec.mode(), "Invoke(%q) mode", tt.name)
		})
	}
}
