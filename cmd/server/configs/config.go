package configs

import (
	"os"
	"strconv"
	"time"
)

const (
	DefaultAddress             = "localhost:8080"
	DefaultMetricStoreInterval = 300 * time.Second
	DefaultStoreFile           = "/tmp/devops-metrics-db.json"
	DefaultDoRestore           = true
	DefaultHashKey             = ""
	DefaultDBConfig            = ""
)

func GetAddress() string {
	return os.Getenv("ADDRESS")
}

func GetStoreInterval() time.Duration {
	storeInterval := os.Getenv("STORE_INTERVAL")
	if storeInterval == "" {
		return 0
	}
	storeIntervalVal, err := strconv.Atoi(storeInterval)
	if err != nil {
		return 0
	}
	duration := time.Duration(storeIntervalVal) * time.Second
	return duration
}

func GetStoreFile() string {
	return os.Getenv("STORE_FILE")
}

func GetHashKey() string {
	return os.Getenv("KEY")
}

func GetDoReStore() *bool {
	doReStore := os.Getenv("RESTORE")
	if doReStore == "" {
		return nil
	}
	doReStoreVal, err := strconv.ParseBool(doReStore)
	if err != nil {
		return nil
	}
	return &doReStoreVal
}

func GetDBConfig() string {
	return os.Getenv("DATABASE_DSN")
}
