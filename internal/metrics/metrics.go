// Package metrics exposes Prometheus scrape handlers and runtime collectors for the daemon.
//
// Owns: custom registry factory, promhttp handler, Go/process collectors, runtime ConstMetrics,
// PR B: Recorder histograms and config reload counter registration.
// Must not: daemon lifecycle, YAML compile, dispatch routing, importing dispatch/config/daemon.
//
// Invariants:
//   - Never use prometheus.DefaultRegisterer.
//   - internal/metrics is a leaf package (stdlib + client_golang only).
//
// Entry: NewRegistry, Handler, RegisterGoAndProcess, RegisterRuntime, NewRecorder.
// See DESIGN.md §1.5 (other), §5.
package metrics

const Namespace = "agentd"
