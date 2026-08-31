package metrics_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/metrics"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		wantSubstr string
	}{
		{name: "go collector", wantSubstr: "go_goroutines"},
		{name: "process collector", wantSubstr: "process_cpu_seconds_total"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reg := metrics.NewRegistry()
			metrics.RegisterGoAndProcess(reg)
			srv := httptest.NewServer(metrics.Handler(reg))
			t.Cleanup(srv.Close)

			resp, err := srv.Client().Get(srv.URL)
			require.NoError(t, err, "TestHandler(%q)", tt.name)
			t.Cleanup(func() { _ = resp.Body.Close() })
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "TestHandler(%q)", tt.name)
			require.Contains(t, string(body), tt.wantSubstr, "TestHandler(%q)", tt.name)
		})
	}
}

func TestHandler_openMetrics(t *testing.T) {
	t.Parallel()

	reg := metrics.NewRegistry()
	metrics.RegisterGoAndProcess(reg)
	srv := httptest.NewServer(metrics.Handler(reg))
	t.Cleanup(srv.Close)

	resp, err := srv.Client().Get(srv.URL)
	require.NoError(t, err, "TestHandler_openMetrics")
	t.Cleanup(func() { _ = resp.Body.Close() })
	ct := resp.Header.Get("Content-Type")
	require.True(t, strings.Contains(ct, "openmetrics") || strings.Contains(ct, "text/plain"),
		"TestHandler_openMetrics content-type %q", ct)
}
