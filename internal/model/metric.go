package model

type GaugeVT float64

type CounterVT uint64

type Metric struct {
	Name  string
	MType MetricType
	Delta CounterVT
	Value GaugeVT
}

func NewGaugeMetric(name string, value GaugeVT) *Metric {
	return &Metric{
		Name:  name,
		MType: GaugeType,
		Value: value,
	}
}

func NewCounterMetric(name string, value CounterVT) *Metric {
	return &Metric{
		Name:  name,
		MType: CounterType,
		Delta: value,
	}
}
