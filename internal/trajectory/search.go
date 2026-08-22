package trajectory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const defaultSearchLimit = 100

// SearchOptions filters ledger JSONL scans. Search walks every matching session
// file line-by-line — O(total bytes) with no index (DESIGN §12 Q6).
type SearchOptions struct {
	Root      string
	Provider  string
	SessionID string
	Types     []string
	Source    string
	Query     string
	Limit     int
}

// SearchHit is one matching ledger event.
type SearchHit struct {
	Provider  string `json:"provider"`
	SessionID string `json:"session_id"`
	Seq       uint64 `json:"seq"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Snippet   string `json:"snippet,omitempty"`
	Path      string `json:"path"`
}

// Search scans session JSONL files under root and returns matching events.
func Search(opts SearchOptions) ([]SearchHit, error) {
	root := opts.Root
	if root == "" {
		root = DefaultSessionsDir()
	}
	if root == "" {
		return nil, ErrSessionsDirUnavailable
	}
	limit := opts.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	typeSet := map[string]bool{}
	for _, t := range opts.Types {
		if t = strings.TrimSpace(t); t != "" {
			typeSet[t] = true
		}
	}
	query := strings.ToLower(strings.TrimSpace(opts.Query))
	sourceFilter := strings.TrimSpace(opts.Source)
	providerFilter := CanonicalProvider(opts.Provider)
	sessionFilter := strings.TrimSpace(opts.SessionID)

	summaries, err := ListSessions(root, providerFilter)
	if err != nil {
		return nil, err
	}
	var hits []SearchHit
	for _, sum := range summaries {
		if sessionFilter != "" && sum.SessionID != sessionFilter {
			continue
		}
		fileHits, err := searchFile(sum.Path, sum.Provider, sum.SessionID, typeSet, sourceFilter, query, limit-len(hits))
		if err != nil {
			return nil, err
		}
		hits = append(hits, fileHits...)
		if len(hits) >= limit {
			break
		}
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Provider != hits[j].Provider {
			return hits[i].Provider < hits[j].Provider
		}
		if hits[i].SessionID != hits[j].SessionID {
			return hits[i].SessionID < hits[j].SessionID
		}
		return hits[i].Seq < hits[j].Seq
	})
	if len(hits) > limit {
		hits = hits[:limit]
	}
	return hits, nil
}

func searchFile(path, provider, sessionID string, types map[string]bool, source, query string, limit int) ([]SearchHit, error) {
	if limit <= 0 {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var hits []SearchHit
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("decode event in %s: %w", filepath.Base(path), err)
		}
		if len(types) > 0 && !types[e.Type] {
			continue
		}
		if source != "" && e.Source != source {
			continue
		}
		if query != "" && !eventMatchesQuery(e, query) {
			continue
		}
		hits = append(hits, SearchHit{
			Provider:  provider,
			SessionID: sessionID,
			Seq:       e.Seq,
			Type:      e.Type,
			Source:    e.Source,
			Snippet:   snippet(e),
			Path:      path,
		})
		if len(hits) >= limit {
			break
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return hits, nil
}

func eventMatchesQuery(e Event, query string) bool {
	if strings.Contains(strings.ToLower(e.Type), query) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Source), query) {
		return true
	}
	if len(e.Data) > 0 && strings.Contains(strings.ToLower(string(e.Data)), query) {
		return true
	}
	if len(e.Raw) > 0 && strings.Contains(strings.ToLower(string(e.Raw)), query) {
		return true
	}
	return false
}

func snippet(e Event) string {
	const maxSnippet = 120
	var s string
	if len(e.Data) > 0 {
		s = string(e.Data)
	} else if len(e.Raw) > 0 {
		s = string(e.Raw)
	} else {
		s = e.Type
	}
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxSnippet {
		return s[:maxSnippet] + "…"
	}
	return s
}
