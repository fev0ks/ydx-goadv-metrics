package model

import (
	"fmt"
	"strconv"
)

type GaugeVT float64

type CounterVT uint64

type Metric struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *CounterVT `json:"delta,omitempty"`
	Value *GaugeVT   `json:"value,omitempty"`
}

//func (m *Metric) UnmarshalJSON(data []byte) error {
//
//	// чтобы избежать рекурсии при json.Unmarshal, объявляем новый тип
//	type MetricAlias Metric
//
//	metricValue := &struct {
//		*MetricAlias
//		// переопределяем поле внутри анонимной структуры
//		MType string `json:"type"`
//	}{
//		// задаём указатель на целевой объект
//		MetricAlias: (*MetricAlias)(m),
//	}
//
//	// вызываем стандартный Unmarshal
//	if err := json.Unmarshal(data, metricValue); err != nil {
//		return err
//	}
//	mType := MTypeValueOf(metricValue.MType)
//	m.MType = mType
//
//	return nil
//}

func (m *Metric) String() string {
	return fmt.Sprintf("ID: %s, Type: %s, Value: %v", m.ID, m.MType, m.GetValue())
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
		metric = &Metric{
			ID:    name,
			MType: NanType,
		}
	}
	return
}

func NewGaugeMetric(name string, value GaugeVT) *Metric {
	return &Metric{
		ID:    name,
		MType: GaugeType,
		Value: &value,
	}
}

func NewCounterMetric(name string, value CounterVT) *Metric {
	return &Metric{
		ID:    name,
		MType: CounterType,
		Delta: &value,
	}
}

func NewNanMetric(name string) *Metric {
	return &Metric{
		ID:    name,
		MType: NanType,
	}
}

func (m *Metric) GetValue() string {
	switch m.MType {
	case GaugeType:
		return fmt.Sprintf("%f", *m.Value)
	case CounterType:
		return fmt.Sprintf("%d", *m.Delta)
	default:
		return NanVal
	}
}

func (m *Metric) GetGenericValue() (value interface{}) {
	switch m.MType {
	case GaugeType:
		value = m.Value
	case CounterType:
		value = m.Delta
	default:
		value = ""
	}
	return
}
