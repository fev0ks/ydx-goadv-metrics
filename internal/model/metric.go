package model

import (
	"errors"
	"fmt"
	"strconv"
)

type GaugeVT float64

type CounterVT uint64

type Metric struct {
	Name    string
	MType   MetricType
	Counter CounterVT
	Gauge   GaugeVT
}

func NewMetric(name string, mType MetricType, value string) (metric *Metric, err error) {
	switch mType {
	case GaugeType:
		vt, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		metric = NewGaugeMetric(name, GaugeVT(vt))
	case CounterType:
		vt, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return nil, err
		}
		metric = NewCounterMetric(name, CounterVT(vt))
	default:
		err = errors.New("metric type NaN is not supported")
	}
	return
}

func NewGaugeMetric(name string, value GaugeVT) *Metric {
	return &Metric{
		Name:  name,
		MType: GaugeType,
		Gauge: value,
	}
}

func NewCounterMetric(name string, value CounterVT) *Metric {
	return &Metric{
		Name:    name,
		MType:   CounterType,
		Counter: value,
	}
}

func (m *Metric) GetValue() string {
	switch m.MType {
	case GaugeType:
		return fmt.Sprintf("%v", fmt.Sprintf("%f", m.Gauge))
	case CounterType:
		return fmt.Sprintf("%v", m.Counter)
	default:
		return NanVal
	}
}
