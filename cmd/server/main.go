package main

import (
	"context"
	backup2 "github.com/fev0ks/ydx-goadv-metrics/cmd/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/configs"
	"log"
	"net/http"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
)

func main() {
	ctx := context.Background()

	repository := repositories.NewCommonRepository()
	mh := rest.NewMetricsHandler(ctx, repository)

	router := rest.NewRouter()
	rest.HandleMetricRequests(router, mh)

	backup := backup2.NewAutoBackup(configs.GetStoreFile(), configs.GetStoreInterval(), repository)
	if configs.GetDoReStore() {
		err := backup.Restore()
		if err != nil {
			log.Fatalf("failed to restore metrics backup: %v\n", err)
		}
	}
	backupCh := backup.Start()

	internal.ProperExitDefer(&internal.ExitHandler{ToStop: []chan struct{}{backupCh}, ToExecute: []func() error{backup.Backup}})
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(configs.GetAddress(), router))
}
