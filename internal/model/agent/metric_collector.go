package agent

type MetricCollector interface {
	CollectMetrics(done chan struct{})
}
