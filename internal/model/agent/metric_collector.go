package agent

type MetricCollector interface {
	CollectMetrics() chan bool
}
