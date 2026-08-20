package config

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
	Mode string `yaml:"mode"`
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
