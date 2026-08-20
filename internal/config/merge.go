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
	if base != nil {
		if base.Secrets != nil {
			s := *base.Secrets
			out.Secrets = &s
		}
		if base.Shell != nil {
			s := *base.Shell
			out.Shell = &s
		}
		if base.MCP != nil {
			m := *base.MCP
			out.MCP = &m
		}
		if base.Paths != nil {
			p := *base.Paths
			out.Paths = &p
		}
	}
	if user == nil {
		return &out
	}
	out.Secrets = mergeSecretsGuardPtr(out.Secrets, user.Secrets)
	out.Shell = mergeShellGuardPtr(out.Shell, user.Shell)
	out.MCP = mergeMCPGuardPtr(out.MCP, user.MCP)
	out.Paths = mergePathsGuardPtr(out.Paths, user.Paths)
	return &out
}

func mergeSecretsGuardPtr(base, user *fileSecretsGuard) *fileSecretsGuard {
	if base == nil && user == nil {
		return nil
	}
	out := fileSecretsGuard{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Enabled != nil {
		out.Enabled = user.Enabled
	}
	if user.Action != "" {
		out.Action = user.Action
	}
	if user.Rules != nil {
		out.Rules = append([]string(nil), user.Rules...)
	}
	return &out
}

func mergeShellGuardPtr(base, user *fileShellGuard) *fileShellGuard {
	if base == nil && user == nil {
		return nil
	}
	out := fileShellGuard{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Enabled != nil {
		out.Enabled = user.Enabled
	}
	if user.DenyPatterns != nil {
		out.DenyPatterns = append([]string(nil), user.DenyPatterns...)
	}
	if user.AskOn != nil {
		out.AskOn = append([]string(nil), user.AskOn...)
	}
	return &out
}

func mergeMCPGuardPtr(base, user *fileMCPGuard) *fileMCPGuard {
	if base == nil && user == nil {
		return nil
	}
	out := fileMCPGuard{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Enabled != nil {
		out.Enabled = user.Enabled
	}
	if user.DenyServers != nil {
		out.DenyServers = append([]string(nil), user.DenyServers...)
	}
	return &out
}

func mergePathsGuardPtr(base, user *filePathsGuard) *filePathsGuard {
	if base == nil && user == nil {
		return nil
	}
	out := filePathsGuard{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Enabled != nil {
		out.Enabled = user.Enabled
	}
	if user.DenyRead != nil {
		out.DenyRead = append([]string(nil), user.DenyRead...)
	}
	if user.DenyWrite != nil {
		out.DenyWrite = append([]string(nil), user.DenyWrite...)
	}
	return &out
}
