package tui_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
	"github.com/macrox-pro/agentd/internal/install/tui"
	"github.com/macrox-pro/agentd/internal/provider"
)

func TestRunWizard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		findings []install.Finding
		contain  string
	}{
		{
			name:     "empty_discover",
			findings: nil,
			contain:  "doctor",
		},
		{
			name: "medium_only_skipped",
			findings: []install.Finding{{
				Provider:   provider.ClaudeCode,
				Confidence: install.ConfidenceMedium,
			}},
			contain: "doctor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := tui.RunWizard(context.Background(), tt.findings, tui.Deps{}, tui.WizardOptions{Out: &buf})
			require.NoError(t, err, "RunWizard(%q)", tt.name)
			assert.Contains(t, buf.String(), tt.contain, "RunWizard(%q)", tt.name)
		})
	}
}
