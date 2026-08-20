package config

import "time"

// FailMode is policy.fail.
type FailMode string

const (
	FailOpen   FailMode = "fail_open"
	FailClosed FailMode = "fail_closed"
)

// UnsupportedMode is policy.unsupported.
type UnsupportedMode string

const (
	UnsupportedDegrade UnsupportedMode = "degrade"
	UnsupportedStrict  UnsupportedMode = "strict"
)

// AskFallback is policy.ask_fallback.
type AskFallback string

const (
	AskFallbackDeny       AskFallback = "deny"
	AskFallbackNoDecision AskFallback = "no_decision"
)

// OverflowMode is async.on_overflow.
type OverflowMode string

const (
	OverflowDrop OverflowMode = "drop"
	OverflowLog  OverflowMode = "log"
)

// GuardAction is guards.*.action.
type GuardAction string

const (
	GuardAsk  GuardAction = "ask"
	GuardDeny GuardAction = "deny"
)

// DispatchMode is a route dispatch mode.
type DispatchMode string

const (
	ModeSyncOnly      DispatchMode = "sync_only"
	ModeAsyncOnly     DispatchMode = "async_only"
	ModeParallel      DispatchMode = "parallel"
	ModeAfterSync     DispatchMode = "after_sync"
	ModeSyncThenAsync DispatchMode = "sync_then_async" // alias for after_sync
)

// Policy is the compiled fail/ask policy.
type Policy struct {
	Fail         FailMode
	Unsupported  UnsupportedMode
	AskFallback  AskFallback
	Offline      FailMode
}

// AsyncConfig is the compiled async queue settings.
type AsyncConfig struct {
	QueueCapacity int
	WorkerLimit   int
	TargetTimeout time.Duration
	OnOverflow    OverflowMode
}

// SecretsGuard is compiled secrets guard settings.
type SecretsGuard struct {
	Enabled bool
	Action  GuardAction
	Rules   []string
}

// Guards holds compiled guard settings.
type Guards struct {
	Secrets SecretsGuard
}

// KindDefault is a per-kind dispatch default.
type KindDefault struct {
	Mode     DispatchMode
	Blocking bool
}

// TargetKind identifies a dispatch target type.
type TargetKind string

const (
	TargetBuiltin TargetKind = "builtin"
)

// CompiledTarget is one sync or async target on a route.
type CompiledTarget struct {
	Kind    TargetKind
	Guards  []string // for builtin sync
	Observe bool     // for builtin async
}

// CompiledRoute is a compiled dispatch route (M2: from dispatch_defaults).
type CompiledRoute struct {
	Name   string
	Kind   string // event kind key, e.g. "tool.pre"
	Mode   DispatchMode
	Sync   []CompiledTarget
	Async  []CompiledTarget
}

// fileConfig is the on-disk YAML shape (partial; M2 subset).
type fileConfig struct {
	Version          int                       `yaml:"version"`
	Policy           *filePolicy               `yaml:"policy"`
	Async            *fileAsync                `yaml:"async"`
	Guards           *fileGuards               `yaml:"guards"`
	DispatchDefaults map[string]fileKindDefault `yaml:"dispatch_defaults"`
}

type filePolicy struct {
	Fail        string `yaml:"fail"`
	Unsupported string `yaml:"unsupported"`
	AskFallback string `yaml:"ask_fallback"`
	Offline     string `yaml:"offline"`
}

type fileAsync struct {
	QueueCapacity int    `yaml:"queue_capacity"`
	WorkerLimit   int    `yaml:"worker_limit"`
	TargetTimeout string `yaml:"target_timeout"`
	OnOverflow    string `yaml:"on_overflow"`
}

type fileGuards struct {
	Secrets *fileSecretsGuard `yaml:"secrets"`
}

type fileSecretsGuard struct {
	Enabled *bool    `yaml:"enabled"`
	Action  string   `yaml:"action"`
	Rules   []string `yaml:"rules"`
}

type fileKindDefault struct {
	Mode     string `yaml:"mode"`
	Blocking *bool  `yaml:"blocking"`
}
