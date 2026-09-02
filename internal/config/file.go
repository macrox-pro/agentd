package config

// fileConfig is the on-disk YAML shape.
type fileConfig struct {
	Version          int                        `yaml:"version,omitempty"`
	Policy           *filePolicy                `yaml:"policy,omitempty"`
	Async            *fileAsync                 `yaml:"async,omitempty"`
	Guards           *fileGuards                `yaml:"guards,omitempty"`
	Approvals        *fileApprovals             `yaml:"approvals,omitempty"`
	Blocks           *fileBlocks                `yaml:"blocks,omitempty"`
	DispatchDefaults map[string]fileKindDefault `yaml:"dispatch_defaults,omitempty"`
	Dispatch         []fileRoute                `yaml:"dispatch,omitempty"`
	Trajectory       *fileTrajectory            `yaml:"trajectory,omitempty"`
	Logging          *fileLogging               `yaml:"logging,omitempty"`
	Metrics          *fileMetrics               `yaml:"metrics,omitempty"`
}

type fileApprovals struct {
	Secrets []fileApproval `yaml:"secrets,omitempty"`
	Shell   []fileApproval `yaml:"shell,omitempty"`
}

type fileApproval struct {
	Fingerprint string `yaml:"fingerprint,omitempty"`
	Scope       string `yaml:"scope,omitempty"`
	Project     string `yaml:"project,omitempty"`
	SessionID   string `yaml:"session_id,omitempty"`
	ExpiresAt   string `yaml:"expires_at,omitempty"`
	GrantedBy   string `yaml:"granted_by,omitempty"`
}

type fileBlocks struct {
	Temporary []fileTemporaryBlock `yaml:"temporary,omitempty"`
}

type fileTemporaryBlock struct {
	Tool    string `yaml:"tool,omitempty"`
	Pattern string `yaml:"pattern,omitempty"`
	Reason  string `yaml:"reason,omitempty"`
	Until   string `yaml:"until,omitempty"`
}

type filePolicy struct {
	Fail        string `yaml:"fail,omitempty"`
	AskFallback string `yaml:"ask_fallback,omitempty"`
	Offline     string `yaml:"offline,omitempty"`
}

type fileAsync struct {
	QueueCapacity int    `yaml:"queue_capacity,omitempty"`
	WorkerLimit   int    `yaml:"worker_limit,omitempty"`
	TargetTimeout string `yaml:"target_timeout,omitempty"`
	OnOverflow    string `yaml:"on_overflow,omitempty"`
}

type fileGuards struct {
	Secrets *fileSecretsGuard `yaml:"secrets,omitempty"`
	Shell   *fileShellGuard   `yaml:"shell,omitempty"`
	MCP     *fileMCPGuard     `yaml:"mcp,omitempty"`
	Paths   *filePathsGuard   `yaml:"paths,omitempty"`
}

type fileSecretsGuard struct {
	Enabled *bool    `yaml:"enabled,omitempty"`
	Action  string   `yaml:"action,omitempty"`
	Rules   []string `yaml:"rules,omitempty"`
}

type fileShellGuard struct {
	Enabled      *bool    `yaml:"enabled,omitempty"`
	DenyPatterns []string `yaml:"deny_patterns,omitempty"`
	AskOn        []string `yaml:"ask_on,omitempty"`
}

type fileMCPGuard struct {
	Enabled     *bool    `yaml:"enabled,omitempty"`
	DenyServers []string `yaml:"deny_servers,omitempty"`
}

type filePathsGuard struct {
	Enabled   *bool    `yaml:"enabled,omitempty"`
	DenyRead  []string `yaml:"deny_read,omitempty"`
	DenyWrite []string `yaml:"deny_write,omitempty"`
}

type fileKindDefault struct {
	Mode string `yaml:"mode,omitempty"`
}

type fileRoute struct {
	Name        string       `yaml:"name,omitempty"`
	Match       fileMatch    `yaml:"match,omitempty"`
	Mode        string       `yaml:"mode,omitempty"`
	SyncTimeout string       `yaml:"sync_timeout,omitempty"`
	Sync        []fileTarget `yaml:"sync,omitempty"`
	Async       []fileTarget `yaml:"async,omitempty"`
}

type fileMatch struct {
	Kind     []string `yaml:"kind,omitempty"`
	Provider []string `yaml:"provider,omitempty"`
	Tools    []string `yaml:"tools,omitempty"`
}

type fileTarget struct {
	Target   string   `yaml:"target,omitempty"`
	Guards   []string `yaml:"guards,omitempty"`
	Observe  bool     `yaml:"observe,omitempty"`
	URL      string   `yaml:"url,omitempty"`
	Command  []string `yaml:"command,omitempty"`
	Stdin    string   `yaml:"stdin,omitempty"`
	Level    string   `yaml:"level,omitempty"`
	Path     string   `yaml:"path,omitempty"`
	Retry    *int     `yaml:"retry,omitempty"`
	Timeout  string   `yaml:"timeout,omitempty"`
	Endpoint string   `yaml:"endpoint,omitempty"`
	OnError  string   `yaml:"on_error,omitempty"`
	Merge    string   `yaml:"merge,omitempty"`
}
