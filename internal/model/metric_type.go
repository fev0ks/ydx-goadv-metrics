package model

const (
	GaugeTypeVal   = "gauge"
	CounterTypeVal = "counter"
	NanVal         = "NaN"
)

var (
	GaugeType   = MetricType(GaugeTypeVal)
	CounterType = MetricType(CounterTypeVal)
	NanType     = MetricType(NanVal)
)

type MetricType string

func MTypeValueOf(value string) MetricType {
	var mType MetricType
	switch value {
	case GaugeTypeVal:
		mType = GaugeType
	case CounterTypeVal:
		mType = CounterType
	default:
		mType = NanType
	}
	return mType
}

//func (mType *MetricType) Value() (driver.Value, error) {
//	// преобразуем []string в string
//	switch mType {
//	case &GaugeType:
//		return GaugeTypeVal, nil
//	case &CounterType:
//		return CounterTypeVal, nil
//	default:
//		return NanType, nil
//	}
//}
//
//func (mType *MetricType) Scan(value interface{}) error {
//	switch value {
//	case GaugeTypeVal:
//		*mType = GaugeType
//	case CounterTypeVal:
//		*mType = CounterType
//	default:
//		*mType = NanType
//	}
//	return nil
//}
