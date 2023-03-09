package repositories

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"
)

type CommonRepository struct {
	*sync.RWMutex
	storage map[string]*model.Metric
}

// NewCommonRepository - инициализация хранилища CommonRepository, реализующего repository.IMetricRepository, в виде Map структуры
// не сохраняет данные при выключении сервиса
func NewCommonRepository() repository.IMetricRepository {
	return &CommonRepository{
		&sync.RWMutex{},
		make(map[string]*model.Metric),
	}
}

func (cr *CommonRepository) HealthCheck(_ context.Context) error {
	return nil
}

func (cr *CommonRepository) SaveMetrics(metrics []*model.Metric) error {
	for _, metric := range metrics {
		err := cr.SaveMetric(metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cr *CommonRepository) SaveMetric(metric *model.Metric) error {
	cr.Lock()
	defer cr.Unlock()
	if metric == nil {
		return nil
	}
	switch metric.MType {
	case model.GaugeType:
		cr.storage[metric.ID] = metric
		log.Printf("Saved %v '%s' metric: '%s'", metric.MType, metric.ID, metric.GetValue())
	case model.CounterType:
		if current, ok := cr.storage[metric.ID]; ok {
			newValue := *current.Delta + *metric.Delta
			current.Delta = &newValue
			log.Printf("%s = %+v", current.ID, newValue)
		} else {
			cr.storage[metric.ID] = metric
		}
		log.Printf("Updated %v '%s' metric: '%s'", metric.MType, metric.ID, metric.GetValue())
	default:
		return fmt.Errorf("failed to save '%s' metric: '%v' type is not supported", metric.ID, metric.MType)
	}
	return nil
}

func (cr *CommonRepository) GetMetrics() (map[string]*model.Metric, error) {
	cr.RLock()
	defer cr.RUnlock()
	return cr.storage, nil
}

func (cr *CommonRepository) GetMetricsList() ([]*model.Metric, error) {
	cr.RLock()
	defer cr.RUnlock()
	metrics := make([]*model.Metric, 0, len(cr.storage))
	for _, v := range cr.storage {
		metrics = append(metrics, v)
	}
	return metrics, nil
}

func (cr *CommonRepository) GetMetric(name string) (*model.Metric, error) {
	cr.RLock()
	defer cr.RUnlock()
	return cr.storage[name], nil
}

func (cr *CommonRepository) Clear() error {
	cr.Lock()
	defer cr.Unlock()
	cr.storage = make(map[string]*model.Metric)
	return nil
}

func (cr *CommonRepository) Close() error {
	return nil
}
