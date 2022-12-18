package configs

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultAddress             = "localhost:8080"
	defaultMetricStoreInterval = 300 * time.Second
	defaultStoreFile           = "/tmp/devops-metrics-db.json"
	defaultDoRestore           = true
)

func GetAddress() string {
	host := os.Getenv("ADDRESS")
	if host == "" {
		return defaultAddress
	}
	return host
}

func GetStoreInterval() time.Duration {
	storeInterval := os.Getenv("STORE_INTERVAL")
	if storeInterval == "" {
		return defaultMetricStoreInterval
	}
	storeIntervalVal, err := strconv.Atoi(storeInterval)
	if err != nil {
		return defaultMetricStoreInterval
	}
	return time.Duration(storeIntervalVal) * time.Second
}

func GetStoreFile() string {
	host := os.Getenv("STORE_FILE ")
	if host == "" {
		return defaultStoreFile
	}
	return host
}

func GetDoReStore() bool {
	doReStore := os.Getenv("RESTORE")
	if doReStore == "" {
		return defaultDoRestore
	}
	doReStoreVal, err := strconv.ParseBool(doReStore)
	if err != nil {
		return defaultDoRestore
	}
	return doReStoreVal
}
