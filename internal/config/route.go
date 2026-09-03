package config

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/speakeasy-api/agenthooks"
)

// KindDefault is a per-kind dispatch default.
type KindDefault struct {
	Mode DispatchMode
}

// eventKinds is the routable wire vocabulary, sourced from agenthooks so
// config and the decoder cannot drift. model.request and mcp.inventory are
// omitted: no provider renderer maps them.
var eventKinds = []string{
	string(agenthooks.KindSessionStart),
	string(agenthooks.KindSessionEnd),
	string(agenthooks.KindPromptSubmitted),
	string(agenthooks.KindToolPre),
	string(agenthooks.KindToolPost),
	string(agenthooks.KindToolError),
	string(agenthooks.KindPermission),
	string(agenthooks.KindStop),
	string(agenthooks.KindSubagentStart),
	string(agenthooks.KindSubagentStop),
	string(agenthooks.KindCompactPre),
	string(agenthooks.KindCompactPost),
	string(agenthooks.KindFileEdited),
	string(agenthooks.KindModelResponse),
	string(agenthooks.KindNotification),
	string(agenthooks.KindOther),
}

func validateEventKind(kind string) error {
	if slices.Contains(eventKinds, kind) {
		return nil
	}
	return fmt.Errorf("unknown kind %q (valid: %s)", kind, strings.Join(eventKinds, ", "))
}

// TargetKind identifies a dispatch target type.
type TargetKind string

const (
	TargetBuiltin TargetKind = "builtin"
	TargetExec    TargetKind = "exec"
	TargetHTTP    TargetKind = "http"
	TargetLog     TargetKind = "log"
	TargetFile    TargetKind = "file"
	TargetGRPC    TargetKind = "grpc"
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
	Name        string
	Kind        string // legacy single-kind key for default routes; empty when Match is set
	Match       RouteMatch
	Mode        DispatchMode
	SyncTimeout time.Duration // 0 = no route cap; used with provider timeout margin
	Sync        []CompiledTarget
	Async       []CompiledTarget
	Default     bool // true for routes synthesized from dispatch_defaults
}

func parseKindDefaults(in map[string]fileKindDefault, def map[string]KindDefault) (map[string]KindDefault, error) {
	out := make(map[string]KindDefault, len(def))
	maps.Copy(out, def)
	for k, v := range in {
		if err := validateEventKind(k); err != nil {
			return nil, fmt.Errorf("dispatch_defaults.%s: %w", k, err)
		}
		cur := out[k]
		if v.Mode != "" {
			m, err := parseDispatchMode(v.Mode)
			if err != nil {
				return nil, fmt.Errorf("dispatch_defaults.%s.mode: %w", k, err)
			}
			cur.Mode = m
		}
		out[k] = cur
	}
	return out, nil
}
