package dispatch_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

type recordObserver struct {
	mu      sync.Mutex
	invokes []struct {
		provider  string
		eventKind string
		decision  string
		outcome   string
	}
	async []struct {
		targetKind string
		result     string
	}
}

func (o *recordObserver) ObserveInvoke(provider, eventKind, decision, outcome string, seconds float64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.invokes = append(o.invokes, struct {
		provider  string
		eventKind string
		decision  string
		outcome   string
	}{provider, eventKind, decision, outcome})
}

func (o *recordObserver) ObserveAsync(targetKind, result string, seconds float64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.async = append(o.async, struct {
		targetKind string
		result     string
	}{targetKind, result})
}

func (o *recordObserver) lastInvoke() (provider, eventKind, decision, outcome string, ok bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.invokes) == 0 {
		return "", "", "", "", false
	}
	last := o.invokes[len(o.invokes)-1]
	return last.provider, last.eventKind, last.decision, last.outcome, true
}

func TestEngineObserveInvoke(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		run          func(t *testing.T) (obs *recordObserver, wantOutcome, wantDecision string, wantErr bool)
		wantObserved bool
	}{
		{
			name: "ok deny",
			run: func(t *testing.T) (*recordObserver, string, string, bool) {
				obs := &recordObserver{}
				q := dispatch.NewQueue(config.AsyncConfig{QueueCapacity: 8, WorkerLimit: 2, TargetTimeout: time.Second}, nil)
				t.Cleanup(func() { q.Close(2 * time.Second) })
				eng := dispatch.NewEngine(q, nil, obs)
				snap := paritySnap(t, []config.CompiledTarget{{
					Kind:   config.TargetBuiltin,
					Guards: []string{"secrets"},
				}}, nil, config.ModeSyncOnly)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"),
					Snap:       snap,
				})
				require.NoError(t, err, "Invoke")
				return obs, "ok", agentdv1.DecisionKind_DECISION_KIND_ASK.String(), false
			},
			wantObserved: true,
		},
		{
			name: "error",
			run: func(t *testing.T) (*recordObserver, string, string, bool) {
				obs := &recordObserver{}
				eng := dispatch.NewEngine(nil, nil, obs)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       nil,
				})
				require.Error(t, err, "Invoke")
				return obs, "error", agentdv1.DecisionKind_DECISION_KIND_UNSPECIFIED.String(), true
			},
			wantObserved: true,
		},
		{
			name: "timeout",
			run: func(t *testing.T) (*recordObserver, string, string, bool) {
				obs := &recordObserver{}
				dir, err := os.MkdirTemp("/tmp", "agentd-hang-")
				require.NoError(t, err, "MkdirTemp")
				t.Cleanup(func() { _ = os.RemoveAll(dir) })
				sock := filepath.Join(dir, "s.sock")
				ln, err := net.Listen("unix", sock)
				require.NoError(t, err, "Listen")
				t.Cleanup(func() { _ = ln.Close() })
				go func() {
					for {
						c, err := ln.Accept()
						if err != nil {
							return
						}
						go func(conn net.Conn) {
							time.Sleep(time.Minute)
							_ = conn.Close()
						}(c)
					}
				}()
				eng := dispatch.NewEngine(nil, nil, obs)
				snap := paritySnap(t, []config.CompiledTarget{{
					Kind:     config.TargetGRPC,
					Endpoint: "unix://" + sock,
					Timeout:  500 * time.Millisecond,
					OnError:  config.FailClosed,
				}}, nil, config.ModeSyncOnly)
				_, err = eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       snap,
					Deadline:   time.Now().Add(50 * time.Millisecond),
				})
				require.NoError(t, err, "Invoke")
				return obs, "timeout", agentdv1.DecisionKind_DECISION_KIND_DENY.String(), false
			},
			wantObserved: true,
		},
		{
			name: "cancel",
			run: func(t *testing.T) (*recordObserver, string, string, bool) {
				obs := &recordObserver{}
				eng := dispatch.NewEngine(nil, nil, obs)
				snap := paritySnap(t, []config.CompiledTarget{{
					Kind:     config.TargetGRPC,
					Endpoint: "unix://" + filepath.Join(t.TempDir(), "missing.sock"),
					Timeout:  500 * time.Millisecond,
					OnError:  config.FailClosed,
				}}, nil, config.ModeSyncOnly)
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := eng.Invoke(ctx, dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       snap,
				})
				require.Error(t, err, "Invoke")
				return obs, "error", agentdv1.DecisionKind_DECISION_KIND_UNSPECIFIED.String(), true
			},
			wantObserved: true,
		},
		{
			name: "no-route",
			run: func(t *testing.T) (*recordObserver, string, string, bool) {
				obs := &recordObserver{}
				eng := dispatch.NewEngine(nil, nil, obs)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       testSnap(t),
				})
				require.NoError(t, err, "Invoke")
				return obs, "ok", agentdv1.DecisionKind_DECISION_KIND_NO_DECISION.String(), false
			},
			wantObserved: true,
		},
		{
			name: "decode error",
			run: func(t *testing.T) (*recordObserver, string, string, bool) {
				obs := &recordObserver{}
				eng := dispatch.NewEngine(nil, nil, obs)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: []byte(`not-json`),
					Snap:       testSnap(t),
				})
				require.Error(t, err, "Invoke")
				return obs, "error", agentdv1.DecisionKind_DECISION_KIND_UNSPECIFIED.String(), true
			},
			wantObserved: true,
		},
		{
			name: "nil observer no panic",
			run: func(t *testing.T) (*recordObserver, string, string, bool) {
				eng := dispatch.NewEngine(nil, nil, nil)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       testSnap(t),
				})
				require.NoError(t, err, "Invoke")
				return nil, "", "", false
			},
			wantObserved: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			obs, wantOutcome, wantDecision, _ := tt.run(t)
			if !tt.wantObserved {
				return
			}
			_, _, decision, outcome, ok := obs.lastInvoke()
			require.True(t, ok, "ObserveInvoke recorded")
			assert.Equal(t, wantOutcome, outcome, "outcome")
			assert.Equal(t, wantDecision, decision, "decision")
			if outcome == "ok" && wantDecision != agentdv1.DecisionKind_DECISION_KIND_UNSPECIFIED.String() {
				assert.NotEmpty(t, obs.invokes[len(obs.invokes)-1].eventKind, "event_kind")
			}
		})
	}
}

func TestEngineObserveAsync(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		run        func(t *testing.T) *recordObserver
		wantResult string
		wantCount  int
	}{
		{
			name: "http ok",
			run: func(t *testing.T) *recordObserver {
				obs := &recordObserver{}
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				t.Cleanup(srv.Close)
				q := dispatch.NewQueue(config.AsyncConfig{QueueCapacity: 8, WorkerLimit: 2, TargetTimeout: time.Second}, nil)
				t.Cleanup(func() { q.Close(2 * time.Second) })
				eng := dispatch.NewEngine(q, nil, obs)
				snap := paritySnap(t, nil, []config.CompiledTarget{{
					Kind: config.TargetHTTP,
					URL:  srv.URL,
				}}, config.ModeAsyncOnly)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       snap,
				})
				require.NoError(t, err, "Invoke")
				return obs
			},
			wantResult: "ok",
			wantCount:  1,
		},
		{
			name: "http error",
			run: func(t *testing.T) *recordObserver {
				obs := &recordObserver{}
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				t.Cleanup(srv.Close)
				q := dispatch.NewQueue(config.AsyncConfig{QueueCapacity: 8, WorkerLimit: 2, TargetTimeout: time.Second}, nil)
				t.Cleanup(func() { q.Close(2 * time.Second) })
				eng := dispatch.NewEngine(q, nil, obs)
				snap := paritySnap(t, nil, []config.CompiledTarget{{
					Kind: config.TargetHTTP,
					URL:  srv.URL,
				}}, config.ModeAsyncOnly)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       snap,
				})
				require.NoError(t, err, "Invoke")
				return obs
			},
			wantResult: "error",
			wantCount:  1,
		},
		{
			name: "not observed when queue drops",
			run: func(t *testing.T) *recordObserver {
				obs := &recordObserver{}
				block := make(chan struct{})
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					<-block
					w.WriteHeader(http.StatusOK)
				}))
				t.Cleanup(func() {
					close(block)
					srv.Close()
				})
				q := dispatch.NewQueue(config.AsyncConfig{
					QueueCapacity: 1,
					WorkerLimit:   1,
					TargetTimeout: time.Second,
				}, nil)
				t.Cleanup(func() { q.Close(2 * time.Second) })
				eng := dispatch.NewEngine(q, nil, obs)
				snap := paritySnap(t, nil, []config.CompiledTarget{
					{Kind: config.TargetHTTP, URL: srv.URL},
					{Kind: config.TargetHTTP, URL: srv.URL},
				}, config.ModeAsyncOnly)
				_, err := eng.Invoke(context.Background(), dispatch.InvokeInput{
					Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
					RawPayload: claudeToolPre(t, "echo ok"),
					Snap:       snap,
				})
				require.NoError(t, err, "Invoke")
				time.Sleep(100 * time.Millisecond)
				return obs
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := tt.run(t)
			deadline := time.Now().Add(3 * time.Second)
			for {
				obs.mu.Lock()
				n := len(obs.async)
				obs.mu.Unlock()
				if n >= tt.wantCount || time.Now().After(deadline) {
					break
				}
				time.Sleep(20 * time.Millisecond)
			}
			obs.mu.Lock()
			n := len(obs.async)
			var result string
			if n > 0 {
				result = obs.async[n-1].result
			}
			obs.mu.Unlock()
			assert.Equal(t, tt.wantCount, n, "async observations")
			if tt.wantResult != "" {
				assert.Equal(t, tt.wantResult, result, "result")
			}
		})
	}
}
