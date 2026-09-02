package config

import "fmt"

type FailMode string

const (
	FailOpen   FailMode = "fail_open"
	FailClosed FailMode = "fail_closed"
)

type AskFallback string

const (
	AskFallbackDeny       AskFallback = "deny"
	AskFallbackNoDecision AskFallback = "no_decision"
)

// Policy is the compiled fail/ask policy.
type Policy struct {
	Fail        FailMode
	AskFallback AskFallback
	Offline     FailMode
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

func parseAskFallback(s string) (AskFallback, error) {
	switch AskFallback(s) {
	case AskFallbackDeny, AskFallbackNoDecision:
		return AskFallback(s), nil
	default:
		return "", fmt.Errorf("unknown %q", s)
	}
}
