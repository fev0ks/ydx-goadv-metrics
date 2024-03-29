package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/configs"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/middlewares"
	backup2 "github.com/fev0ks/ydx-goadv-metrics/internal/model/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"
	"github.com/fev0ks/ydx-goadv-metrics/internal/shutdown"
)

var (
	BuildVersion      = "N/A"
	BuildDate         = "N/A"
	BuildCommit       = "N/A"
	configPathEnvVar  = "CONFIG"
	defaultConfigPath = "cmd/server/config.json"
)

// go build -ldflags "-X github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.BuildVersion=v1 -X 'github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.BuildDate=$(date)' -X 'github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.BuildCommit=$(git rev-parse HEAD)'" github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.go
func main() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
	ctx := context.Background()
	log.Printf("Server args: %s", os.Args[1:])
	configPath := os.Getenv(configPathEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
	}
	appConfig, err := configs.InitAppConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	exitHandler := shutdown.NewExitHandler()
	var autoBackup backup2.IAutoBackup

	stopCh := make([]chan struct{}, 0)
	toExecute := make([]func() error, 0)

	var metricRepo repository.IMetricRepository
	if appConfig.DBConfig != "" {
		metricRepo, err = repositories.NewPgRepository(appConfig.DBConfig, ctx)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		metricRepo = repositories.NewCommonRepository()
		autoBackup = backup.NewFileAutoBackup(metricRepo, appConfig)
		if *appConfig.DoRestore {
			log.Println("trying to restore metrics...")
			err := autoBackup.Restore()
			if err != nil {
				log.Fatalf("failed to restore metrics: %v", err)
			}
		}
		stopCh = append(stopCh, autoBackup.Start())
		toExecute = append(toExecute, autoBackup.Backup)
	}
	exitHandler.ToStop = stopCh
	exitHandler.ToExecute = toExecute
	exitHandler.ToClose = []io.Closer{metricRepo}

	mh := rest.NewMetricsHandler(ctx, metricRepo, appConfig.HashKey)
	hc := rest.NewHealthChecker(ctx, metricRepo)

	router := rest.NewRouter()

	shutdownBlocker := middlewares.NewShutdownBlocker(exitHandler)
	router.Use(shutdownBlocker.BlockTillFinish)

	decrypter := middlewares.NewDecrypter(appConfig.PrivateKey)
	rest.HandleEncryptedMetricRequests(router, mh, decrypter)
	rest.HandleMetricRequests(router, mh)
	rest.HandleHeathCheck(router, hc)
	rest.HandlePprof(router)

	log.Printf("Server started on %s", appConfig.ServerAddress)
	shutdown.ProperExitDefer(exitHandler)
	server := &http.Server{Addr: appConfig.ServerAddress, Handler: router}
	exitHandler.ShutdownServerBeforeExit(server)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server closed with msg: '%v'", err)
	}
	<-ctx.Done()
}
