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
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	backup2 "github.com/fev0ks/ydx-goadv-metrics/internal/model/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

// go build -ldflags "-X github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.BuildVersion=v1 -X 'github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.BuildDate=$(date)' -X 'github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.BuildCommit=$(git rev-parse HEAD)'" github.com/fev0ks/ydx-goadv-metrics/cmd/server/main.go
func main() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
	ctx := context.Background()
	var err error
	log.Printf("Server args: %s", os.Args[1:])
	appConfig := configs.InitAppConfig()

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
		if appConfig.DoRestore {
			log.Println("trying to restore metrics...")
			err := autoBackup.Restore()
			if err != nil {
				log.Fatalf("failed to restore metrics: %v", err)
			}
		}
		stopCh = append(stopCh, autoBackup.Start())
		toExecute = append(toExecute, autoBackup.Backup)
	}

	mh := rest.NewMetricsHandler(ctx, metricRepo, appConfig.HashKey)
	hc := rest.NewHealthChecker(ctx, metricRepo)

	router := rest.NewRouter()
	rest.HandleMetricRequests(router, mh)
	rest.HandleHeathCheck(router, hc)
	rest.HandlePprof(router)
	internal.ProperExitDefer(&internal.ExitHandler{
		ToStop:    stopCh,
		ToExecute: toExecute,
		ToClose:   []io.Closer{metricRepo},
	})
	log.Printf("Server started on %s", appConfig.ServerAddress)
	log.Fatal(http.ListenAndServe(appConfig.ServerAddress, router))
}
