package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var invokeDurationBuckets = []float64{
	0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30,
}

// Recorder observes invoke_sync and async_side histograms on a dedicated registry.
type Recorder struct {
	invoke *prometheus.HistogramVec
	async  *prometheus.HistogramVec
}

// NewRecorder registers invoke and async histograms on reg.
func NewRecorder(reg prometheus.Registerer) *Recorder {
	r := promauto.With(reg)
	return &Recorder{
		invoke: r.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "invoke_duration_seconds",
			Help:      "Hook Invoke sync pipeline duration in seconds.",
			Buckets:   invokeDurationBuckets,
		}, []string{"provider", "event_kind", "decision", "outcome"}),
		async: r.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "async_job_duration_seconds",
			Help:      "Async dispatch target job duration in seconds.",
			Buckets:   invokeDurationBuckets,
		}, []string{"target_kind", "result"}),
	}
}

// ObserveInvoke records one Invoke sync pipeline observation.
func (rec *Recorder) ObserveInvoke(provider, eventKind, decision, outcome string, seconds float64) {
	if rec == nil {
		return
	}
	rec.invoke.WithLabelValues(provider, eventKind, decision, outcome).Observe(seconds)
}

// ObserveAsync records one async target job observation.
func (rec *Recorder) ObserveAsync(targetKind, result string, seconds float64) {
	if rec == nil {
		return
	}
	rec.async.WithLabelValues(targetKind, result).Observe(seconds)
}
