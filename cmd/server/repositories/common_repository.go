package repositories

import (
	"context"
	"errors"
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
	ctx     *context.Context
	storage map[string]*model.Metric
}

func GetCommonRepository(ctx *context.Context) *CommonRepository {
	crInitOnce.Do(func() {
		crInstance = &CommonRepository{
			ctx:     ctx,
			storage: make(map[string]*model.Metric, 0),
		}
	})
	return crInstance
}

func (cr CommonRepository) SaveMetric(metric *model.Metric) error {
	switch metric.MType {
	case model.GaugeType:
		cr.storage[metric.Name] = metric
	case model.CounterType:
		if current, ok := cr.storage[metric.Name]; ok {
			//TODO counter, тип int64, новое значение должно добавляться к предыдущему (если оно ранее уже было известно серверу). - why?????
			//counter is not cleared on agent every reportInterval => 10 - 30 - 70 etc...
			current.Delta += metric.Delta
			log.Printf("%s = %+v", current.Name, current.Delta)
		} else {
			cr.storage[metric.Name] = metric
		}
	default:
		return errors.New(fmt.Sprintf("failed to save '%s' metric: '%v' type is not supported", metric.Name, metric.MType))
	}
	return nil
}

func (cr CommonRepository) GetMetrics() map[string]*model.Metric {
	return cr.storage
}

func (cr CommonRepository) GetMetric(name string) *model.Metric {
	return cr.storage[name]
}

func (cr CommonRepository) Clear() {
	cr.storage = make(map[string]*model.Metric, 0)
}
