package trajectory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/macrox-pro/agentd/internal/provider"
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
	if id, ok := provider.Lookup(name); ok {
		return string(id)
	}
	return strings.TrimSpace(name)
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

// ResolveSessionKeyID is ResolveSessionKey with a validated provider.ID.
func ResolveSessionKeyID(id provider.ID, sessionID, projectRoot, cwd string) SessionKey {
	return ResolveSessionKey(string(id), sessionID, projectRoot, cwd)
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
