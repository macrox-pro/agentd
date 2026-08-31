package install

import (
	"context"
	"fmt"
	"io"

	ahinstall "github.com/speakeasy-api/agenthooks/install"
)

// HookStatus classifies hook config freshness for one target.
type HookStatus string

const (
	HookStatusMissing HookStatus = "missing"
	HookStatusCurrent HookStatus = "current"
	HookStatusStale   HookStatus = "stale"
)

// PlanEntry is the planned outcome for one install target.
type PlanEntry struct {
	Target  Target
	Changes []FileChange
	Status  HookStatus
}

// Plan diffs manifest against each target without writing.
func Plan(ctx context.Context, targets []Target, command []string) ([]PlanEntry, error) {
	m, err := Manifest(command)
	if err != nil {
		return nil, err
	}
	out := make([]PlanEntry, 0, len(targets))
	for _, t := range targets {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		ahTarget, _, err := toAHTargetWithLabel(t)
		if err != nil {
			return nil, err
		}
		changes, err := ahinstall.Diff(m, ahTarget)
		if err != nil {
			return nil, err
		}
		mapped := mapChanges(changes)
		out = append(out, PlanEntry{
			Target:  t,
			Changes: mapped,
			Status:  hookStatusFromChanges(mapped),
		})
	}
	return out, nil
}

// RunAll installs hooks for each target.
func RunAll(ctx context.Context, targets []Target, command []string, dryRun bool) ([]Result, error) {
	m, err := Manifest(command)
	if err != nil {
		return nil, err
	}
	var installOpts []ahinstall.InstallOption
	if dryRun {
		installOpts = append(installOpts, ahinstall.WithDryRun())
	}
	out := make([]Result, 0, len(targets))
	for _, t := range targets {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		ahTarget, scopeLabel, err := toAHTargetWithLabel(t)
		if err != nil {
			return nil, err
		}
		changes, err := ahinstall.Diff(m, ahTarget)
		if err != nil {
			return nil, err
		}
		if err := ahinstall.Install(ctx, m, ahTarget, installOpts...); err != nil {
			return nil, err
		}
		out = append(out, Result{
			Provider: t.Provider.String(),
			Scope:    scopeLabel,
			Dir:      t.Dir,
			Changes:  mapChanges(changes),
		})
	}
	return out, nil
}

// WriteReports prints summaries for multiple install results.
func WriteReports(w io.Writer, results []Result) error {
	for i, r := range results {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if err := WriteReport(w, r); err != nil {
			return err
		}
	}
	return nil
}

func toAHTargetWithLabel(t Target) (ahinstall.Target, string, error) {
	ahProv, err := t.Provider.Agenthooks()
	if err != nil {
		return ahinstall.Target{}, "", err
	}
	scope, scopeLabel, err := resolveAHScope(t.Scope)
	if err != nil {
		return ahinstall.Target{}, "", err
	}
	return ahinstall.Target{
		Provider: ahProv,
		Scope:    scope,
		Dir:      t.Dir,
	}, scopeLabel, nil
}

func hookStatusFromChanges(changes []FileChange) HookStatus {
	if len(changes) == 0 {
		return HookStatusMissing
	}
	allUnchanged := true
	allCreate := true
	anyCreate := false
	for _, c := range changes {
		if c.State != StateUnchanged {
			allUnchanged = false
		}
		if c.State != StateCreate {
			allCreate = false
		}
		if c.State == StateCreate {
			anyCreate = true
		}
	}
	if allUnchanged {
		return HookStatusCurrent
	}
	if allCreate {
		return HookStatusMissing
	}
	if !anyCreate {
		return HookStatusCurrent
	}
	return HookStatusStale
}
