package server

import "github.com/fev0ks/ydx-goadv-metrics/internal/model"

type MetricRepository interface {
	SaveMetric(metric *model.Metric) error
	GetMetrics() map[string]*model.Metric
	GetMetric(name string) *model.Metric
	Clear()
}
