package agent

import "github.com/fev0ks/ydx-goadv-metrics/internal/model"

type MetricPoller interface {
	PollMetrics() chan struct{}
	SendMetric(metric *model.Metric) error
}
