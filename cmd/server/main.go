package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/configs"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	ctx := context.Background()
	var err error
	log.Printf("Server args: %s", os.Args[1:])
	appConfig := configs.InitAppConfig()

	var autoBackup backup.AutoBackup
	stopCh := make([]chan struct{}, 0)
	toExecute := make([]func() error, 0)
	var repository server.MetricRepository
	if appConfig.DBConfig != "" {
		repository, err = repositories.NewPgRepository(appConfig.DBConfig, ctx)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		repository = repositories.NewCommonRepository()
		autoBackup = backup.NewFileAutoBackup(repository, appConfig)
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

	mh := rest.NewMetricsHandler(ctx, repository, appConfig.HashKey)
	hc := rest.NewHealthChecker(ctx, repository)

	router := rest.NewRouter()
	rest.HandleMetricRequests(router, mh)
	rest.HandleHeathCheck(router, hc)
	rest.HandlePprof(router)

	internal.ProperExitDefer(&internal.ExitHandler{
		ToStop:    stopCh,
		ToExecute: toExecute,
		ToClose:   []io.Closer{repository},
	})
	log.Printf("Server started on %s", appConfig.ServerAddress)
	log.Fatal(http.ListenAndServe(appConfig.ServerAddress, router))
}
