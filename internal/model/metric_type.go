package model

const (
	gaugeTypeVal   = "gauge"
	counterTypeVal = "counter"
	NanVal         = "NaN"
)

var (
	GaugeType   = MetricType(gaugeTypeVal)
	CounterType = MetricType(counterTypeVal)
	NanType     = MetricType(NanVal)
)

type MetricType string

func MTypeValueOf(value string) MetricType {
	var mType MetricType
	switch value {
	case gaugeTypeVal:
		mType = GaugeType
	case counterTypeVal:
		mType = CounterType
	default:
		mType = NanType
	}
	return mType
}
