//go:build unix

package transport_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/transport"
)

func TestDefaultSocketPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		wantNonEmp bool
		wantSubstr string
	}{
		{name: "non empty", wantNonEmp: true},
		{name: "contains agentd", wantSubstr: "agentd"},
		{name: "ends with sock", wantSubstr: "agentd.sock"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := transport.DefaultSocketPath()
			if tt.wantNonEmp {
				assert.NotEmpty(t, got, "DefaultSocketPath()")
			}
			if tt.wantSubstr != "" {
				assert.Contains(t, got, tt.wantSubstr, "DefaultSocketPath()")
			}
			if tt.name == "ends with sock" {
				assert.Equal(t, "agentd.sock", filepath.Base(got), "DefaultSocketPath() base")
			}
		})
	}
}
