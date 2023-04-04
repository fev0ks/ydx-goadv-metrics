package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/configs"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/rest"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/rest/clients"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service/sender"
	"github.com/fev0ks/ydx-goadv-metrics/internal/shutdown"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

const (
	configPathEnvVar  = "CONFIG"
	defaultConfigPath = "cmd/agent/config.json"
)

// go run -ldflags "-X github.com/fev0ks/ydx-goadv-metrics/cmd/agent/main.BuildVersion=v1 -X 'github.com/fev0ks/ydx-goadv-metrics/cmd/agent/main.BuildDate=$(date)' -X 'github.com/fev0ks/ydx-goadv-metrics/cmd/agent/main.BuildCommit=$(git rev-parse HEAD)'" github.com/fev0ks/ydx-goadv-metrics/cmd/agent/main.go
func main() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
	ctx := context.Background()
	log.Printf("Agent args: %s", os.Args[1:])
	configPath := os.Getenv(configPathEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
	}
	appConfig, err := configs.InitAppConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	exitHandler := shutdown.NewExitHandler()

	done := make(chan struct{})
	exitHandler.ToStop = []chan struct{}{done}

	repository := repositories.NewCommonMetricsRepository()
	metricFactory := service.NewMetricFactory(appConfig.HashKey)
	metricCollector := service.NewCommonMetricCollector(repository, metricFactory, appConfig.ReportInterval)
	metricCollector.CollectMetrics(done)

	client := clients.CreateClient(appConfig.ServerAddress)
	encryptor := rest.NewEncryptor(appConfig.PublicKey)

	var metricSender sender.Sender
	if appConfig.UseBuffSenderClient {
		metricSender = sender.NewBulkMetricSender(client, appConfig.BuffBatchLimit, encryptor)
	} else {
		metricSender = sender.NewJSONMetricSender(client, encryptor)
	}

	mpCtx, mpCancel := context.WithCancel(ctx)
	exitHandler.ToCancel = []context.CancelFunc{mpCancel}

	metricPoller := service.NewCommonMetricPoller(client, metricSender, repository, appConfig.PollInterval)
	metricPoller.PollMetrics(mpCtx, exitHandler, done)

	log.Println("Agent started")
	shutdown.ProperExitDefer(exitHandler)
	log.Fatal(http.ListenAndServe(appConfig.AgentAddress, nil))
}
