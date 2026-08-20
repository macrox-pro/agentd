package guard_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/guard"
)

const fakeAWSKey = "AKIAIOSFODNN7EXAMPLE"

func TestScan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sample   string
		wantRule string
		rules    []string
		wantHit  bool
	}{
		{name: "aws", sample: "key=" + fakeAWSKey, wantRule: "AWS access key ID", wantHit: true},
		{name: "github", sample: "ghp_" + strings.Repeat("a1", 18), wantRule: "GitHub token", wantHit: true},
		{name: "clean", sample: "ls -la", wantHit: false},
		{name: "aws filtered out", sample: "key=" + fakeAWSKey, rules: []string{"jwt"}, wantHit: false},
		{name: "aws filtered in", sample: "key=" + fakeAWSKey, rules: []string{"aws_key"}, wantRule: "AWS access key ID", wantHit: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input, err := json.Marshal(map[string]string{"v": tt.sample})
			require.NoError(t, err, "Marshal")
			findings := guard.Scan(input, tt.rules)
			if !tt.wantHit {
				assert.Empty(t, findings, "Scan(%q)", tt.sample)
				return
			}
			require.NotEmpty(t, findings, "Scan(%q)", tt.sample)
			found := false
			for _, f := range findings {
				if f.Rule == tt.wantRule {
					found = true
					assert.NotContains(t, f.Masked, fakeAWSKey, "mask must not leak full secret")
				}
			}
			assert.True(t, found, "Scan(%q) want rule %q got %v", tt.sample, tt.wantRule, findings)
		})
	}
}

func TestRuleIDsAlignWithDefaults(t *testing.T) {
	t.Parallel()
	assert.Equal(t, config.DefaultSecretsRules, guard.RuleIDs(), "guard rule ids must match config.DefaultSecretsRules")
}
