package agent

import "context"

type MetricCollector interface {
	CollectMetrics(ctx context.Context) chan struct{}
}
