package repository

import (
	"context"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

// IMetricRepository - интерфейс для работы с хранилищем метрик
//
//go:generate mockgen -source=server_repository.go -destination=../../../test/mock/dynamic/repository/server_repository.go -package=repository
type IMetricRepository interface {
	// SaveMetric - сохранение состояния метрики
	SaveMetric(ctx context.Context, metric *model.Metric) error
	// SaveMetrics - сохранение состояния метрик
	SaveMetrics(ctx context.Context, metrics []*model.Metric) error
	// GetMetrics - получение метрик в виде Map структуры, где key - имя метрики, value - сама метрика
	GetMetrics(ctx context.Context) (map[string]*model.Metric, error)
	// GetMetricsList - получение списка метрик
	GetMetricsList(ctx context.Context) ([]*model.Metric, error)
	// GetMetric - получение метрики по ее имени
	GetMetric(ctx context.Context, name string) (*model.Metric, error)
	// HealthCheck - проверка доступности хранилища метрик
	HealthCheck(ctx context.Context) error
	// Clear - удаление всех метрик
	Clear(ctx context.Context) error
	// Close - завершение работы с хранилищем метрик
	Close() error
}
