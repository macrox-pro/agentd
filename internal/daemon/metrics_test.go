package daemon_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/daemon"
	"github.com/macrox-pro/agentd/internal/metrics"
)

func freeListenAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "Listen")
	addr := ln.Addr().String()
	require.NoError(t, ln.Close(), "Close")
	return addr
}

func TestStartMetricsHTTP(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "disabled no listen",
			run: func(t *testing.T) {
				t.Parallel()
				en, listen, err := config.MetricsConfig{Enabled: false}.EffectiveListen("")
				require.NoError(t, err, "EffectiveListen")
				assert.False(t, en, "enabled")
				assert.Empty(t, listen, "listen")
			},
		},
		{
			name: "shutdown ErrServerClosed",
			run: func(t *testing.T) {
				reg := metrics.NewRegistry()
				metrics.RegisterGoAndProcess(reg)
				addr, shutdown, err := daemon.StartMetricsServerForTest(context.Background(), freeListenAddr(t), metrics.Handler(reg))
				require.NoError(t, err, "StartMetricsServerForTest")
				resp, err := http.Get("http://" + addr + "/metrics")
				require.NoError(t, err, "GET /metrics")
				require.NoError(t, resp.Body.Close(), "Close body")
				assert.Equal(t, http.StatusOK, resp.StatusCode, "status")

				shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				shutdown(shutdownCtx)

				req, err := http.NewRequestWithContext(shutdownCtx, http.MethodGet, "http://"+addr+"/metrics", nil)
				require.NoError(t, err, "NewRequest")
				_, err = http.DefaultClient.Do(req)
				require.Error(t, err, "GET after shutdown")
			},
		},
		{
			name: "ctx cancel stops server",
			run: func(t *testing.T) {
				reg := metrics.NewRegistry()
				metrics.RegisterGoAndProcess(reg)
				runCtx, cancel := context.WithCancel(context.Background())
				addr, _, err := daemon.StartMetricsServerForTest(runCtx, freeListenAddr(t), metrics.Handler(reg))
				require.NoError(t, err, "StartMetricsServerForTest")
				resp, err := http.Get("http://" + addr + "/metrics")
				require.NoError(t, err, "GET /metrics")
				require.NoError(t, resp.Body.Close(), "Close body")

				cancel()
				deadline := time.Now().Add(2 * time.Second)
				for time.Now().Before(deadline) {
					req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://"+addr+"/metrics", nil)
					require.NoError(t, err, "NewRequest")
					_, err = http.DefaultClient.Do(req)
					if err != nil {
						return
					}
					time.Sleep(20 * time.Millisecond)
				}
				t.Fatal("metrics server still listening after ctx cancel")
			},
		},
		{
			name: "bind in use",
			run: func(t *testing.T) {
				addr := freeListenAddr(t)
				hold, err := net.Listen("tcp", addr)
				require.NoError(t, err, "Listen hold")
				t.Cleanup(func() { _ = hold.Close() })

				reg := metrics.NewRegistry()
				_, _, err = daemon.StartMetricsServerForTest(context.Background(), addr, metrics.Handler(reg))
				require.Error(t, err, "StartMetricsServerForTest bind in use")
			},
		},
		{
			name: "reload does not rebind",
			run: func(t *testing.T) {
				socket, cfg := testSocket(t)
				metricsAddr := freeListenAddr(t)
				cfgYAML := fmt.Sprintf(`version: 1
metrics:
  enabled: true
  listen: %q
`, metricsAddr)
				require.NoError(t, os.WriteFile(cfg, []byte(cfgYAML), 0o600))

				errCh := make(chan error, 1)
				go func() {
					errCh <- daemon.Start(t.Context(), daemon.StartOptions{
						Socket:     socket,
						ConfigPath: cfg,
						Foreground: true,
						Version:    "test",
					})
				}()
				require.NoError(t, waitReady(t, socket, errCh, 5*time.Second), "Start")

				first := scrapeListen(t, metricsAddr)
				require.NotEmpty(t, first, "metrics body")

				newAddr := freeListenAddr(t)
				require.NoError(t, os.WriteFile(cfg, []byte(fmt.Sprintf(`version: 1
metrics:
  enabled: true
  listen: %q
`, newAddr)), 0o600))

				reloadCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_, reloadErr := daemon.Reload(reloadCtx, socket)
				require.NoError(t, reloadErr, "Reload")
				cancel()

				second := scrapeListen(t, metricsAddr)
				require.NotEmpty(t, second, "metrics still on original addr")
				_, getErr := http.Get("http://" + newAddr + "/metrics")
				require.Error(t, getErr, "new addr should not listen")

				stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer stopCancel()
				require.NoError(t, daemon.Stop(stopCtx, socket, 10*time.Second), "Stop")
				<-errCh
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func scrapeListen(t *testing.T, addr string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+addr+"/metrics", nil)
	require.NoError(t, err, "NewRequest")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "ReadAll")
	return string(b)
}