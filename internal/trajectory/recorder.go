package trajectory

import (
	"log/slog"
	"time"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

// Recorder enqueues trajectory events from HookService.Invoke (async side).
type Recorder struct {
	queue *Queue
	hub   *Hub
	log   *slog.Logger
}

// NewRecorder wires store + persist + hub into a bounded queue.
func NewRecorder(sessionsDir string, capacity int, log *slog.Logger) *Recorder {
	store := NewStore()
	persist := NewPersister(sessionsDir, log)
	hub := NewHub(log)
	q := NewQueue(capacity, store, persist, hub, log)
	return &Recorder{queue: q, hub: hub, log: log}
}

// Hub returns the live event fan-out registry.
func (r *Recorder) Hub() *Hub {
	if r == nil {
		return nil
	}
	return r.hub
}

// Queue returns the underlying queue (Status drop counter).
func (r *Recorder) Queue() *Queue {
	if r == nil {
		return nil
	}
	return r.queue
}

// Close drains workers, closes hub subscriptions, and flushes pending JSONL.
func (r *Recorder) Close(timeout time.Duration) {
	if r == nil {
		return
	}
	if r.hub != nil {
		r.hub.Close()
	}
	if r.queue == nil {
		return
	}
	r.queue.Close(timeout)
}

// RecordInput is one Invoke worth of trajectory data.
type RecordInput struct {
	Provider       agentdv1.Provider
	InvocationMode agentdv1.InvocationMode
	CWD            string
	ProjectRoot    string
	RawPayload     []byte
	Result         dispatch.InvokeResult
	Snap           *config.Snapshot
}

// Record enqueues ledger events when trajectory is enabled in snap.
func (r *Recorder) Record(in RecordInput) {
	if r == nil || r.queue == nil || in.Snap == nil || !in.Snap.Trajectory.Enabled {
		return
	}
	cfg := in.Snap.Trajectory
	meta := in.Result.Meta
	key := ResolveSessionKey(meta.Provider, meta.SessionID, in.ProjectRoot, in.CWD)
	mode := InvocationModeString(in.InvocationMode)
	now := time.Now().UTC()
	raw := PrepareRaw(in.RawPayload, cfg)

	events := []Event{
		{
			Type:           TypeSessionOpen,
			Source:         SourceSystem,
			TS:             now,
			Provider:       key.Provider,
			InvocationMode: mode,
			SessionID:      key.SessionID,
			ProjectRoot:    key.ProjectRoot,
			CWD:            in.CWD,
			Data: mustJSON(SessionOpenData{
				Provider:    key.Provider,
				CWD:         in.CWD,
				ProjectRoot: in.ProjectRoot,
			}),
		},
		{
			Type:           TypeHookInvoked,
			Source:         SourceHook,
			TS:             now,
			Provider:       key.Provider,
			InvocationMode: mode,
			SessionID:      key.SessionID,
			ProjectRoot:    key.ProjectRoot,
			CWD:            in.CWD,
			Data: mustJSON(HookInvokedData{
				Kind:      meta.EventKind,
				ToolName:  meta.ToolName,
				ToolUseID: meta.ToolUseID,
				HasRoute:  meta.HasRoute,
			}),
			Raw: raw,
		},
	}

	decision := agentdv1.DecisionKind_DECISION_KIND_NO_DECISION
	reason := ""
	if d := in.Result.Decision; d != nil {
		decision = d.GetKind()
		reason = d.GetReason()
	}
	events = append(events, Event{
		Type:           TypeHookDecided,
		Source:         SourceDecision,
		TS:             now,
		Provider:       key.Provider,
		InvocationMode: mode,
		SessionID:      key.SessionID,
		ProjectRoot:    key.ProjectRoot,
		CWD:            in.CWD,
		Data: mustJSON(HookDecidedData{
			Kind:                 meta.EventKind,
			Decision:             decision.String(),
			Reason:               reason,
			ConfigGeneration:     in.Snap.Generation,
			ConfigFingerprint:    in.Snap.Fingerprint,
			AsyncDispatchedCount: in.Result.AsyncDispatchedCount,
		}),
	})

	if in.Result.AsyncDispatchedCount > 0 {
		events = append(events, Event{
			Type:           TypeAsyncDispatched,
			Source:         SourceHook,
			TS:             now,
			Provider:       key.Provider,
			InvocationMode: mode,
			SessionID:      key.SessionID,
			ProjectRoot:    key.ProjectRoot,
			CWD:            in.CWD,
			Data:           mustJSON(AsyncDispatchedData{Count: in.Result.AsyncDispatchedCount}),
		})
	}

	if !r.queue.Enqueue(key, events) {
		dropKey := key
		r.queue.Enqueue(dropKey, []Event{{
			Type:      TypeAsyncDropped,
			Source:    SourceSystem,
			TS:        time.Now().UTC(),
			Provider:  key.Provider,
			SessionID: key.SessionID,
			Ignorable: true,
			Data:      mustJSON(AsyncDroppedData{Reason: "trajectory_queue_overflow"}),
		}})
	}
}

// InvocationModeString maps proto enum to ledger strings.
func InvocationModeString(m agentdv1.InvocationMode) string {
	switch m {
	case agentdv1.InvocationMode_INVOCATION_MODE_STDIN:
		return "stdin"
	case agentdv1.InvocationMode_INVOCATION_MODE_ARGV:
		return "argv"
	case agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY:
		return "notify"
	default:
		return ""
	}
}
