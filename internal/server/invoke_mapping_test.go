package server

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
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

func dialHookService(t *testing.T, h *hookService) *grpc.ClientConn {
	t.Helper()
	s := grpc.NewServer()
	agentdv1.RegisterHookServiceServer(s, h)
	lis := bufconn.Listen(1024 * 1024)
	go func() { _ = s.Serve(lis) }()
	t.Cleanup(func() {
		s.Stop()
		_ = lis.Close()
	})
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "dial bufconn")
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

func TestHookServiceInvokeMapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseSnap := &config.Snapshot{Generation: 42, Fingerprint: "fp-test"}

	tests := []struct {
		name       string
		snap       SnapshotSource
		inv        Invoker
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
			inv:      fakeInvoker{result: dispatch.InvokeResult{Decision: dispatch.NeutralDecision()}},
			wantKind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
			wantGen:  99,
			wantFP:   "proj-fp",
		},
		{
			name: "async_count_forwarded",
			snap: fakeSnapshotSource{snap: baseSnap},
			inv: fakeInvoker{result: dispatch.InvokeResult{
				Decision:             dispatch.NeutralDecision(),
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
			conn := dialHookService(t, &hookService{
				snap:   tt.snap,
				engine: tt.inv,
			})
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
