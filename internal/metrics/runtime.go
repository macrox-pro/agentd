package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// RuntimeStats supplies scrape-time gauge/counter values via closures.
type RuntimeStats struct {
	AsyncQueueDepth       func() float64
	AsyncQueueCapacity    func() float64
	AsyncQueueDropped     func() float64
	TrajectoryQueueDropped func() float64
	SessionsActive        func() float64
	ConfigGeneration      func() float64
	CompiledRouteCount    func() float64
}

type runtimeCollector struct {
	stats RuntimeStats
}

func (c *runtimeCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *runtimeCollector) Collect(ch chan<- prometheus.Metric) {
	emitGauge := func(name, help string, fn func() float64) {
		if fn == nil {
			return
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(Namespace+"_"+name, help, nil, nil),
			prometheus.GaugeValue,
			fn(),
		)
	}
	emitCounter := func(name, help string, fn func() float64) {
		if fn == nil {
			return
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(Namespace+"_"+name, help, nil, nil),
			prometheus.CounterValue,
			fn(),
		)
	}
	emitGauge("async_queue_depth", "Queued async jobs waiting for workers.", c.stats.AsyncQueueDepth)
	emitGauge("async_queue_capacity", "Async queue capacity at daemon start.", c.stats.AsyncQueueCapacity)
	emitCounter("async_queue_dropped_total", "Async queue overflow drops.", c.stats.AsyncQueueDropped)
	emitCounter("trajectory_queue_dropped_total", "Trajectory ledger queue overflow drops.", c.stats.TrajectoryQueueDropped)
	emitGauge("sessions_active", "Sync Invoke session locks currently held.", c.stats.SessionsActive)
	emitGauge("config_generation", "Current compiled config generation.", c.stats.ConfigGeneration)
	emitGauge("compiled_route_count", "Routes in the active config snapshot.", c.stats.CompiledRouteCount)
}

// RegisterGoAndProcess registers standard Go runtime and process collectors on reg.
func RegisterGoAndProcess(reg *prometheus.Registry) {
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
}

// RegisterRuntime registers scrape-time daemon stats on reg.
func RegisterRuntime(reg *prometheus.Registry, stats RuntimeStats) {
	reg.MustRegister(&runtimeCollector{stats: stats})
}

// RegisterReloadCounter registers agentd_config_reload_total on reg.
func RegisterReloadCounter(reg prometheus.Registerer) *prometheus.CounterVec {
	return promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "config_reload_total",
		Help:      "Config reload attempts via Store.Reload.",
	}, []string{"result"})
}
