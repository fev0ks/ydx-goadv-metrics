package service

import (
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type MetricFactory interface {
	NewGaugeMetric(name string, value model.GaugeVT) *model.Metric
	NewCounterMetric(name string, value model.CounterVT) *model.Metric
}

type metricFactory struct {
	hashKey string
}

func NewMetricFactory(hashKey string) MetricFactory {
	return &metricFactory{hashKey: hashKey}
}

func (mf metricFactory) NewGaugeMetric(name string, value model.GaugeVT) *model.Metric {
	return newGaugeMetric(
		name,
		value,
		mf.hashKey,
	)
}

func (mf metricFactory) NewCounterMetric(name string, value model.CounterVT) *model.Metric {
	return newCounterMetric(
		name,
		value,
		mf.hashKey,
	)
}

func newGaugeMetric(name string, value model.GaugeVT, hashKey string) *model.Metric {
	metric := &model.Metric{
		ID:    name,
		MType: model.GaugeType,
		Value: &value,
	}
	metric.UpdateHash(hashKey)
	return metric
}

func newCounterMetric(name string, value model.CounterVT, hashKey string) *model.Metric {
	metric := &model.Metric{
		ID:    name,
		MType: model.CounterType,
		Delta: &value,
	}
	metric.UpdateHash(hashKey)
	return metric
}
