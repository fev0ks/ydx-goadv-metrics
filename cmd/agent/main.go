package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/config"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	var repository model.MetricRepository
	repository = repositories.GetCommonMetricsRepository()

	var metricCollector model.MetricCollector
	mcCtx, mcCancel := context.WithCancel(ctx)
	metricCollector = service.NewCommonMetricCollector(mcCtx, repository, config.GetReportInterval())
	stopCollectMetricsCh := metricCollector.CollectMetrics()

	var metricPoller model.MetricPoller
	mpCtx, mpCancel := context.WithCancel(ctx)
	metricPoller = service.NewCommonMetricPoller(mpCtx, repository, config.GetHost(), config.GetPort(), config.GetPollInterval())
	stopPollMetricsCh := metricPoller.PollMetrics()

	log.Println("Started")
	properExitDefer(&model.ExitHandler{
		Cancel: []context.CancelFunc{mcCancel, mpCancel},
		Stop:   []chan bool{stopCollectMetricsCh, stopPollMetricsCh},
	})
}

func properExitDefer(exitHandler *model.ExitHandler) {
	log.Println("Graceful func execution activated")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)
	select {
	case s := <-signals:
		log.Printf("Received a signal '%s', cancel active contexts\n", s)
		for _, cancel := range exitHandler.Cancel {
			cancel()
		}
		log.Printf("Received a signal '%s', stop active goroutines\n", s)
		for _, stop := range exitHandler.Stop {
			stop <- true
		}
	}
}
