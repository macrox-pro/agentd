package daemon

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	metricsReadHeaderTimeout = 5 * time.Second
	metricsShutdownTimeout   = 5 * time.Second
)

type metricsServer struct {
	srv    *http.Server
	listen string
}

func startMetricsServer(ctx context.Context, listen string, handler http.Handler) (*metricsServer, error) {
	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", listen)
	if err != nil {
		return nil, fmt.Errorf("metrics listen: %w", err)
	}
	srv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: metricsReadHeaderTimeout,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	ms := &metricsServer{srv: srv, listen: ln.Addr().String()}
	go func() {
		_ = srv.Serve(ln)
	}()
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), metricsShutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(shutCtx)
	}()
	return ms, nil
}

func (m *metricsServer) shutdown(ctx context.Context) {
	if m == nil || m.srv == nil {
		return
	}
	_ = m.srv.Shutdown(ctx)
}

func (m *metricsServer) addr() string {
	if m == nil {
		return ""
	}
	return m.listen
}
