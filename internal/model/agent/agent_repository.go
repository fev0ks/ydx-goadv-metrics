package agent

import "github.com/fev0ks/ydx-goadv-metrics/internal/model"

type MetricRepository interface {
	SaveMetric(metrics *model.Metric)
	SaveMetrics(metrics []*model.Metric)
	GetMetrics() []*model.Metric
}
