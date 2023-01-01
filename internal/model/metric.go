package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
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
	Hash  string     `json:"hash,omitempty"`
}

func (m *Metric) String() string {
	return fmt.Sprintf("ID: %s, Type: %s, Value: %v", m.ID, m.MType, m.GetValue())
}

func (m *Metric) UpdateHash(hashKey string) {
	m.Hash = m.GetHash(hashKey)
}
func (m *Metric) GetHash(hashKey string) string {
	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		var hashCore string
		switch m.MType {
		case GaugeType:
			hashCore = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
		case CounterType:
			hashCore = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
		}
		h.Write([]byte(hashCore))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		return hash
	}
	return ""
}

func ParseMetric(name string, mType MetricType, value string) (metric *Metric, err error) {
	switch mType {
	case GaugeType:
		vt, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		gaugeVT := GaugeVT(vt)
		metric = &Metric{
			ID:    name,
			MType: GaugeType,
			Value: &gaugeVT,
		}
	case CounterType:
		vt, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return nil, err
		}
		counterVT := CounterVT(vt)
		metric = &Metric{
			ID:    name,
			MType: CounterType,
			Delta: &counterVT,
		}
	default:
		metric = &Metric{
			ID:    name,
			MType: NanType,
		}
	}
	return
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

func (m *Metric) CheckHash(hashKey string) error {
	if hashKey != "" {
		if m.GetHash(hashKey) != m.Hash {
			return errors.New("received metric hash is not matched")
		}
	}
	return nil
}
