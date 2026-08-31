package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRegistry returns a dedicated Prometheus registry (never the global default).
func NewRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

// Handler returns an HTTP handler that exposes reg via Prometheus text format.
func Handler(reg prometheus.Gatherer) http.Handler {
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true})
}
