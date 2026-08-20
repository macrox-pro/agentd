package config

import "maps"

func mergeFile(base *fileConfig, user *fileConfig) *fileConfig {
	if base == nil {
		base = &fileConfig{}
	}
	if user == nil {
		return base
	}
	out := *base
	if user.Version != 0 {
		out.Version = user.Version
	}
	if user.Policy != nil {
		p := mergePolicyPtr(out.Policy, user.Policy)
		out.Policy = p
	}
	if user.Async != nil {
		a := mergeAsyncPtr(out.Async, user.Async)
		out.Async = a
	}
	if user.Guards != nil {
		g := mergeGuardsPtr(out.Guards, user.Guards)
		out.Guards = g
	}
	if len(user.DispatchDefaults) > 0 {
		dd := map[string]fileKindDefault{}
		maps.Copy(dd, out.DispatchDefaults)
		for k, v := range user.DispatchDefaults {
			cur := dd[k]
			if v.Mode != "" {
				cur.Mode = v.Mode
			}
			dd[k] = cur
		}
		out.DispatchDefaults = dd
	}
	if user.Dispatch != nil {
		out.Dispatch = append([]fileRoute(nil), user.Dispatch...)
	}
	return &out
}

func mergePolicyPtr(base, user *filePolicy) *filePolicy {
	if base == nil && user == nil {
		return nil
	}
	out := filePolicy{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Fail != "" {
		out.Fail = user.Fail
	}
	if user.Unsupported != "" {
		out.Unsupported = user.Unsupported
	}
	if user.AskFallback != "" {
		out.AskFallback = user.AskFallback
	}
	if user.Offline != "" {
		out.Offline = user.Offline
	}
	return &out
}

func mergeAsyncPtr(base, user *fileAsync) *fileAsync {
	if base == nil && user == nil {
		return nil
	}
	out := fileAsync{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.QueueCapacity != 0 {
		out.QueueCapacity = user.QueueCapacity
	}
	if user.WorkerLimit != 0 {
		out.WorkerLimit = user.WorkerLimit
	}
	if user.TargetTimeout != "" {
		out.TargetTimeout = user.TargetTimeout
	}
	if user.OnOverflow != "" {
		out.OnOverflow = user.OnOverflow
	}
	return &out
}

func mergeGuardsPtr(base, user *fileGuards) *fileGuards {
	if base == nil && user == nil {
		return nil
	}
	out := fileGuards{}
	if base != nil && base.Secrets != nil {
		s := *base.Secrets
		out.Secrets = &s
	}
	if user == nil || user.Secrets == nil {
		return &out
	}
	if out.Secrets == nil {
		out.Secrets = &fileSecretsGuard{}
	}
	if user.Secrets.Enabled != nil {
		out.Secrets.Enabled = user.Secrets.Enabled
	}
	if user.Secrets.Action != "" {
		out.Secrets.Action = user.Secrets.Action
	}
	if user.Secrets.Rules != nil {
		out.Secrets.Rules = append([]string(nil), user.Secrets.Rules...)
	}
	return &out
}
