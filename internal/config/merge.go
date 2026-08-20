package config

import (
	"fmt"
	"time"
)

func mergeFile(base *fileConfig, user *fileConfig) *fileConfig {
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
		for k, v := range out.DispatchDefaults {
			dd[k] = v
		}
		for k, v := range user.DispatchDefaults {
			cur := dd[k]
			if v.Mode != "" {
				cur.Mode = v.Mode
			}
			if v.Blocking != nil {
				cur.Blocking = v.Blocking
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

func parsePolicy(fp *filePolicy, def Policy) (Policy, error) {
	out := def
	if fp == nil {
		return out, nil
	}
	if fp.Fail != "" {
		m, err := parseFailMode(fp.Fail)
		if err != nil {
			return Policy{}, fmt.Errorf("policy.fail: %w", err)
		}
		out.Fail = m
	}
	if fp.Unsupported != "" {
		m, err := parseUnsupported(fp.Unsupported)
		if err != nil {
			return Policy{}, fmt.Errorf("policy.unsupported: %w", err)
		}
		out.Unsupported = m
	}
	if fp.AskFallback != "" {
		m, err := parseAskFallback(fp.AskFallback)
		if err != nil {
			return Policy{}, fmt.Errorf("policy.ask_fallback: %w", err)
		}
		out.AskFallback = m
	}
	if fp.Offline != "" {
		m, err := parseFailMode(fp.Offline)
		if err != nil {
			return Policy{}, fmt.Errorf("policy.offline: %w", err)
		}
		out.Offline = m
	}
	return out, nil
}

func parseFailMode(s string) (FailMode, error) {
	switch FailMode(s) {
	case FailOpen, FailClosed:
		return FailMode(s), nil
	default:
		return "", fmt.Errorf("unknown %q", s)
	}
}

func parseUnsupported(s string) (UnsupportedMode, error) {
	switch UnsupportedMode(s) {
	case UnsupportedDegrade, UnsupportedStrict:
		return UnsupportedMode(s), nil
	default:
		return "", fmt.Errorf("unknown %q", s)
	}
}

func parseAskFallback(s string) (AskFallback, error) {
	switch AskFallback(s) {
	case AskFallbackDeny, AskFallbackNoDecision:
		return AskFallback(s), nil
	default:
		return "", fmt.Errorf("unknown %q", s)
	}
}

func parseAsync(fa *fileAsync, def AsyncConfig) (AsyncConfig, error) {
	out := def
	if fa == nil {
		return out, nil
	}
	if fa.QueueCapacity != 0 {
		if fa.QueueCapacity < 1 {
			return AsyncConfig{}, fmt.Errorf("async.queue_capacity must be >= 1")
		}
		out.QueueCapacity = fa.QueueCapacity
	}
	if fa.WorkerLimit != 0 {
		if fa.WorkerLimit < 1 {
			return AsyncConfig{}, fmt.Errorf("async.worker_limit must be >= 1")
		}
		out.WorkerLimit = fa.WorkerLimit
	}
	if fa.TargetTimeout != "" {
		d, err := time.ParseDuration(fa.TargetTimeout)
		if err != nil {
			return AsyncConfig{}, fmt.Errorf("async.target_timeout: %w", err)
		}
		if d <= 0 {
			return AsyncConfig{}, fmt.Errorf("async.target_timeout must be > 0")
		}
		out.TargetTimeout = d
	}
	if fa.OnOverflow != "" {
		switch OverflowMode(fa.OnOverflow) {
		case OverflowDrop, OverflowLog:
			out.OnOverflow = OverflowMode(fa.OnOverflow)
		default:
			return AsyncConfig{}, fmt.Errorf("async.on_overflow: unknown %q", fa.OnOverflow)
		}
	}
	return out, nil
}

func parseGuards(fg *fileGuards, def Guards) (Guards, error) {
	out := def
	if fg == nil || fg.Secrets == nil {
		return out, nil
	}
	s := fg.Secrets
	if s.Enabled != nil {
		out.Secrets.Enabled = *s.Enabled
	}
	if s.Action != "" {
		switch GuardAction(s.Action) {
		case GuardAsk, GuardDeny:
			out.Secrets.Action = GuardAction(s.Action)
		default:
			return Guards{}, fmt.Errorf("guards.secrets.action: unknown %q", s.Action)
		}
	}
	if s.Rules != nil {
		out.Secrets.Rules = append([]string(nil), s.Rules...)
	}
	return out, nil
}

func parseKindDefaults(in map[string]fileKindDefault, def map[string]KindDefault) (map[string]KindDefault, error) {
	out := make(map[string]KindDefault, len(def))
	for k, v := range def {
		out[k] = v
	}
	for k, v := range in {
		cur := out[k]
		if v.Mode != "" {
			m, err := parseDispatchMode(v.Mode)
			if err != nil {
				return nil, fmt.Errorf("dispatch_defaults.%s.mode: %w", k, err)
			}
			cur.Mode = m
		}
		if v.Blocking != nil {
			cur.Blocking = *v.Blocking
		}
		out[k] = cur
	}
	return out, nil
}

func parseDispatchMode(s string) (DispatchMode, error) {
	switch DispatchMode(s) {
	case ModeSyncOnly, ModeAsyncOnly, ModeParallel, ModeAfterSync, ModeSyncThenAsync:
		return DispatchMode(s), nil
	default:
		return "", fmt.Errorf("unknown %q", s)
	}
}

// NormalizeMode maps aliases to canonical modes.
func NormalizeMode(m DispatchMode) DispatchMode {
	if m == ModeSyncThenAsync {
		return ModeAfterSync
	}
	return m
}
