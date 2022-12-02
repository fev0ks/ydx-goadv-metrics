package repositories

import (
	"testing"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGaugesMetrics(t *testing.T) {
	tc := struct {
		name          string
		metrics       []*model.Metric
		expectedValue map[string]float64
	}{
		"Should save and return Gauges",
		[]*model.Metric{
			model.NewGaugeMetric("test1", 123),
			model.NewGaugeMetric("test1", 321),
			model.NewGaugeMetric("test2", 123),
		},
		map[string]float64{"test1": 321, "test2": 123},
	}

	t.Run(tc.name, func(t *testing.T) {
		repository := NewCommonRepository()
		for _, metric := range tc.metrics {
			err := repository.SaveMetric(metric)
			require.NoError(t, err)
		}
		actualMetrics := repository.GetMetrics()
		for key, value := range tc.expectedValue {
			assert.Equal(t, model.GaugeVT(value), actualMetrics[key].Gauge)

			metricByName := repository.GetMetric(key)
			assert.Equal(t, model.GaugeVT(value), metricByName.Gauge)
		}
		repository.Clear()
		assert.Equal(t, 0, len(repository.GetMetrics()))
	})
}

func TestCounterMetrics(t *testing.T) {
	tc := struct {
		name          string
		metrics       []*model.Metric
		expectedValue map[string]uint32
	}{
		"Should save and return sum of Counter metric",
		[]*model.Metric{
			model.NewCounterMetric("test1", 1),
			model.NewCounterMetric("test1", 2),
			model.NewCounterMetric("test2", 0),
		},
		map[string]uint32{"test1": 3, "test2": 0},
	}

	t.Run(tc.name, func(t *testing.T) {
		repository := NewCommonRepository()
		for _, metric := range tc.metrics {
			err := repository.SaveMetric(metric)
			require.NoError(t, err)
		}
		actualMetrics := repository.GetMetrics()
		for key, value := range tc.expectedValue {
			assert.Equal(t, model.CounterVT(value), actualMetrics[key].Counter)

			metricByName := repository.GetMetric(key)
			assert.Equal(t, model.CounterVT(value), metricByName.Counter)
		}
		repository.Clear()
		assert.Equal(t, 0, len(repository.GetMetrics()))
	})
}

func TestNanMetrics(t *testing.T) {
	tc := struct {
		name    string
		metrics []*model.Metric
	}{
		"Should not save Nan metric",
		[]*model.Metric{
			model.NewNanMetric("test1"),
		},
	}

	t.Run(tc.name, func(t *testing.T) {
		repository := NewCommonRepository()
		for _, metric := range tc.metrics {
			err := repository.SaveMetric(metric)
			require.Error(t, err)
		}
	})
}
