package server

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type MetricRepository interface {
	SaveMetric(metric *model.Metric) error
	GetMetrics() map[string]*model.Metric
	GetMetricsList() []*model.Metric
	GetMetric(name string) *model.Metric
	HealthCheck(ctx context.Context) error
	Clear()
}
