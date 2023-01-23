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
	"github.com/spf13/pflag"
)

func main() {
	ctx := context.Background()
	log.Printf("Agent args: %s\n", os.Args[1:])

	address := configs.GetServerAddress()
	var addressF string
	pflag.StringVarP(&addressF, "a", "a", configs.DefaultServerAddress, "Address of the server")

	reportInterval := configs.GetReportInterval()
	var reportIntervalF time.Duration
	pflag.DurationVarP(&reportIntervalF, "r", "r", configs.DefaultMetricReportInterval, "Report to server interval in sec")

	pollInterval := configs.GetPollInterval()
	var pollIntervalF time.Duration
	pflag.DurationVarP(&pollIntervalF, "p", "p", configs.DefaultMetricPollInterval, "Pool metrics interval in sec")

	hashKey := configs.GetHashKey()
	var hashKeyF string
	pflag.StringVarP(&hashKeyF, "k", "k", configs.DefaultHashKey, "Hash key")

	pflag.Parse()

	if address == "" {
		address = addressF
	}
	if reportInterval == 0 {
		reportInterval = reportIntervalF
	}
	if pollInterval == 0 {
		pollInterval = pollIntervalF
	}
	if hashKey == "" {
		hashKey = hashKeyF
	}

	repository := repositories.NewCommonMetricsRepository()
	metricFactory := service.NewMetricFactory(hashKey)
	mcCtx, mcCancel := context.WithCancel(ctx)
	metricCollector := service.NewCommonMetricCollector(mcCtx, repository, metricFactory, reportInterval)
	stopCollectMetricsCh := metricCollector.CollectMetrics()

	client := getClient(address)

	msCtx, msCancel := context.WithCancel(ctx)
	var metricSender sender.MetricSender
	if configs.UseBuffSenderClient() {
		metricSender = sender.NewBulkMetricSender(msCtx, client, configs.GetBuffBatchLimit(), configs.GetBuffSendInterval())
	} else {
		metricSender = sender.NewJSONMetricSender(msCtx, client)
	}

	mpCtx, mpCancel := context.WithCancel(ctx)
	metricPoller := service.NewCommonMetricPoller(mpCtx, client, metricSender, repository, pollInterval)
	stopPollMetricsCh := metricPoller.PollMetrics()

	log.Println("Agent started")
	internal.ProperExitDefer(&internal.ExitHandler{
		ToCancel: []context.CancelFunc{mcCancel, msCancel, mpCancel},
		ToStop:   []chan struct{}{stopCollectMetricsCh, stopPollMetricsCh},
	})

	<-ctx.Done()
}

func getClient(address string) *resty.Client {
	client := resty.New().
		SetBaseURL(fmt.Sprintf("http://%s", address)).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(3 * time.Second)
	return client
}
