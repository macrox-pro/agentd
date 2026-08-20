package config

import (
	"fmt"
	"strings"
	"time"
)

// TemporaryBlock is one non-expired runtime deny rule.
type TemporaryBlock struct {
	Tool    string
	Pattern string
	Reason  string
	Until   time.Time
}

func parseTemporaryBlocks(in *fileBlocks, now time.Time) ([]TemporaryBlock, error) {
	if in == nil || len(in.Temporary) == 0 {
		return nil, nil
	}
	out := make([]TemporaryBlock, 0, len(in.Temporary))
	for i, fb := range in.Temporary {
		tb, ok, err := parseOneTemporaryBlock(fb, now)
		if err != nil {
			return nil, fmt.Errorf("blocks.temporary[%d]: %w", i, err)
		}
		if ok {
			out = append(out, tb)
		}
	}
	return out, nil
}

func parseOneTemporaryBlock(fb fileTemporaryBlock, now time.Time) (TemporaryBlock, bool, error) {
	if fb.Tool == "" {
		return TemporaryBlock{}, false, fmt.Errorf("tool is required")
	}
	if fb.Pattern == "" {
		return TemporaryBlock{}, false, fmt.Errorf("pattern is required")
	}
	if fb.Until == "" {
		return TemporaryBlock{}, false, fmt.Errorf("until is required")
	}
	until, err := time.Parse(time.RFC3339, fb.Until)
	if err != nil {
		return TemporaryBlock{}, false, fmt.Errorf("until: %w", err)
	}
	if !until.After(now) {
		return TemporaryBlock{}, false, nil
	}
	return TemporaryBlock{
		Tool:    fb.Tool,
		Pattern: fb.Pattern,
		Reason:  fb.Reason,
		Until:   until,
	}, true, nil
}

// MatchTemporaryBlock returns the first active block matching tool name and pattern substring.
func MatchTemporaryBlock(blocks []TemporaryBlock, toolName, haystack string, now time.Time) *TemporaryBlock {
	for i := range blocks {
		b := &blocks[i]
		if !b.Until.After(now) {
			continue
		}
		if !toolMatches(b.Tool, toolName) {
			continue
		}
		if b.Pattern != "" && !strings.Contains(haystack, b.Pattern) {
			continue
		}
		return b
	}
	return nil
}

func toolMatches(want, got string) bool {
	if want == "*" || want == "" {
		return true
	}
	return strings.EqualFold(want, got)
}

func temporaryBlockKey(tool, pattern string) string {
	return tool + "\x00" + pattern
}

func upsertTemporaryList(base, overlay []fileTemporaryBlock) []fileTemporaryBlock {
	if overlay == nil {
		return append([]fileTemporaryBlock(nil), base...)
	}
	if len(base) == 0 {
		return append([]fileTemporaryBlock(nil), overlay...)
	}
	out := append([]fileTemporaryBlock(nil), base...)
	index := map[string]int{}
	for i, b := range out {
		index[temporaryBlockKey(b.Tool, b.Pattern)] = i
	}
	for _, b := range overlay {
		key := temporaryBlockKey(b.Tool, b.Pattern)
		if i, ok := index[key]; ok {
			out[i] = b
			continue
		}
		index[key] = len(out)
		out = append(out, b)
	}
	return out
}

func mergeBlocksPtr(base, overlay *fileBlocks) *fileBlocks {
	if base == nil && overlay == nil {
		return nil
	}
	out := fileBlocks{}
	if base != nil {
		out.Temporary = append([]fileTemporaryBlock(nil), base.Temporary...)
	}
	if overlay == nil {
		return &out
	}
	if overlay.Temporary != nil {
		out.Temporary = upsertTemporaryList(out.Temporary, overlay.Temporary)
	}
	return &out
}
