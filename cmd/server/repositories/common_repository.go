package repositories

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
	"log"
	"sync"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type commonRepository struct {
	*sync.RWMutex
	storage map[string]*model.Metric
}

func NewCommonRepository() server.MetricRepository {
	return &commonRepository{
		&sync.RWMutex{},
		make(map[string]*model.Metric),
	}
}

func (cr *commonRepository) HealthCheck(_ context.Context) error {
	return nil
}

func (cr *commonRepository) SaveMetric(metric *model.Metric) error {
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

func (cr *commonRepository) GetMetrics() map[string]*model.Metric {
	cr.RLock()
	defer cr.RUnlock()
	return cr.storage
}

func (cr *commonRepository) GetMetricsList() []*model.Metric {
	cr.RLock()
	defer cr.RUnlock()
	metrics := make([]*model.Metric, 0, len(cr.storage))
	for _, v := range cr.storage {
		metrics = append(metrics, v)
	}
	return metrics
}

func (cr *commonRepository) GetMetric(name string) *model.Metric {
	cr.RLock()
	defer cr.RUnlock()
	return cr.storage[name]
}

func (cr *commonRepository) Clear() {
	cr.Lock()
	defer cr.Unlock()
	cr.storage = make(map[string]*model.Metric)
}
