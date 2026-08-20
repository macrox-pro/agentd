package guard

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

func pathsHandler(cfg config.PathsGuard) func(context.Context, *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
	return func(ctx context.Context, e *agenthooks.ToolPreEvent) (agenthooks.ToolPreDecision, error) {
		patterns, ok := pathDenyPatterns(e.Tool.Canonical, cfg)
		if !ok || len(patterns) == 0 {
			return agenthooks.NoDecision(), nil
		}
		candidates := extractPaths(e.Tool.Input)
		for _, candidate := range candidates {
			if hit := matchPathGlob(candidate, patterns); hit != "" {
				agenthooks.Logger(ctx).Warn("paths guard: deny pattern matched",
					"tool", e.Tool.Name, "path", candidate, "pattern", hit)
				return agenthooks.Deny(fmt.Sprintf(
					"path access blocked: %q matched deny pattern %q", candidate, hit,
				)), nil
			}
		}
		return agenthooks.NoDecision(), nil
	}
}

func pathDenyPatterns(canonical agenthooks.CanonicalTool, cfg config.PathsGuard) ([]string, bool) {
	switch canonical {
	case agenthooks.ToolFileRead:
		return cfg.DenyRead, true
	case agenthooks.ToolFileWrite, agenthooks.ToolFileEdit:
		return cfg.DenyWrite, true
	default:
		return nil, false
	}
}

var pathInputKeys = []string{
	"file_path", "path", "target_file", "filePath", "filename", "file",
}

func extractPaths(input json.RawMessage) []string {
	var obj map[string]any
	if err := json.Unmarshal(input, &obj); err != nil {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	for _, key := range pathInputKeys {
		v, ok := obj[key]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

func matchPathGlob(candidate string, patterns []string) string {
	for _, p := range patterns {
		if p == "" {
			continue
		}
		if matchDoublestar(p, candidate) {
			return p
		}
	}
	return ""
}

// matchDoublestar matches path-style globs where * does not cross '/' and **
// matches across path segments (including zero segments).
func matchDoublestar(pattern, name string) bool {
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	name = strings.ReplaceAll(name, "\\", "/")
	return matchDoublestarParts(splitPath(pattern), splitPath(name))
}

func splitPath(p string) []string {
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

func matchDoublestarParts(pat, name []string) bool {
	for len(pat) > 0 {
		if pat[0] == "**" {
			pat = pat[1:]
			if len(pat) == 0 {
				return true
			}
			for i := 0; i <= len(name); i++ {
				if matchDoublestarParts(pat, name[i:]) {
					return true
				}
			}
			return false
		}
		if len(name) == 0 {
			return false
		}
		if !matchSegment(pat[0], name[0]) {
			return false
		}
		pat = pat[1:]
		name = name[1:]
	}
	return len(name) == 0
}

func matchSegment(pat, seg string) bool {
	// path.Match semantics within one segment (* and ?).
	i, j := 0, 0
	star := -1
	match := 0
	for j < len(seg) {
		if i < len(pat) && (pat[i] == '?' || pat[i] == seg[j]) {
			i++
			j++
			continue
		}
		if i < len(pat) && pat[i] == '*' {
			star = i
			match = j
			i++
			continue
		}
		if star >= 0 {
			i = star + 1
			match++
			j = match
			continue
		}
		return false
	}
	for i < len(pat) && pat[i] == '*' {
		i++
	}
	return i == len(pat)
}
