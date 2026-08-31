package config

import (
	"fmt"
	"sort"
	"time"
)

func baseFileConfig() *fileConfig {
	trueVal := true
	dd := map[string]fileKindDefault{}
	for k, v := range defaultKindDefaults() {
		dd[k] = fileKindDefault{Mode: string(v.Mode)}
	}
	async := defaultAsync()
	pol := defaultPolicy()
	return &fileConfig{
		Version: 1,
		Policy: &filePolicy{
			Fail:        string(pol.Fail),
			Unsupported: string(pol.Unsupported),
			AskFallback: string(pol.AskFallback),
			Offline:     string(pol.Offline),
		},
		Async: &fileAsync{
			QueueCapacity: async.QueueCapacity,
			WorkerLimit:   async.WorkerLimit,
			TargetTimeout: async.TargetTimeout.String(),
			OnOverflow:    string(async.OnOverflow),
		},
		Guards: &fileGuards{
			Secrets: &fileSecretsGuard{
				Enabled: &trueVal,
				Action:  string(GuardAsk),
				Rules:   append([]string(nil), DefaultSecretsRules...),
			},
		},
		DispatchDefaults: dd,
		Logging: &fileLogging{
			Level: string(LogLevelInfo),
		},
	}
}

// CompileResult is the compiled Snapshot fields from a layer merge.
type CompileResult struct {
	Policy          Policy
	Async           AsyncConfig
	Guards          Guards
	Approvals       Approvals
	TemporaryBlocks []TemporaryBlock
	Trajectory      TrajectoryConfig
	Logging         LoggingConfig
	Metrics         MetricsConfig
	Routes          []CompiledRoute
	Merged          *fileConfig
}

// CompileMerged merges defaults ⊕ user ⊕ project ⊕ runtime and compiles Snapshot fields.
func CompileMerged(user, project, runtime *fileConfig) (CompileResult, error) {
	merged := mergeFile(baseFileConfig(), user)
	merged = mergeFile(merged, project)
	merged = mergeFile(merged, runtime)

	now := time.Now().UTC()
	pol, err := parsePolicy(merged.Policy, defaultPolicy())
	if err != nil {
		return CompileResult{}, err
	}
	async, err := parseAsync(merged.Async, defaultAsync())
	if err != nil {
		return CompileResult{}, err
	}
	guards, err := parseGuards(merged.Guards, defaultGuards())
	if err != nil {
		return CompileResult{}, err
	}
	approvals, err := parseApprovals(merged.Approvals, now)
	if err != nil {
		return CompileResult{}, err
	}
	blocks, err := parseTemporaryBlocks(merged.Blocks, now)
	if err != nil {
		return CompileResult{}, err
	}
	kinds, err := parseKindDefaults(merged.DispatchDefaults, defaultKindDefaults())
	if err != nil {
		return CompileResult{}, err
	}
	userRoutes, err := compileUserRoutes(merged.Dispatch)
	if err != nil {
		return CompileResult{}, err
	}
	defaults := compileDefaultRoutes(kinds, guards)
	routes := append(userRoutes, defaults...)
	traj, err := parseTrajectory(merged.Trajectory, defaultTrajectory())
	if err != nil {
		return CompileResult{}, err
	}
	logging, err := parseLogging(merged.Logging, defaultLogging())
	if err != nil {
		return CompileResult{}, err
	}
	metricsCfg, err := parseMetrics(merged.Metrics, defaultMetrics())
	if err != nil {
		return CompileResult{}, err
	}
	return CompileResult{
		Policy:          pol,
		Async:           async,
		Guards:          guards,
		Approvals:       approvals,
		TemporaryBlocks: blocks,
		Trajectory:      traj,
		Logging:         logging,
		Metrics:         metricsCfg,
		Routes:          routes,
		Merged:          merged,
	}, nil
}

func compileUserRoutes(in []fileRoute) ([]CompiledRoute, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make([]CompiledRoute, 0, len(in))
	seen := map[string]bool{}
	for i, fr := range in {
		if fr.Name == "" {
			return nil, fmt.Errorf("dispatch[%d]: name is required", i)
		}
		if seen[fr.Name] {
			return nil, fmt.Errorf("dispatch: duplicate route name %q", fr.Name)
		}
		seen[fr.Name] = true
		mode, err := parseDispatchMode(fr.Mode)
		if err != nil {
			return nil, fmt.Errorf("dispatch[%q].mode: %w", fr.Name, err)
		}
		mode = NormalizeMode(mode)
		syncTargets, err := compileTargets(fr.Sync, true, fr.Name)
		if err != nil {
			return nil, err
		}
		asyncTargets, err := compileTargets(fr.Async, false, fr.Name)
		if err != nil {
			return nil, err
		}
		var syncTimeout time.Duration
		if fr.SyncTimeout != "" {
			d, err := time.ParseDuration(fr.SyncTimeout)
			if err != nil {
				return nil, fmt.Errorf("dispatch[%q].sync_timeout: %w", fr.Name, err)
			}
			if d <= 0 {
				return nil, fmt.Errorf("dispatch[%q].sync_timeout must be > 0", fr.Name)
			}
			syncTimeout = d
		}
		out = append(out, CompiledRoute{
			Name:        fr.Name,
			Match:       RouteMatch{Kinds: append([]string(nil), fr.Match.Kind...), Providers: append([]string(nil), fr.Match.Provider...), Tools: append([]string(nil), fr.Match.Tools...)},
			Mode:        mode,
			SyncTimeout: syncTimeout,
			Sync:        syncTargets,
			Async:       asyncTargets,
		})
	}
	return out, nil
}

func compileTargets(in []fileTarget, sync bool, routeName string) ([]CompiledTarget, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make([]CompiledTarget, 0, len(in))
	for i, ft := range in {
		ct, err := compileOneTarget(ft, sync)
		if err != nil {
			return nil, fmt.Errorf("dispatch[%q] target[%d]: %w", routeName, i, err)
		}
		out = append(out, ct)
	}
	return out, nil
}

func compileOneTarget(ft fileTarget, sync bool) (CompiledTarget, error) {
	kind := TargetKind(ft.Target)
	switch kind {
	case TargetBuiltin, TargetExec, TargetHTTP, TargetLog, TargetFile, TargetGRPC:
	case "":
		return CompiledTarget{}, fmt.Errorf("target is required")
	default:
		return CompiledTarget{}, fmt.Errorf("unknown target %q", ft.Target)
	}
	if sync && kind != TargetBuiltin && kind != TargetGRPC {
		return CompiledTarget{}, fmt.Errorf("sync target %q not supported (builtin or grpc only)", kind)
	}
	ct := CompiledTarget{
		Kind:     kind,
		Guards:   append([]string(nil), ft.Guards...),
		Observe:  ft.Observe,
		URL:      ft.URL,
		Command:  append([]string(nil), ft.Command...),
		Stdin:    ft.Stdin,
		Level:    ft.Level,
		Path:     ft.Path,
		Endpoint: ft.Endpoint,
	}
	if ft.Retry != nil {
		ct.Retry = *ft.Retry
	}
	if ft.Timeout != "" {
		d, err := time.ParseDuration(ft.Timeout)
		if err != nil {
			return CompiledTarget{}, fmt.Errorf("timeout: %w", err)
		}
		ct.Timeout = d
	}
	switch kind {
	case TargetBuiltin:
		for _, g := range ct.Guards {
			if !knownGuardName(g) {
				return CompiledTarget{}, fmt.Errorf("unknown guard %q", g)
			}
		}
		if !sync && !ct.Observe && len(ct.Guards) == 0 {
			ct.Observe = true
		}
	case TargetHTTP:
		if ct.URL == "" {
			return CompiledTarget{}, fmt.Errorf("http target requires url")
		}
		if ct.Retry != 0 {
			return CompiledTarget{}, fmt.Errorf("http retry must be 0 in M3")
		}
	case TargetExec:
		if len(ct.Command) == 0 {
			return CompiledTarget{}, fmt.Errorf("exec target requires command")
		}
		if ct.Stdin != "" && ct.Stdin != "raw" {
			return CompiledTarget{}, fmt.Errorf("exec stdin must be empty or %q", "raw")
		}
	case TargetFile:
		if ct.Path == "" {
			return CompiledTarget{}, fmt.Errorf("file target requires path")
		}
	case TargetLog:
		if ct.Level == "" {
			ct.Level = "info"
		}
		switch ct.Level {
		case "debug", "info", "warn", "error":
		default:
			return CompiledTarget{}, fmt.Errorf("log level unknown %q", ct.Level)
		}
	case TargetGRPC:
		if ct.Endpoint == "" {
			return CompiledTarget{}, fmt.Errorf("grpc target requires endpoint")
		}
		onErr := FailClosed
		if ft.OnError != "" {
			m, err := parseFailMode(ft.OnError)
			if err != nil {
				return CompiledTarget{}, fmt.Errorf("grpc on_error: %w", err)
			}
			onErr = m
		}
		ct.OnError = onErr
		if ft.Merge != "" {
			if SyncMerge(ft.Merge) != MergeFirstConclusive {
				return CompiledTarget{}, fmt.Errorf("grpc merge unknown %q (only first_conclusive)", ft.Merge)
			}
			ct.Merge = MergeFirstConclusive
		}
	}
	return ct, nil
}

func compileDefaultRoutes(kinds map[string]KindDefault, guards Guards) []CompiledRoute {
	keys := make([]string, 0, len(kinds))
	for k := range kinds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	routes := make([]CompiledRoute, 0, len(keys))
	for _, kind := range keys {
		kd := kinds[kind]
		mode := NormalizeMode(kd.Mode)
		r := CompiledRoute{
			Name:    fmt.Sprintf("default-%s", kind),
			Kind:    kind,
			Match:   RouteMatch{Kinds: []string{kind}},
			Mode:    mode,
			Default: true,
		}
		switch mode {
		case ModeSyncOnly, ModeParallel, ModeAfterSync:
			names := enabledGuardNames(guards)
			if len(names) > 0 {
				r.Sync = []CompiledTarget{{
					Kind:   TargetBuiltin,
					Guards: names,
				}}
			} else {
				r.Sync = []CompiledTarget{{Kind: TargetBuiltin}}
			}
		}
		switch mode {
		case ModeAsyncOnly, ModeParallel, ModeAfterSync:
			r.Async = []CompiledTarget{{
				Kind:    TargetBuiltin,
				Observe: true,
			}}
		}
		routes = append(routes, r)
	}
	return routes
}
