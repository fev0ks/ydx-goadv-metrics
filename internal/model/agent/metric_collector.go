package agent

type MetricCollector interface {
	CollectMetrics() chan struct{}
}
