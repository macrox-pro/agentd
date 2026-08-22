package trajectory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

func writeSessionJSONL(t *testing.T, root, provider, sessionID, content string) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	path := filepath.Join(dir, sessionID+".jsonl")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
}

func TestSearchFilterByProviderAndQuery(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeSessionJSONL(t, root, "claude-code", "s1", `{"seq":1,"type":"hook/invoked","source":"hook","session_id":"s1","provider":"claude-code","data":{"kind":"tool.pre"}}
{"seq":2,"type":"transcript/message","source":"transcript","session_id":"s1","provider":"claude-code","data":{"text":"hello thinking world"}}
`)
	writeSessionJSONL(t, root, "cursor", "s2", `{"seq":1,"type":"hook/invoked","source":"hook","session_id":"s2","provider":"cursor","data":{"kind":"tool.pre"}}
`)

	hits, err := trajectory.Search(trajectory.SearchOptions{
		Root:     root,
		Provider: "claude-code",
		Query:    "thinking",
	})
	require.NoError(t, err, "Search")
	require.Len(t, hits, 1)
	assert.Equal(t, "transcript/message", hits[0].Type)
	assert.Equal(t, "s1", hits[0].SessionID)
}

func TestSearchFilterByKindAndSource(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeSessionJSONL(t, root, "claude-code", "s1", `{"seq":1,"type":"hook/invoked","source":"hook","session_id":"s1","provider":"claude-code"}
{"seq":2,"type":"transcript/thinking","source":"transcript","session_id":"s1","provider":"claude-code","data":{"text":"reason"}}
`)

	hits, err := trajectory.Search(trajectory.SearchOptions{
		Root:   root,
		Types:  []string{trajectory.TypeTranscriptThinking},
		Source: trajectory.SourceTranscript,
	})
	require.NoError(t, err, "Search")
	require.Len(t, hits, 1)
	assert.Equal(t, uint64(2), hits[0].Seq)
}

func TestSearchEmptyRootNoSessions(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	hits, err := trajectory.Search(trajectory.SearchOptions{Root: root, Query: "x"})
	require.NoError(t, err, "Search")
	assert.Empty(t, hits)
}

func TestSearchCorruptLineReturnsError(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeSessionJSONL(t, root, "gemini", "s3", `{not json`)
	_, err := trajectory.Search(trajectory.SearchOptions{Root: root})
	require.Error(t, err, "Search corrupt")
}

func TestSearchMultiSessionRespectsLimit(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeSessionJSONL(t, root, "codex", "a", `{"seq":1,"type":"hook/invoked","source":"hook","session_id":"a","provider":"codex","data":{"tool":"x"}}
`)
	writeSessionJSONL(t, root, "codex", "b", `{"seq":1,"type":"hook/invoked","source":"hook","session_id":"b","provider":"codex","data":{"tool":"x"}}
`)

	hits, err := trajectory.Search(trajectory.SearchOptions{Root: root, Query: "tool", Limit: 1})
	require.NoError(t, err, "Search")
	require.Len(t, hits, 1)
}
