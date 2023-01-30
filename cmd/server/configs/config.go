package configs

import (
	"github.com/spf13/pflag"
	"os"
	"strconv"
	"time"
)

const (
	defaultAddress             = "localhost:8080"
	defaultMetricStoreInterval = 300 * time.Second
	defaultStoreFile           = "/tmp/devops-metrics-db.json"
	defaultDoRestore           = true
	defaultHashKey             = ""
	defaultDBConfig            = ""
)

type AppConfig struct {
	ServerAddress string
	StoreInterval time.Duration
	DoRestore     bool
	StoreFile     string
	HashKey       string
	DBConfig      string
}

func InitAppConfig() *AppConfig {
	address := getAddress()
	var addressF string
	pflag.StringVarP(&addressF, "a", "a", defaultAddress, "Address of the server")

	restore := getDoReStore()
	var restoreF bool
	pflag.BoolVarP(&restoreF, "r", "r", defaultDoRestore, "Do autoBackup restore?")

	storeInterval := getStoreInterval()
	var storeIntervalF time.Duration
	pflag.DurationVarP(&storeIntervalF, "i", "i", defaultMetricStoreInterval, "Backup interval in sec")

	storeFile := getStoreFile()
	var storeFileF string
	pflag.StringVarP(&storeFileF, "f", "f", defaultStoreFile, "Path of Backup store file")

	hashKey := getHashKey()
	var hashKeyF string
	pflag.StringVarP(&hashKeyF, "k", "k", defaultHashKey, "Hash key")

	dbConfig := getDBConfig()
	var dbDsnF string
	pflag.StringVarP(&dbDsnF, "d", "d", defaultDBConfig, "Postgres DB DSN")

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

	return &AppConfig{
		ServerAddress: address,
		StoreInterval: storeInterval,
		DoRestore:     *restore,
		StoreFile:     storeFile,
		HashKey:       hashKey,
		DBConfig:      dbConfig,
	}
}

func getAddress() string {
	return os.Getenv("ADDRESS")
}

func getStoreInterval() time.Duration {
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

func getStoreFile() string {
	return os.Getenv("STORE_FILE")
}

func getHashKey() string {
	return os.Getenv("KEY")
}

func getDoReStore() *bool {
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

func getDBConfig() string {
	return os.Getenv("DATABASE_DSN")
}
