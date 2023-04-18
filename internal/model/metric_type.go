package model

import "strings"

const (
	GaugeTypeVal   = "gauge"
	CounterTypeVal = "counter"
	NanVal         = "nan"
)

var (
	GaugeType   = MetricType(GaugeTypeVal)
	CounterType = MetricType(CounterTypeVal)
	NanType     = MetricType(NanVal)
)

type MetricType string

func MTypeValueOf(value string) MetricType {
	var mType MetricType
	switch strings.ToLower(value) {
	case GaugeTypeVal:
		mType = GaugeType
	case CounterTypeVal:
		mType = CounterType
	default:
		mType = NanType
	}
	return mType
}
