package main

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/configs"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service/sender"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	repository := repositories.NewCommonMetricsRepository()

	var metricCollector agent.MetricCollector
	mcCtx, mcCancel := context.WithCancel(ctx)
	metricCollector = service.NewCommonMetricCollector(mcCtx, repository, configs.GetReportInterval())
	stopCollectMetricsCh := metricCollector.CollectMetrics()

	client := getClient()

	var metricPoller agent.MetricPoller
	mpCtx, mpCancel := context.WithCancel(ctx)
	metricSender := sender.NewJsonMetricSender(mpCtx, client)
	metricPoller = service.NewCommonMetricPoller(mpCtx, client, metricSender, repository, configs.GetPollInterval())
	stopPollMetricsCh := metricPoller.PollMetrics()

	log.Println("Agent started")
	internal.ProperExitDefer(&internal.ExitHandler{
		ToCancel: []context.CancelFunc{mcCancel, mpCancel},
		ToStop:   []chan struct{}{stopCollectMetricsCh, stopPollMetricsCh},
	})

	<-ctx.Done()
}

func getClient() *resty.Client {
	client := resty.New().
		SetBaseURL(fmt.Sprintf("http://%s", configs.GetServerAddress())).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(3 * time.Second)
	return client
}
