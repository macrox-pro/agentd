package guard

import (
	"encoding/json"
	"regexp"
)

// Finding is one detected credential-shaped string. The value itself is never
// carried — only a masked preview safe for prompts and logs.
type Finding struct {
	Rule   string
	Masked string
}

type rule struct {
	id   string
	name string
	re   *regexp.Regexp
}

// Order matters where prefixes overlap (Anthropic keys before OpenAI sk-).
var allRules = []rule{
	{id: "aws_key", name: "AWS access key ID", re: regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`)},
	{id: "github_pat", name: "GitHub token", re: regexp.MustCompile(`\bgh[opusr]_[A-Za-z0-9]{36,}\b`)},
	{id: "github_fine_grained", name: "GitHub fine-grained PAT", re: regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{22,}\b`)},
	{id: "slack_token", name: "Slack token", re: regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{10,}\b`)},
	{id: "stripe_live", name: "Stripe live key", re: regexp.MustCompile(`\b[sr]k_live_[A-Za-z0-9]{16,}\b`)},
	{id: "anthropic_key", name: "Anthropic API key", re: regexp.MustCompile(`\bsk-ant-[A-Za-z0-9-]{20,}\b`)},
	{id: "openai_key", name: "OpenAI API key", re: regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{20,}\b`)},
	{id: "google_api_key", name: "Google API key", re: regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}\b`)},
	{id: "private_key", name: "private key block", re: regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
	{id: "jwt", name: "JWT", re: regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`)},
	{id: "assigned_secret", name: "assigned secret literal", re: regexp.MustCompile(`(?i)(?:api[_-]?key|secret|token|password|passwd)["']?\s*[:=]\s*["'][^"']{8,}["']`)},
}

// Scan walks string values in tool-input JSON and reports credential-shaped matches.
// enabledRules filters by rule id; empty means all rules.
func Scan(input json.RawMessage, enabledRules []string) []Finding {
	var decoded any
	if err := json.Unmarshal(input, &decoded); err != nil {
		return nil
	}
	enabled := ruleSet(enabledRules)
	seen := map[string]bool{}
	var findings []Finding
	walkStrings(decoded, func(s string) {
		for _, r := range allRules {
			if enabled != nil && !enabled[r.id] {
				continue
			}
			if seen[r.id] {
				continue
			}
			if m := r.re.FindString(s); m != "" {
				seen[r.id] = true
				findings = append(findings, Finding{Rule: r.name, Masked: mask(m)})
			}
		}
	})
	return findings
}

func ruleSet(ids []string) map[string]bool {
	if len(ids) == 0 {
		return nil
	}
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}

func walkStrings(v any, visit func(string)) {
	switch t := v.(type) {
	case string:
		visit(t)
	case map[string]any:
		for _, val := range t {
			walkStrings(val, visit)
		}
	case []any:
		for _, val := range t {
			walkStrings(val, visit)
		}
	}
}

func mask(s string) string {
	if len(s) <= 12 {
		return s[:4] + "…"
	}
	return s[:6] + "…" + s[len(s)-4:]
}
