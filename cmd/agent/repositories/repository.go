package repositories

import (
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"golang.org/x/exp/maps"
	"sync"
)

type commonMetricRepository struct {
	*sync.RWMutex
	cache map[string]*model.Metric
}

func NewCommonMetricsRepository() agent.MetricRepository {
	return &commonMetricRepository{
		&sync.RWMutex{},
		make(map[string]*model.Metric, 0),
	}
}

func (cmr *commonMetricRepository) SaveMetric(metric *model.Metric) {
	cmr.Lock()
	defer cmr.Unlock()
	cmr.cache[metric.ID] = metric
}

func (cmr *commonMetricRepository) SaveMetrics(metrics []*model.Metric) {
	cmr.Lock()
	defer cmr.Unlock()
	for _, metric := range metrics {
		cmr.cache[metric.ID] = metric
	}
}

func (cmr *commonMetricRepository) GetMetrics() []*model.Metric {
	cmr.Lock()
	defer cmr.Unlock()
	return maps.Values(cmr.cache)
}
