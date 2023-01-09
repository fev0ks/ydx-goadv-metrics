package server

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type MetricRepository interface {
	SaveMetric(metric *model.Metric) error
	GetMetrics() (map[string]*model.Metric, error)
	GetMetricsList() ([]*model.Metric, error)
	GetMetric(name string) (*model.Metric, error)
	HealthCheck(ctx context.Context) error
	Clear() error
	Close() error
}
