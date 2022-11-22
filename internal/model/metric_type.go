package model

const (
	gaugeTypeVal   = "gauge"
	counterTypeVal = "counter"
	nanVal         = "NaN"
)

var (
	GaugeType   = initGaugeType()
	CounterType = initCounterType()
	NanType     = initNanType()
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

func initGaugeType() MetricType {
	var metricType MetricType
	metricType = gaugeTypeVal
	return metricType
}

func initCounterType() MetricType {
	var metricType MetricType
	metricType = counterTypeVal
	return metricType
}

func initNanType() MetricType {
	var metricType MetricType
	metricType = nanVal
	return metricType
}
