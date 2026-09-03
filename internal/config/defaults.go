package config

import (
	"time"

	"github.com/speakeasy-api/agenthooks"
)

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
		AskFallback: AskFallbackDeny,
		Offline:     FailOpen,
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
		Shell: ShellGuard{Enabled: false},
		MCP:   MCPGuard{Enabled: false},
		Paths: PathsGuard{Enabled: false},
	}
}

func defaultKindDefaults() map[string]KindDefault {
	return map[string]KindDefault{
		string(agenthooks.KindToolPre):         {Mode: ModeParallel},
		string(agenthooks.KindPromptSubmitted): {Mode: ModeSyncOnly},
		string(agenthooks.KindStop):            {Mode: ModeSyncThenAsync},
		string(agenthooks.KindToolPost):        {Mode: ModeParallel},
		string(agenthooks.KindNotification):    {Mode: ModeAsyncOnly},
		string(agenthooks.KindOther):           {Mode: ModeAsyncOnly},
		string(agenthooks.KindSessionStart):    {Mode: ModeAsyncOnly},
		string(agenthooks.KindSessionEnd):      {Mode: ModeAsyncOnly},
		string(agenthooks.KindToolError):       {Mode: ModeAsyncOnly},
		string(agenthooks.KindPermission):      {Mode: ModeAsyncOnly},
		string(agenthooks.KindSubagentStart):   {Mode: ModeAsyncOnly},
		string(agenthooks.KindSubagentStop):    {Mode: ModeAsyncOnly},
		string(agenthooks.KindCompactPre):      {Mode: ModeAsyncOnly},
		string(agenthooks.KindCompactPost):     {Mode: ModeAsyncOnly},
		string(agenthooks.KindFileEdited):      {Mode: ModeAsyncOnly},
		string(agenthooks.KindModelResponse):   {Mode: ModeAsyncOnly},
	}
}
