package repositories

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type metricGenerator struct {
	factory service.IMetricFactory
}

// 81          20681857 ns/op           32458 B/op        866 allocs/op
func BenchmarkCommonRepository_SaveMetrics(b *testing.B) {
	repo := NewCommonRepository()
	generator := metricGenerator{service.NewMetricFactory("")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		metrics := generator.generateMetrics(100)
		b.StartTimer()
		err := repo.SaveMetrics(metrics)
		if err != nil {
			b.Errorf("benchmark failed %v", err)
		}
	}
}

func (g *metricGenerator) generateMetrics(count int) []*model.Metric {
	metrics := make([]*model.Metric, 0, count)
	rand.Seed(time.Now().Unix())
	for i := 0; i < count; i++ {
		if rand.Intn(10) > 5 {
			metrics = append(metrics, g.factory.NewCounterMetric(fmt.Sprintf("counter %d", i), model.CounterVT(rand.Uint64())))
		} else {
			metrics = append(metrics, g.factory.NewGaugeMetric(fmt.Sprintf("gauge %d", i), model.GaugeVT(rand.Float64())))
		}
	}
	return metrics
}

func TestGaugesMetrics(t *testing.T) {
	//tc := struct {
	//	name          string
	//	metrics       []*model.Metric
	//	expectedValue map[string]float64
	//}{
	//	"Should save and return Gauges",
	//	[]*model.Metric{
	//		model.NewGaugeMetric("test1", 123),
	//		model.NewGaugeMetric("test1", 321),
	//		model.NewGaugeMetric("test2", 123),
	//	},
	//	map[string]float64{"test1": 321, "test2": 123},
	//}
	//
	//t.Run(tc.name, func(t *testing.T) {
	//	repository := NewCommonRepository()
	//	for _, metric := range tc.metrics {
	//		err := repository.SaveMetric(metric)
	//		require.NoError(t, err)
	//	}
	//	actualMetrics := repository.GetMetrics()
	//	for key, value := range tc.expectedValue {
	//		assert.Equal(t, model.GaugeVT(value), *actualMetrics[key].Value)
	//
	//		metricByName := repository.GetMetric(key)
	//		assert.Equal(t, model.GaugeVT(value), *metricByName.Value)
	//	}
	//	repository.Clear()
	//	assert.Equal(t, 0, len(repository.GetMetrics()))
	//})
}

func TestCounterMetrics(t *testing.T) {
	//tc := struct {
	//	name          string
	//	metrics       []*model.Metric
	//	expectedValue map[string]uint32
	//}{
	//	"Should save and return sum of Delta metric",
	//	[]*model.Metric{
	//		model.NewCounterMetric("test1", 1),
	//		model.NewCounterMetric("test1", 2),
	//		model.NewCounterMetric("test2", 0),
	//	},
	//	map[string]uint32{"test1": 3, "test2": 0},
	//}
	//
	//t.Run(tc.name, func(t *testing.T) {
	//	repository := NewCommonRepository()
	//	for _, metric := range tc.metrics {
	//		err := repository.SaveMetric(metric)
	//		require.NoError(t, err)
	//	}
	//	actualMetrics := repository.GetMetrics()
	//	for key, value := range tc.expectedValue {
	//		assert.Equal(t, model.CounterVT(value), *actualMetrics[key].Delta)
	//
	//		metricByName := repository.GetMetric(key)
	//		assert.Equal(t, model.CounterVT(value), *metricByName.Delta)
	//	}
	//	repository.Clear()
	//	assert.Equal(t, 0, len(repository.GetMetrics()))
	//})
}
