package config

import (
	"fmt"
	"sort"
)

func baseFileConfig() *fileConfig {
	trueVal := true
	dd := map[string]fileKindDefault{}
	for k, v := range defaultKindDefaults() {
		blocking := v.Blocking
		dd[k] = fileKindDefault{
			Mode:     string(v.Mode),
			Blocking: &blocking,
		}
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
	}
}

// Compile merges defaults with optional user fileConfig and produces Snapshot fields.
func Compile(user *fileConfig) (Policy, AsyncConfig, Guards, []CompiledRoute, error) {
	merged := mergeFile(baseFileConfig(), user)

	pol, err := parsePolicy(merged.Policy, defaultPolicy())
	if err != nil {
		return Policy{}, AsyncConfig{}, Guards{}, nil, err
	}
	async, err := parseAsync(merged.Async, defaultAsync())
	if err != nil {
		return Policy{}, AsyncConfig{}, Guards{}, nil, err
	}
	guards, err := parseGuards(merged.Guards, defaultGuards())
	if err != nil {
		return Policy{}, AsyncConfig{}, Guards{}, nil, err
	}
	kinds, err := parseKindDefaults(merged.DispatchDefaults, defaultKindDefaults())
	if err != nil {
		return Policy{}, AsyncConfig{}, Guards{}, nil, err
	}
	routes := compileRoutes(kinds, guards)
	return pol, async, guards, routes, nil
}

func compileRoutes(kinds map[string]KindDefault, guards Guards) []CompiledRoute {
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
			Name: fmt.Sprintf("default-%s", kind),
			Kind: kind,
			Mode: mode,
		}
		switch mode {
		case ModeSyncOnly, ModeParallel, ModeAfterSync:
			if guards.Secrets.Enabled {
				r.Sync = []CompiledTarget{{
					Kind:   TargetBuiltin,
					Guards: []string{"secrets"},
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
