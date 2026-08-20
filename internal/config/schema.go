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
	Fail        FailMode
	Unsupported UnsupportedMode
	AskFallback AskFallback
	Offline     FailMode
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
	TargetExec    TargetKind = "exec"
	TargetHTTP    TargetKind = "http"
	TargetLog     TargetKind = "log"
	TargetFile    TargetKind = "file"
	TargetGRPC TargetKind = "grpc"
)

// SyncMerge is the sync merge policy for a grpc target (engine uses first_conclusive).
type SyncMerge string

const (
	MergeFirstConclusive SyncMerge = "first_conclusive"
)

// RouteMatch is compiled match criteria for a declarative route.
type RouteMatch struct {
	Kinds     []string // empty = any
	Providers []string // empty or ["*"] = any
	Tools     []string // empty = any; tool name or canonical class
}

// CompiledTarget is one sync or async target on a route.
type CompiledTarget struct {
	Kind     TargetKind
	Guards   []string // builtin sync
	Observe  bool     // builtin async
	URL      string   // http
	Command  []string // exec
	Stdin    string   // exec: "raw" or empty
	Level    string   // log: info|warn|error|debug
	Path     string   // file
	Retry    int      // http (M3: always 0)
	Timeout  time.Duration
	Endpoint string   // grpc
	OnError  FailMode // grpc sync: fail_closed (default) | fail_open
	Merge    SyncMerge
}

// CompiledRoute is a compiled dispatch route.
type CompiledRoute struct {
	Name    string
	Kind    string // legacy single-kind key for default routes; empty when Match is set
	Match   RouteMatch
	Mode    DispatchMode
	Sync    []CompiledTarget
	Async   []CompiledTarget
	Default bool // true for routes synthesized from dispatch_defaults
}

// fileConfig is the on-disk YAML shape.
type fileConfig struct {
	Version          int                        `yaml:"version"`
	Policy           *filePolicy                `yaml:"policy"`
	Async            *fileAsync                 `yaml:"async"`
	Guards           *fileGuards                `yaml:"guards"`
	DispatchDefaults map[string]fileKindDefault `yaml:"dispatch_defaults"`
	Dispatch         []fileRoute                `yaml:"dispatch"`
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

type fileRoute struct {
	Name  string       `yaml:"name"`
	Match fileMatch    `yaml:"match"`
	Mode  string       `yaml:"mode"`
	Sync  []fileTarget `yaml:"sync"`
	Async []fileTarget `yaml:"async"`
}

type fileMatch struct {
	Kind     []string `yaml:"kind"`
	Provider []string `yaml:"provider"`
	Tools    []string `yaml:"tools"`
}

type fileTarget struct {
	Target   string   `yaml:"target"`
	Guards   []string `yaml:"guards"`
	Observe  bool     `yaml:"observe"`
	URL      string   `yaml:"url"`
	Command  []string `yaml:"command"`
	Stdin    string   `yaml:"stdin"`
	Level    string   `yaml:"level"`
	Path     string   `yaml:"path"`
	Retry    *int     `yaml:"retry"`
	Timeout  string   `yaml:"timeout"`
	Endpoint string   `yaml:"endpoint"`
	OnError  string   `yaml:"on_error"`
	Merge    string   `yaml:"merge"`
}
