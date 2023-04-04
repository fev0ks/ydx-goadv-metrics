package agent

import (
	"context"

	"github.com/fev0ks/ydx-goadv-metrics/internal/shutdown"
)

type MetricPoller interface {
	PollMetrics(ctx context.Context, exitHandler *shutdown.ExitHandler, done chan struct{})
}
