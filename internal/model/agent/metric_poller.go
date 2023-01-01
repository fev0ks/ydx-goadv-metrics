package agent

type MetricPoller interface {
	PollMetrics() chan struct{}
}
