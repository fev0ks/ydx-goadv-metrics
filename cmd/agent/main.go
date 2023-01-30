package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/configs"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service/sender"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	"github.com/go-resty/resty/v2"
)

func main() {
	ctx := context.Background()
	log.Printf("Agent args: %s", os.Args[1:])
	appConfig := configs.InitAppConfig()

	done := make(chan struct{})
	repository := repositories.NewCommonMetricsRepository()
	metricFactory := service.NewMetricFactory(appConfig.HashKey)
	mcCtx, mcCancel := context.WithCancel(ctx)
	metricCollector := service.NewCommonMetricCollector(repository, metricFactory, appConfig.ReportInterval)
	metricCollector.CollectMetrics(mcCtx, done)

	client := getClient(appConfig.ServerAddress)

	var metricSender sender.Sender
	if appConfig.UseBuffSenderClient {
		metricSender = sender.NewBulkMetricSender(client, appConfig.BuffBatchLimit)
	} else {
		metricSender = sender.NewJSONMetricSender(client)
	}

	mpCtx, mpCancel := context.WithCancel(ctx)
	metricPoller := service.NewCommonMetricPoller(client, metricSender, repository, appConfig.PollInterval)
	metricPoller.PollMetrics(mpCtx, done)

	log.Println("Agent started")
	internal.ProperExitDefer(&internal.ExitHandler{
		ToCancel: []context.CancelFunc{mcCancel, mpCancel},
		ToStop:   []chan struct{}{done},
	})

	<-ctx.Done()
}

func getClient(address string) *resty.Client {
	client := resty.New().
		SetBaseURL(fmt.Sprintf("http://%s", address)).
		SetRetryCount(1).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
	return client
}
