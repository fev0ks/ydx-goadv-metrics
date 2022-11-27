package repositories

import (
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"log"
	"sync"
)

var (
	crInitOnce sync.Once
	crInstance *CommonRepository
)

type CommonRepository struct {
	storage map[string]*model.Metric
}

func NewCommonRepository() *CommonRepository {
	return &CommonRepository{
		storage: make(map[string]*model.Metric),
	}
}

func GetCommonRepository() *CommonRepository {
	crInitOnce.Do(func() {
		crInstance = &CommonRepository{
			storage: make(map[string]*model.Metric),
		}
	})
	return crInstance
}

func (cr *CommonRepository) SaveMetric(metric *model.Metric) error {
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

func (cr *CommonRepository) GetMetrics() map[string]*model.Metric {
	return cr.storage
}

func (cr *CommonRepository) GetMetric(name string) *model.Metric {
	return cr.storage[name]
}

func (cr *CommonRepository) Clear() {
	cr.storage = make(map[string]*model.Metric)
}
