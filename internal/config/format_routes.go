package config

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// FormatRoutes writes compiled routes as a tab line each, or indented JSON.
func FormatRoutes(w io.Writer, routes []CompiledRoute, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(routes)
	}
	for _, r := range routes {
		kinds := r.Kind
		if len(r.Match.Kinds) > 0 {
			kinds = strings.Join(r.Match.Kinds, ",")
		}
		syncKinds := targetKinds(r.Sync)
		asyncKinds := targetKinds(r.Async)
		if _, err := fmt.Fprintf(w, "%s\tmatch.kind=%s\tmode=%s\tsync=[%s]\tasync=[%s]\n",
			r.Name, kinds, r.Mode, syncKinds, asyncKinds); err != nil {
			return err
		}
	}
	return nil
}

func targetKinds(ts []CompiledTarget) string {
	if len(ts) == 0 {
		return ""
	}
	parts := make([]string, len(ts))
	for i, t := range ts {
		parts[i] = string(t.Kind)
	}
	return strings.Join(parts, ",")
}
