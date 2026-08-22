package trajectory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

// SessionKey identifies one ledger stream.
type SessionKey struct {
	Provider    string
	SessionID   string
	ProjectRoot string
}

func (k SessionKey) String() string {
	return fmt.Sprintf("%s:%s:%s", k.Provider, k.SessionID, k.ProjectRoot)
}

// CanonicalProvider normalizes CLI/provider ids for ledger keys.
func CanonicalProvider(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "kimicode", "kimi-code":
		return "kimi-code"
	case "claude-code":
		return "claude-code"
	case "cursor":
		return "cursor"
	case "codex":
		return "codex"
	case "gemini":
		return "gemini"
	case "opencode":
		return "opencode"
	default:
		return strings.TrimSpace(name)
	}
}

// ResolveSessionKey builds a stable ledger key; empty session_id gets a weak synthetic id.
func ResolveSessionKey(provider, sessionID, projectRoot, cwd string) SessionKey {
	prov := CanonicalProvider(provider)
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		sid = weakSessionID(prov, projectRoot, cwd)
	}
	return SessionKey{
		Provider:    prov,
		SessionID:   sid,
		ProjectRoot: strings.TrimSpace(projectRoot),
	}
}

func weakSessionID(provider, projectRoot, cwd string) string {
	sum := sha256.Sum256([]byte(provider + "\x00" + projectRoot + "\x00" + cwd))
	return "weak-" + hex.EncodeToString(sum[:8])
}

// SessionFileName returns the on-disk basename for a session id.
func SessionFileName(sessionID string) string {
	if sessionID == "" {
		return "_anonymous.jsonl"
	}
	safe := strings.ReplaceAll(sessionID, "/", "_")
	safe = strings.ReplaceAll(safe, string(filepath.Separator), "_")
	return safe + ".jsonl"
}
