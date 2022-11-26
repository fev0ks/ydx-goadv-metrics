package repositories

import (
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"log"
	"sync"
)

var (
	cmrInitOnce sync.Once
	instance    *CommonMetricRepository
)

type CommonMetricRepository struct {
	Cache []*model.Metric
}

func NewCommonMetricsRepository() *CommonMetricRepository {
	return &CommonMetricRepository{}
}

func GetCommonMetricsRepository() *CommonMetricRepository {
	cmrInitOnce.Do(func() {
		instance = &CommonMetricRepository{}
	})
	return instance
}

func (cmr *CommonMetricRepository) SaveMetric(metrics []*model.Metric) {
	cmr.Cache = metrics
	log.Printf("Saved %d metrics\n", len(cmr.Cache))
}

func (cmr *CommonMetricRepository) GetMetricsList() []*model.Metric {
	return cmr.Cache
}
