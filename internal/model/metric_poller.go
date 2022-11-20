package model

type MetricPoller interface {
	PollMetrics() chan bool
}
