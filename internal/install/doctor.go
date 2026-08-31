package install

import (
	"context"
	"time"

	"github.com/macrox-pro/agentd/internal/hookclient"
)

const doctorDialTimeout = 2 * time.Second

// DoctorOptions configures Report.
type DoctorOptions struct {
	Cwd     string
	Home    string
	Socket  string // empty skips daemon dial
	Command []string
	Env     DiscoverEnv
}

// DoctorFinding combines discovery with hook plan status.
type DoctorFinding struct {
	Finding
	Plan *PlanEntry
}

// DoctorReport is the read-only doctor output.
type DoctorReport struct {
	Findings        []DoctorFinding
	DaemonReachable bool
}

// Report discovers agents and enriches with hook plan status.
func Report(ctx context.Context, opts DoctorOptions) (DoctorReport, error) {
	env := opts.Env.withDefaults()
	if opts.Cwd != "" {
		env.Cwd = opts.Cwd
	}
	if opts.Home != "" {
		env.Home = opts.Home
	}
	findings, err := Discover(ctx, env)
	if err != nil {
		return DoctorReport{}, err
	}
	targets, err := TargetsFromHighConfidence(findings, env)
	if err != nil {
		return DoctorReport{}, err
	}
	var planEntries []PlanEntry
	if len(opts.Command) > 0 {
		planEntries, err = Plan(ctx, targets, opts.Command)
		if err != nil {
			return DoctorReport{}, err
		}
	}
	planByKey := make(map[string]PlanEntry, len(planEntries))
	for _, pe := range planEntries {
		planByKey[targetKey(pe.Target)] = pe
	}

	out := make([]DoctorFinding, 0, len(findings))
	for _, f := range findings {
		df := DoctorFinding{Finding: f}
		ts, err := targetsFromFinding(f, env)
		if err != nil {
			return DoctorReport{}, err
		}
		if len(ts) > 0 {
			if pe, ok := planByKey[targetKey(ts[0])]; ok {
				peCopy := pe
				df.Plan = &peCopy
			}
		}
		out = append(out, df)
	}

	report := DoctorReport{Findings: out}
	if opts.Socket != "" {
		dialCtx, cancel := context.WithTimeout(ctx, doctorDialTimeout)
		defer cancel()
		cli, err := hookclient.DialReady(dialCtx, opts.Socket)
		if err == nil {
			report.DaemonReachable = true
			_ = cli.Close()
		}
	}
	return report, nil
}

func targetKey(t Target) string {
	return t.Provider.String() + ":" + t.Scope + ":" + t.Dir
}
