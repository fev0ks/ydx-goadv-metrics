package model

type MetricRepository interface {
	SaveMetricsList(metrics []*Metric)
	GetMetricsList() []*Metric
	//GetMetrics() map[string]model.Metric
	//GetMetric(string) model.Metric
}
