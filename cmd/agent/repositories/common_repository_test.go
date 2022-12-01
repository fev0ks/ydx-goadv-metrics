package repositories

import (
	"testing"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestCommonMetricRepository(t *testing.T) {
	testCases := []struct {
		name    string
		metrics []*model.Metric
	}{
		{
			"Should save and return Gauges",
			[]*model.Metric{
				model.NewGaugeMetric("test1", 123),
				model.NewGaugeMetric("test2", 123),
			},
		},
		{
			"Should save and return Counters",
			[]*model.Metric{
				model.NewCounterMetric("test1", 0),
				model.NewCounterMetric("test1", 123),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := NewCommonMetricsRepository()
			repository.SaveMetric(tc.metrics)
			actualMetrics := repository.GetMetricsList()
			assert.Equal(t, tc.metrics, actualMetrics)
		})
	}
}
