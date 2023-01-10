package main

import (
	"context"
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

	repository := repositories.NewCommonRepository()

	mh := rest.NewMetricsHandler(ctx, repository, hashKey)

	router := rest.NewRouter()
	rest.HandleMetricRequests(router, mh)

	var autoBackup backup.AutoBackup
	if dbConfig != "" {
		storage := repositories.NewPgRepository(dbConfig, ctx)
		hc := rest.NewHealthChecker(ctx, storage)
		rest.HandleHeathCheck(router, hc)
		autoBackup = backup.NewPgAutoBackup(storeInterval, repository, storage)
	} else {
		autoBackup = backup.NewFileAutoBackup(storeInterval, repository, storeFile)
	}
	if *restore {
		log.Println("trying to restore metrics...")
		err := autoBackup.Restore()
		if err != nil {
			log.Fatalf("failed to restore metrics: %v\n", err)
		}
	}
	backupCh := autoBackup.Start()

	internal.ProperExitDefer(&internal.ExitHandler{
		ToStop:    []chan struct{}{backupCh},
		ToExecute: []func() error{autoBackup.Backup},
		ToClose:   []io.Closer{repository},
	})
	log.Printf("Server started on %s\n", address)
	log.Fatal(http.ListenAndServe(address, router))
}
