package config

import "time"

// DefaultSecretsRules are the built-in rule names enabled when secrets.rules is empty.
var DefaultSecretsRules = []string{
	"aws_key",
	"github_pat",
	"github_fine_grained",
	"slack_token",
	"stripe_live",
	"anthropic_key",
	"openai_key",
	"google_api_key",
	"private_key",
	"jwt",
	"assigned_secret",
}

func defaultPolicy() Policy {
	return Policy{
		Fail:        FailClosed,
		Unsupported: UnsupportedDegrade,
		AskFallback: AskFallbackDeny,
		Offline:     FailClosed,
	}
}

func defaultAsync() AsyncConfig {
	return AsyncConfig{
		QueueCapacity: 1024,
		WorkerLimit:   8,
		TargetTimeout: 30 * time.Second,
		OnOverflow:    OverflowDrop,
	}
}

func defaultGuards() Guards {
	rules := make([]string, len(DefaultSecretsRules))
	copy(rules, DefaultSecretsRules)
	return Guards{
		Secrets: SecretsGuard{
			Enabled: true,
			Action:  GuardAsk,
			Rules:   rules,
		},
	}
}

func defaultKindDefaults() map[string]KindDefault {
	return map[string]KindDefault{
		"tool.pre":         {Mode: ModeParallel},
		"prompt.submitted": {Mode: ModeSyncOnly},
		"agent.stop":       {Mode: ModeSyncThenAsync},
		"tool.post":        {Mode: ModeParallel},
		"notification":     {Mode: ModeAsyncOnly},
		"other":            {Mode: ModeAsyncOnly},
	}
}
