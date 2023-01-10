package repositories

import (
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"log"
	"sync"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type commonMetricRepository struct {
	*sync.RWMutex
	Cache []*model.Metric
}

func NewCommonMetricsRepository() agent.MetricRepository {
	return &commonMetricRepository{
		&sync.RWMutex{},
		make([]*model.Metric, 0),
	}
}

func (cmr *commonMetricRepository) SaveMetric(metrics []*model.Metric) {
	cmr.Lock()
	defer cmr.Unlock()
	cmr.Cache = metrics
	log.Printf("Saved %d metrics\n", len(cmr.Cache))
}

func (cmr *commonMetricRepository) GetMetricsList() []*model.Metric {
	cmr.RLock()
	defer cmr.RUnlock()
	return cmr.Cache
}
