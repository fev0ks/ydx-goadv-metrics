package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/config"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"log"
)

func main() {
	ctx := context.Background()
	var repository agent.MetricRepository
	repository = repositories.GetCommonMetricsRepository()

	var metricCollector agent.MetricCollector
	mcCtx, mcCancel := context.WithCancel(ctx)
	metricCollector = service.NewCommonMetricCollector(mcCtx, repository, config.GetReportInterval())
	stopCollectMetricsCh := metricCollector.CollectMetrics()

	var metricPoller agent.MetricPoller
	mpCtx, mpCancel := context.WithCancel(ctx)
	metricPoller = service.NewCommonMetricPoller(mpCtx, repository, config.GetHost(), config.GetPort(), config.GetPollInterval())
	stopPollMetricsCh := metricPoller.PollMetrics()

	log.Println("Agent started")
	internal.ProperExitDefer(&internal.ExitHandler{
		Cancel: []context.CancelFunc{mcCancel, mpCancel},
		Stop:   []chan bool{stopCollectMetricsCh, stopPollMetricsCh},
	})
}
