package model

type MetricCollector interface {
	CollectMetrics() chan bool
}
