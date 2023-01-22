package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/configs"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest"
	"github.com/fev0ks/ydx-goadv-metrics/internal"

	"github.com/spf13/pflag"
)

func main() {
	ctx := context.Background()
	log.Printf("Server args: %s\n", os.Args[1:])
	address := configs.GetAddress()
	var addressF string
	pflag.StringVarP(&addressF, "a", "a", configs.DefaultAddress, "Address of the server")

	restore := configs.GetDoReStore()
	var restoreF bool
	pflag.BoolVarP(&restoreF, "r", "r", configs.DefaultDoRestore, "Do autoBackup restore?")

	storeInterval := configs.GetStoreInterval()
	var storeIntervalF time.Duration
	pflag.DurationVarP(&storeIntervalF, "i", "i", configs.DefaultMetricStoreInterval, "Backup interval in sec")

	storeFile := configs.GetStoreFile()
	var storeFileF string
	pflag.StringVarP(&storeFileF, "f", "f", configs.DefaultStoreFile, "Path of Backup store file")

	hashKey := configs.GetHashKey()
	var hashKeyF string
	pflag.StringVarP(&hashKeyF, "k", "k", configs.DefaultHashKey, "Hash key")

	dbConfig := configs.GetDBConfig()
	var dbDsnF string
	pflag.StringVarP(&dbDsnF, "d", "d", configs.DefaultDBConfig, "Postgres DB DSN")

	pflag.Parse()

	if address == "" {
		address = addressF
	}
	if restore == nil {
		restore = &restoreF
	}
	if storeInterval == 0 {
		storeInterval = storeIntervalF
	}
	if address == "" {
		address = addressF
	}
	if restore == nil {
		restore = &restoreF
	}
	if storeFile == "" {
		storeFile = storeFileF
	}
	if hashKey == "" {
		hashKey = hashKeyF
	}
	if dbConfig == "" {
		dbConfig = dbDsnF
	}

	var autoBackup backup.AutoBackup
	stopCh := make([]chan struct{}, 0)
	toExecute := make([]func() error, 0)
	var repository server.MetricRepository
	if dbConfig != "" {
		repository = repositories.NewPgRepository(dbConfig, ctx)

	} else {
		repository = repositories.NewCommonRepository()
		autoBackup = backup.NewFileAutoBackup(storeInterval, repository, storeFile)
		if *restore {
			log.Println("trying to restore metrics...")
			err := autoBackup.Restore()
			if err != nil {
				log.Fatalf("failed to restore metrics: %v\n", err)
			}
		}
		stopCh = append(stopCh, autoBackup.Start())
		toExecute = append(toExecute, autoBackup.Backup)
	}

	mh := rest.NewMetricsHandler(ctx, repository, hashKey)
	hc := rest.NewHealthChecker(ctx, repository)

	router := rest.NewRouter()
	rest.HandleMetricRequests(router, mh)
	rest.HandleHeathCheck(router, hc)

	internal.ProperExitDefer(&internal.ExitHandler{
		ToStop:    stopCh,
		ToExecute: toExecute,
		ToClose:   []io.Closer{repository},
	})
	log.Printf("Server started on %s\n", address)
	log.Fatal(http.ListenAndServe(address, router))
}
