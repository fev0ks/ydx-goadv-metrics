package agent

import "context"

type MetricPoller interface {
	PollMetrics(ctx context.Context, done chan struct{})
}
