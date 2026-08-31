package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"

	"github.com/macrox-pro/agentd/internal/install"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// Mode selects wizard depth (short install vs full setup).
type Mode int

const (
	ModeShort Mode = iota
	ModeFull
)

// Deps supplies plan/apply callbacks for the wizard.
type Deps struct {
	Plan  func(context.Context, []install.Target) ([]install.PlanEntry, error)
	Apply func(context.Context, []install.Target) ([]install.Result, error)
}

// WizardOptions configures RunWizard.
type WizardOptions struct {
	Mode   Mode
	Yes    bool
	DryRun bool
	Env    install.DiscoverEnv
	Out    io.Writer
}

// RunWizard guides target selection, preview, and optional apply.
func RunWizard(ctx context.Context, findings []install.Finding, deps Deps, opts WizardOptions) error {
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}
	high := highConfidence(findings)
	if len(high) == 0 {
		_, err := fmt.Fprintln(out, "no high-confidence agents detected; run agentd doctor")
		return err
	}

	choices := make([]huh.Option[string], len(high))
	byKey := make(map[string]install.Finding, len(high))
	defaultKeys := make([]string, len(high))
	for i, f := range high {
		key := findingKey(f)
		defaultKeys[i] = key
		byKey[key] = f
		choices[i] = huh.NewOption(findingLabel(f), key)
	}
	selected := append([]string(nil), defaultKeys...)

	selectTitle := "Select agents to configure"
	if opts.Mode == ModeFull {
		selectTitle = titleStyle.Render("agentd setup") + "\n" + selectTitle
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(selectTitle).
				Description(mutedStyle.Render("High-confidence detections only")).
				Options(choices...).
				Value(&selected),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	if len(selected) == 0 {
		_, err := fmt.Fprintln(out, "no agents selected")
		return err
	}

	targets, err := targetsForKeys(selected, byKey, opts.Env)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		_, err := fmt.Fprintln(out, "no install targets for selection")
		return err
	}

	entries, err := deps.Plan(ctx, targets)
	if err != nil {
		return err
	}
	if err := writePlanPreview(out, entries); err != nil {
		return err
	}
	if opts.DryRun {
		return nil
	}

	apply := opts.Yes
	if !apply {
		confirmForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Apply hook installs?").
					Affirmative("Yes").
					Negative("No").
					Value(&apply),
			),
		)
		if err := confirmForm.Run(); err != nil {
			return err
		}
	}
	if !apply {
		return nil
	}

	results, err := deps.Apply(ctx, targets)
	if err != nil {
		return err
	}
	return install.WriteReports(out, results)
}

func highConfidence(findings []install.Finding) []install.Finding {
	var out []install.Finding
	for _, f := range findings {
		if f.Confidence == install.ConfidenceHigh {
			out = append(out, f)
		}
	}
	return out
}

func findingKey(f install.Finding) string {
	return f.Provider.String()
}

func findingLabel(f install.Finding) string {
	var parts []string
	if f.ProjectDir != "" {
		parts = append(parts, "project")
	}
	if f.UserDir != "" {
		parts = append(parts, "user")
	}
	suffix := ""
	if len(parts) > 0 {
		suffix = " (" + strings.Join(parts, "+") + ")"
	}
	return f.Provider.String() + suffix
}

func targetsForKeys(keys []string, byKey map[string]install.Finding, env install.DiscoverEnv) ([]install.Target, error) {
	var out []install.Target
	for _, key := range keys {
		f, ok := byKey[key]
		if !ok {
			continue
		}
		ts, err := install.TargetsFromHighConfidence([]install.Finding{f}, env)
		if err != nil {
			return nil, err
		}
		out = append(out, ts...)
	}
	return out, nil
}

func writePlanPreview(w io.Writer, entries []install.PlanEntry) error {
	if _, err := fmt.Fprintln(w, "\nPlanned changes:"); err != nil {
		return err
	}
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "provider=%s scope=%s dir=%s hooks=%s\n",
			e.Target.Provider, e.Target.Scope, e.Target.Dir, e.Status); err != nil {
			return err
		}
		for _, c := range e.Changes {
			if _, err := fmt.Fprintf(w, "  %-9s %s\n", c.State, c.Path); err != nil {
				return err
			}
		}
	}
	return nil
}
