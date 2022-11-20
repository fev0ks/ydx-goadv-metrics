package model

const (
	gaugeTypeVal   = "gauge"
	counterTypeVal = "counter"
)

var (
	GaugeType   = initGaugeType()
	CounterType = initCounterType()
)

type MetricType string

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
