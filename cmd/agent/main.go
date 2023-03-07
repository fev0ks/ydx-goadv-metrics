package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/go-resty/resty/v2"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/configs"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service/sender"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
)

func main() {
	ctx := context.Background()
	log.Printf("Agent args: %s", os.Args[1:])
	appConfig := configs.InitAppConfig()

	done := make(chan struct{})
	repository := repositories.NewCommonMetricsRepository()
	metricFactory := service.NewMetricFactory(appConfig.HashKey)
	metricCollector := service.NewCommonMetricCollector(repository, metricFactory, appConfig.ReportInterval)
	metricCollector.CollectMetrics(done)

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
		ToCancel: []context.CancelFunc{mpCancel},
		ToStop:   []chan struct{}{done},
	})

	log.Fatal(http.ListenAndServe(appConfig.AgentAddress, nil))
}

func getClient(address string) *resty.Client {
	client := resty.New().
		SetBaseURL(fmt.Sprintf("http://%s", address)).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
	return client
}
