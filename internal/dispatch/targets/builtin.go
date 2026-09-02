package targets

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/guard"
)

// Builtin runs in-process agenthooks guards / observers.
type Builtin struct {
	Guards          config.Guards
	Policy          config.Policy
	Approvals       config.Approvals
	TemporaryBlocks []config.TemporaryBlock
	ProjectRoot     string
	Log             *slog.Logger
}

// Decide runs sync guards via Runner.Decide.
func (b *Builtin) Decide(ctx context.Context, typed any, guardNames []string) (agenthooks.Decision, error) {
	r := b.newRunner(guardNames, false)
	return r.Decide(ctx, typed)
}

// Observe runs async OnAny observers; errors are logged only.
func (b *Builtin) Observe(ctx context.Context, typed any) {
	r := b.newRunner(nil, true)
	// Decide still runs OnAny observers; we ignore the decision.
	_, err := r.Decide(ctx, typed)
	if err != nil && b.logger() != nil {
		b.logger().Warn("builtin observe failed", "error", err)
	}
}

func (b *Builtin) newRunner(guardNames []string, observeOnly bool) *agenthooks.Runner {
	r := agenthooks.New(agenthooks.WithLogger(b.logger()))
	if observeOnly {
		r.OnAny(func(context.Context, *agenthooks.Event) error { return nil })
		return r
	}
	dctx := guard.DecisionContext{
		Approvals:       b.Approvals,
		TemporaryBlocks: b.TemporaryBlocks,
		ProjectRoot:     b.ProjectRoot,
		AskFallback:     b.Policy.AskFallback,
	}
	guard.AttachBlocks(r, b.TemporaryBlocks)
	guard.AttachCheckers(r, b.Guards, dctx, guardNames)
	return r
}

func (b *Builtin) logger() *slog.Logger {
	if b.Log != nil {
		return b.Log
	}
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// ProjectRootOf returns the project directory for approval matching.
func ProjectRootOf(snap *config.Snapshot) string {
	if snap == nil || snap.ProjectPath == "" {
		return ""
	}
	return filepath.Dir(snap.ProjectPath)
}
