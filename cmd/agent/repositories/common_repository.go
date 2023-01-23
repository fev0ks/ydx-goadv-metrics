package repositories

import (
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"sync"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type commonMetricRepository struct {
	*sync.RWMutex
	Cache map[string]*model.Metric
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
	cmr.Cache[metric.ID] = metric
}

func (cmr *commonMetricRepository) SaveMetrics(metrics []*model.Metric) {
	cmr.Lock()
	defer cmr.Unlock()
	for _, metric := range metrics {
		cmr.Cache[metric.ID] = metric
	}
}

func (cmr *commonMetricRepository) PullMetrics() map[string]*model.Metric {
	cmr.Lock()
	defer cmr.Unlock()
	metrics := cmr.Cache
	cmr.Cache = make(map[string]*model.Metric, 0)
	return metrics
}
