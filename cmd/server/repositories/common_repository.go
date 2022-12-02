package repositories

import (
	"fmt"
	"log"
	"sync"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type commonRepository struct {
	*sync.RWMutex
	storage map[string]*model.Metric
}

func NewCommonRepository() *commonRepository {
	return &commonRepository{
		&sync.RWMutex{},
		make(map[string]*model.Metric),
	}
}

func (cr *commonRepository) SaveMetric(metric *model.Metric) error {
	cr.Lock()
	defer cr.Unlock()
	switch metric.MType {
	case model.GaugeType:
		cr.storage[metric.Name] = metric
	case model.CounterType:
		if current, ok := cr.storage[metric.Name]; ok {
			current.Counter += metric.Counter
			log.Printf("%s = %+v", current.Name, current.Counter)
		} else {
			cr.storage[metric.Name] = metric
		}
	default:
		return fmt.Errorf("failed to save '%s' metric: '%v' type is not supported", metric.Name, metric.MType)
	}
	return nil
}

func (cr *commonRepository) GetMetrics() map[string]*model.Metric {
	cr.RLock()
	defer cr.RUnlock()
	return cr.storage
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
