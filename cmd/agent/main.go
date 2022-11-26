package main

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/config"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	repository := repositories.GetCommonMetricsRepository()

	var metricCollector agent.MetricCollector
	mcCtx, mcCancel := context.WithCancel(ctx)
	metricCollector = service.NewCommonMetricCollector(mcCtx, repository, config.GetReportInterval())
	stopCollectMetricsCh := metricCollector.CollectMetrics()

	client := getClient()

	var metricPoller agent.MetricPoller
	mpCtx, mpCancel := context.WithCancel(ctx)
	metricPoller = service.NewCommonMetricPoller(mpCtx, client, repository, config.GetPollInterval())
	stopPollMetricsCh := metricPoller.PollMetrics()

	log.Println("Agent started")
	internal.ProperExitDefer(&internal.ExitHandler{
		Cancel: []context.CancelFunc{mcCancel, mpCancel},
		Stop:   []chan bool{stopCollectMetricsCh, stopPollMetricsCh},
	})
	<-ctx.Done()
}

func getClient() *resty.Client {
	client := resty.New().
		SetBaseURL(fmt.Sprintf("http://%s:%s", config.GetHost(), config.GetPort())).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(3 * time.Second)
	return client
}
