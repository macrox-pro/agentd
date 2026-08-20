package transport_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/transport"
)

func TestParseEndpoint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "unix scheme", in: "unix:///tmp/agentd.sock", want: "/tmp/agentd.sock"},
		{name: "bare path", in: "/tmp/agentd.sock", want: "/tmp/agentd.sock"},
		{name: "npipe scheme", in: `npipe:\\.\pipe\agentd`, want: `\\.\pipe\agentd`},
		{name: "bare pipe", in: `\\.\pipe\agentd-sid`, want: `\\.\pipe\agentd-sid`},
		{name: "empty", in: "", wantErr: true},
		{name: "unix empty path", in: "unix://", wantErr: true},
		{name: "npipe empty path", in: "npipe:", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := transport.ParseEndpoint(tt.in)
			if tt.wantErr {
				require.Error(t, err, "ParseEndpoint(%q)", tt.in)
				return
			}
			require.NoError(t, err, "ParseEndpoint(%q)", tt.in)
			assert.Equal(t, tt.want, got, "ParseEndpoint(%q)", tt.in)
		})
	}
}
