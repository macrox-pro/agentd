package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	approvalFingerprintPrefix = "sha256:"
	grantedByAskUser          = "ask_user"
	projectApprovalTTL        = 24 * time.Hour
)

// ApprovalScope is the lifetime binding for a recorded approval.
type ApprovalScope string

const (
	ApprovalScopeProject ApprovalScope = "project"
	ApprovalScopeSession ApprovalScope = "session"
)

// ApprovalKind names the Ask-capable guard that granted the approval.
type ApprovalKind string

const (
	ApprovalKindSecrets ApprovalKind = "secrets"
	ApprovalKindShell   ApprovalKind = "shell"
)

// Approval is one non-expired runtime approval entry.
type Approval struct {
	Kind        ApprovalKind
	Fingerprint string
	Scope       ApprovalScope
	Project     string
	SessionID   string
	ExpiresAt   time.Time // zero = no wall-clock expiry (session)
	GrantedBy   string
}

// Approvals is the compiled set of active approvals by kind.
type Approvals struct {
	Secrets []Approval
	Shell   []Approval
}

// ApprovalFingerprint builds the stable approval id for kind+tool+stableKey.
// Format: sha256:<kind>/<hex>. stableKey must not contain secret material.
func ApprovalFingerprint(kind ApprovalKind, tool, stableKey string) string {
	sum := sha256.Sum256([]byte(string(kind) + "\x00" + tool + "\x00" + stableKey))
	return approvalFingerprintPrefix + string(kind) + "/" + hex.EncodeToString(sum[:])
}

// ParseApprovalKind extracts the guard kind embedded in an approval fingerprint.
func ParseApprovalKind(fingerprint string) (ApprovalKind, error) {
	rest, ok := strings.CutPrefix(fingerprint, approvalFingerprintPrefix)
	if !ok || rest == "" {
		return "", fmt.Errorf("invalid approval fingerprint")
	}
	kindStr, _, ok := strings.Cut(rest, "/")
	if !ok || kindStr == "" {
		return "", fmt.Errorf("approval fingerprint missing kind")
	}
	switch ApprovalKind(kindStr) {
	case ApprovalKindSecrets, ApprovalKindShell:
		return ApprovalKind(kindStr), nil
	default:
		return "", fmt.Errorf("unknown approval kind %q", kindStr)
	}
}

// SecretsStableKey returns the sorted rule-id key for secrets approvals.
func SecretsStableKey(ruleIDs []string) string {
	if len(ruleIDs) == 0 {
		return ""
	}
	cp := append([]string(nil), ruleIDs...)
	sort.Strings(cp)
	return strings.Join(cp, ",")
}

// HasApproval reports whether a non-expired matching approval exists.
func (a Approvals) HasApproval(kind ApprovalKind, fingerprint, project, sessionID string, now time.Time) bool {
	for _, ap := range a.byKind(kind) {
		if !approvalMatches(ap, fingerprint, project, sessionID, now) {
			continue
		}
		return true
	}
	return false
}

func (a Approvals) byKind(kind ApprovalKind) []Approval {
	switch kind {
	case ApprovalKindSecrets:
		return a.Secrets
	case ApprovalKindShell:
		return a.Shell
	default:
		return nil
	}
}

func approvalMatches(ap Approval, fingerprint, project, sessionID string, now time.Time) bool {
	if ap.Fingerprint != fingerprint {
		return false
	}
	if !ap.ExpiresAt.IsZero() && !ap.ExpiresAt.After(now) {
		return false
	}
	switch ap.Scope {
	case ApprovalScopeProject:
		return project != "" && ap.Project != "" && samePath(ap.Project, project)
	case ApprovalScopeSession:
		return sessionID != "" && ap.SessionID == sessionID
	default:
		return false
	}
}

func samePath(a, b string) bool {
	return strings.TrimRight(a, "/") == strings.TrimRight(b, "/")
}

func parseApprovals(in *fileApprovals, now time.Time) (Approvals, error) {
	if in == nil {
		return Approvals{}, nil
	}
	secrets, err := parseApprovalList(ApprovalKindSecrets, in.Secrets, now)
	if err != nil {
		return Approvals{}, err
	}
	shell, err := parseApprovalList(ApprovalKindShell, in.Shell, now)
	if err != nil {
		return Approvals{}, err
	}
	return Approvals{Secrets: secrets, Shell: shell}, nil
}

func parseApprovalList(kind ApprovalKind, in []fileApproval, now time.Time) ([]Approval, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make([]Approval, 0, len(in))
	for i, fa := range in {
		ap, ok, err := parseOneApproval(kind, fa, now)
		if err != nil {
			return nil, fmt.Errorf("approvals.%s[%d]: %w", kind, i, err)
		}
		if ok {
			out = append(out, ap)
		}
	}
	return out, nil
}

func parseOneApproval(kind ApprovalKind, fa fileApproval, now time.Time) (Approval, bool, error) {
	if fa.Fingerprint == "" {
		return Approval{}, false, fmt.Errorf("fingerprint is required")
	}
	scope, err := parseApprovalScope(fa.Scope)
	if err != nil {
		return Approval{}, false, err
	}
	var expires time.Time
	if fa.ExpiresAt != "" {
		expires, err = time.Parse(time.RFC3339, fa.ExpiresAt)
		if err != nil {
			return Approval{}, false, fmt.Errorf("expires_at: %w", err)
		}
		if !expires.After(now) {
			return Approval{}, false, nil
		}
	}
	switch scope {
	case ApprovalScopeProject:
		if fa.Project == "" {
			return Approval{}, false, fmt.Errorf("project is required for project scope")
		}
	case ApprovalScopeSession:
		if fa.SessionID == "" {
			return Approval{}, false, fmt.Errorf("session_id is required for session scope")
		}
	}
	granted := fa.GrantedBy
	if granted == "" {
		granted = grantedByAskUser
	}
	return Approval{
		Kind:        kind,
		Fingerprint: fa.Fingerprint,
		Scope:       scope,
		Project:     fa.Project,
		SessionID:   fa.SessionID,
		ExpiresAt:   expires,
		GrantedBy:   granted,
	}, true, nil
}

func parseApprovalScope(s string) (ApprovalScope, error) {
	switch ApprovalScope(s) {
	case ApprovalScopeProject, ApprovalScopeSession:
		return ApprovalScope(s), nil
	case "":
		return "", fmt.Errorf("scope is required")
	default:
		return "", fmt.Errorf("unknown scope %q", s)
	}
}

func upsertApprovalList(base, overlay []fileApproval) []fileApproval {
	if overlay == nil {
		return append([]fileApproval(nil), base...)
	}
	if len(base) == 0 {
		return append([]fileApproval(nil), overlay...)
	}
	out := append([]fileApproval(nil), base...)
	index := map[string]int{}
	for i, a := range out {
		if a.Fingerprint != "" {
			index[a.Fingerprint] = i
		}
	}
	for _, a := range overlay {
		if a.Fingerprint == "" {
			out = append(out, a)
			continue
		}
		if i, ok := index[a.Fingerprint]; ok {
			out[i] = a
			continue
		}
		index[a.Fingerprint] = len(out)
		out = append(out, a)
	}
	return out
}

func mergeApprovalsPtr(base, overlay *fileApprovals) *fileApprovals {
	if base == nil && overlay == nil {
		return nil
	}
	out := fileApprovals{}
	if base != nil {
		out.Secrets = append([]fileApproval(nil), base.Secrets...)
		out.Shell = append([]fileApproval(nil), base.Shell...)
	}
	if overlay == nil {
		return &out
	}
	if overlay.Secrets != nil {
		out.Secrets = upsertApprovalList(out.Secrets, overlay.Secrets)
	}
	if overlay.Shell != nil {
		out.Shell = upsertApprovalList(out.Shell, overlay.Shell)
	}
	return &out
}
